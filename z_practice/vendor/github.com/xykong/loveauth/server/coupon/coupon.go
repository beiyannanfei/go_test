package coupon

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	coupon2 "github.com/xykong/loveauth/coupon"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
curl -H "Content-Type: application/json" -X POST \
        -d '{
        "Token":"3D1CE81678F75C458FB813716893582D"
        }' \
        http://localhost:8080/api/v1/verify/token
 */

// Binding from JSON
type RequestGenerate struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Count uint32 `form:"count" json:"count" binding:"required"`
}

//
// in: body
// swagger:parameters coupon_generate
type DoCouponRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestGenerate RequestGenerate `form:"RequestGenerate" json:"RequestGenerate" binding:"required"`
}

// Binding from JSON
type RequestVerify struct {
	Name   string `form:"name" json:"name" binding:"required"`
	Coupon string `form:"coupon" json:"coupon" binding:"required"`
	Mark   bool   `form:"mark" json:"mark" binding:""`
}

//
// in: body
// swagger:parameters coupon_verify
type RequestVerifyBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestVerify RequestVerify `form:"RequestVerify" json:"RequestVerify" binding:"required"`
}

// A coupon.ResponseCoupon is an response message to client.
// swagger:response ResponseCoupon
type Response struct {
	// in: body
	// description: coupon response
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		File    string `json:"file"`
	}
}

func Start(group *gin.RouterGroup) {

	group.POST("/generate", generate)
	group.POST("/verify", verifynew)
	group.POST("/generatenew", generatenew)
	group.POST("/mark", markCoupon)
	group.GET("/query", query)
	group.GET("/download/:filename", download)
}

// swagger:route POST /coupon/generate coupon coupon_generate
//
// Generate coupons to file and database by given arguments.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: ResponseCoupon
//       400: ResponseCoupon
//       500: ResponseCoupon
func generate(c *gin.Context) {

	var json RequestGenerate

	// validation
	if err := c.BindJSON(&json); err != nil {
		resp := Response{}
		resp.Body.Message = "Generate coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	if json.Name == "" {
		resp := Response{}
		resp.Body.Message = "Coupon Name is empty!"
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := Response{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	file, err := ioutil.TempFile(os.TempDir(), "coupon")
	if err != nil {
		resp := Response{}
		resp.Body.Message = "coupon file generate failed: " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	defer file.Close()

	_, _ = file.WriteString(fmt.Sprintf("%d,%s\r\n", json.Count, json.Name))

	var table = storage.Coupon{}
	table.SetName(json.Name)
	db.AutoMigrate(&table)

	var setting = settings.Get("loveauth")
	var salt = setting.GetString("coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(json.Name+salt)))
	var generator = coupon2.NewEncoding(key, 6, 0, 5, 4, 5)

	var coupons [] *storage.Coupon
	for i := uint32(1); i < json.Count+1; i++ {

		var code = generator.EncodeToString(uint64(i))
		//var code = generator.EncodePartsToString(1, 1, uint64(i))
		var coupon = storage.Coupon{CouponId: i, Coupon: code, Used: time.Time{}}

		coupons = append(coupons, &coupon)
		_, _ = file.WriteString(fmt.Sprintf("%d,%s\r\n", i, code))
	}

	batchCount := 3000
	end := batchCount

	for i := 0; i < len(coupons); {

		if end > len(coupons) {
			end = len(coupons)
		}

		err = table.BulkInsert(coupons[i:end])
		if err != nil {
			resp := Response{}
			resp.Body.Message = "coupon insert failed: " + err.Error()
			resp.Body.Code = int64(errors.Failed)
			c.JSON(http.StatusInternalServerError, resp.Body)

			return
		}

		i = end
		end += batchCount
	}

	resp := Response{}
	resp.Body.Message = "Coupon generate start successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.File = fmt.Sprintf("%s/%s", strings.Replace(c.Request.URL.Path, "generate", "download", -1),
		filepath.Base(file.Name()))
	c.JSON(http.StatusOK, resp.Body)
}

// swagger:route POST /coupon/verify coupon coupon_verify
//
// Verify coupon received from server.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: ResponseCoupon
//       400: ResponseCoupon
//       500: ResponseCoupon
func verify(c *gin.Context) {

	var json RequestVerify

	// validation
	if err := c.BindJSON(&json); err != nil {
		resp := Response{}
		resp.Body.Message = "Verify coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	if json.Name == "" {
		resp := Response{}
		resp.Body.Message = "Coupon Name is empty."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := Response{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	var table = storage.Coupon{}
	table.SetName(json.Name)

	var setting = settings.Get("loveauth")
	var salt = setting.GetString("coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(json.Name+salt)))
	var generator = coupon2.NewEncoding(key, 4, 1, 4, 4, 4)

	var code = json.Coupon
	index, err := generator.DecodeString(code)
	if err != nil {
		resp := Response{}
		resp.Body.Message = "Coupon decode failed: " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	table.CouponId = uint32(index)
	if table.QueryCoupon() == nil {
		resp := Response{}
		resp.Body.Message = fmt.Sprintf("Coupon verify failed, not found: %d, %s", index, code)
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	zero := time.Time{}
	if table.Used != zero {
		resp := Response{}
		resp.Body.Message = "Coupon verify failed, coupon used at: " + table.Used.String()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	if json.Mark {
		if table.MarkCoupon() == nil {
			resp := Response{}
			resp.Body.Message = "Coupon mark failed."
			resp.Body.Code = int64(errors.Failed)
			c.JSON(http.StatusOK, resp.Body)

			return
		}
	}

	resp := Response{}
	resp.Body.Message = "Coupon verify successfully!"
	resp.Body.Code = int64(errors.Ok)
	c.JSON(http.StatusOK, resp.Body)
}

// swagger:route GET /coupon/download/{:filename} coupon coupon_download
//
// Download coupon file generated by server
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200:
//       403:
func download(c *gin.Context) {

	fileName := c.Param("filename")
	targetPath := filepath.Join(os.TempDir(), fileName)
	//This check is for example, I not sure is it can prevent all possible filename attacks
	// - will be much better if real filename will not come from user side. I not even tried this code
	if !strings.HasPrefix(filepath.Clean(targetPath), os.TempDir()) {
		c.String(403, "Looks like you attacking me")
		return
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		c.String(403, "File not exist.")
		return
	}

	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/octet-stream")
	c.File(targetPath)
}
