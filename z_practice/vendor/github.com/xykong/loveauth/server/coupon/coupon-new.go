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

// Binding from JSON
type RequestGenerateNew struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Count    uint32 `form:"count" json:"count" binding:"required"`
	Channel  uint64 `form:"channel" json:"channel"`
	Activity uint64 `form:"activity" json:"activity"`
}

//
// in: body
// swagger:parameters coupon_generatenew
type DoCouponRequestNewBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestGenerateNew RequestGenerateNew `form:"RequestGenerateNew" json:"RequestGenerateNew" binding:"required"`
}

// A coupon.ResponseCoupon is an response message to client.
// swagger:response ResponseCouponNew
type ResponseNew struct {
	// in: body
	// description: coupon responsenew
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		File    string `json:"file"`
	}
}

//curl -X POST "http://10.1.16.69:8080/api/v1/coupon/generatenew" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"activity\": 123, \"channel\": 123, \"count\": 12332, \"name\": \"qqqqqqq\"}"

// swagger:route POST /coupon/generatenew coupon coupon_generatenew
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
//       200: ResponseCouponNew
//       400: ResponseCouponNew
//       500: ResponseCouponNew
func generatenew(c *gin.Context) {

	var json RequestGenerateNew

	// validation
	if err := c.BindJSON(&json); err != nil {
		resp := Response{}
		resp.Body.Message = "Generate coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	if json.Name == "" {
		resp := ResponseNew{}
		resp.Body.Message = "Coupon Name is empty!"
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := ResponseNew{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	file, err := ioutil.TempFile(os.TempDir(), "coupon")
	if err != nil {
		resp := ResponseNew{}
		resp.Body.Message = "coupon file generate failed: " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	defer file.Close()

	_, _ = file.WriteString(fmt.Sprintf("%s,渠道[%d],活动[%d],兑换码个数[%d]\r\n", json.Name, json.Channel, json.Activity, json.Count))

	var table = storage.Coupon{}
	table.SetName(makeCouponName(json.Channel, json.Activity))
	db.AutoMigrate(&table)

	var setting = settings.Get("loveauth")
	var salt = setting.GetString("coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(salt)))
	var generator = coupon2.NewEncoding(key, 6, 0, 5, 4, 5)

	var coupons [] *storage.Coupon
	for i := uint32(1); i < json.Count+1; i++ {

		var code = generator.EncodePartsToString(json.Channel, json.Activity, uint64(i))

		channel, activity, index, err := generator.DecodePartsString(code)
		if err != nil || channel != json.Channel || activity != json.Activity || uint64(i) != index {

			resp := ResponseNew{}
			resp.Body.Message = "coupon check failed: " + err.Error()
			resp.Body.Code = int64(errors.Failed)
			c.JSON(http.StatusInternalServerError, resp.Body)

			return
		}
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

	resp := ResponseNew{}
	resp.Body.Message = "Coupon generate start successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.File = fmt.Sprintf("%s/%s", strings.Replace(c.Request.URL.Path, "generate", "download", -1),
		filepath.Base(file.Name()))
	c.JSON(http.StatusOK, resp.Body)
}

// Binding from JSON
type RequestVerifyNew struct {
	Coupon string `form:"coupon" json:"coupon" binding:"required"`
}

//
// in: body
// swagger:parameters coupon_verifynew
type RequestVerifyNewBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestVerifyNew RequestVerifyNew `form:"RequestVerifyNew" json:"RequestVerifyNew" binding:"required"`
}

// A coupon.VerifyResponseNew is an response message to client.
// swagger:response VerifyResponseNew
type VerifyResponseNew struct {
	// in: body
	// description: coupon verifyresponsenew
	Body struct {
		Code     int64  `json:"code"`
		Message  string `json:"message"`
		Count    uint64 `json:"count"`
		Channel  uint64 `json:"channel"`
		Activity uint64 `json:"activity"`
		Index    uint32 `json:"index"`
	}
}

//curl -X POST "http://10.1.16.69:8080/api/v1/coupon/verifynew" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"coupon\": \"3YRA3-YVKK-FAW71\", \"globalId\": 123456}"

// swagger:route POST /coupon/verifynew coupon coupon_verifynew
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
//       200: ResponseCouponNew
//       400: ResponseCouponNew
//       500: ResponseCouponNew
func verifynew(c *gin.Context) {

	var json RequestVerifyNew

	// validation
	if err := c.BindJSON(&json); err != nil {
		resp := VerifyResponseNew{}
		resp.Body.Message = "Verify coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var setting = settings.Get("loveauth")
	var salt = setting.GetString("coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(salt)))
	var generator = coupon2.NewEncoding(key, 6, 0, 5, 4, 5)

	channel, activity, index, err := generator.DecodePartsString(json.Coupon)
	if err != nil {
		resp := VerifyResponseNew{}
		resp.Body.Message = fmt.Sprintf("Coupon decode failed: " + err.Error())
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := VerifyResponseNew{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	var table = storage.Coupon{}
	table.SetName(makeCouponName(channel, activity))

	table.CouponId = uint32(index)
	if table.QueryCoupon() == nil {
		resp := VerifyResponseNew{}
		resp.Body.Message = fmt.Sprintf("Coupon verify failed, not found: %d, %s", index, json.Coupon)
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	//table.Used = time.Now()
	table.Count++

	//updates := map[string]interface{}{"used": table.Used, "count": table.Count}
	//if table.GlobalId == 0 {
	//	updates["global_id"] = json.GlobalId
	//}
	//if err = db.Table(table.TableName()).Where("id = ?", table.ID).Updates(updates).Error; err != nil {
	//
	//	resp := Response{}
	//	resp.Body.Message = "Coupon mark failed."
	//	resp.Body.Code = int64(errors.Failed)
	//	c.JSON(http.StatusOK, resp.Body)
	//
	//	return
	//}

	resp := VerifyResponseNew{}
	resp.Body.Message = "Coupon verify successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Count = table.Count
	resp.Body.Channel = channel
	resp.Body.Activity = activity
	resp.Body.Index = table.CouponId

	c.JSON(http.StatusOK, resp.Body)
}

// Binding from JSON
type RequestMark struct {
	GlobalId uint64 `form:"globalId" json:"globalId" binding:"required"`
	Index    uint32 `form:"index" json:"index" binding:"required"`
	Channel  uint64 `form:"channel" json:"channel"`
	Activity uint64 `form:"activity" json:"activity"`
}

//
// in: body
// swagger:parameters coupon_mark
type RequestMarkBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestMark RequestMark `form:"RequestMark" json:"RequestMark" binding:"required"`
}

// A coupon.MarkResponse is an response message to client.
// swagger:response MarkResponse
type MarkResponse struct {
	// in: body
	// description: coupon markresponse
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
}

//curl -X POST "http://10.1.16.69:8080/api/v1/coupon/mark" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"activity\": 123, \"channel\": 123, \"globalId\": 123123123, \"index\": 1}"

// swagger:route POST /coupon/mark coupon coupon_mark
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
//       200: MarkResponse
//       400: MarkResponse
//       500: MarkResponse
func markCoupon(c *gin.Context) {
	var json RequestMark

	// validation
	if err := c.BindJSON(&json); err != nil {
		resp := MarkResponse{}
		resp.Body.Message = "mark coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := MarkResponse{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	var table = storage.Coupon{}
	table.SetName(makeCouponName(json.Channel, json.Activity))

	table.CouponId = json.Index
	if table.QueryCoupon() == nil {
		resp := MarkResponse{}
		resp.Body.Message = fmt.Sprintf("Coupon mark failed, not found: %d, %d_%d", json.Index, json.Channel, json.Activity)
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	table.Count++

	updates := map[string]interface{}{"used": time.Now(), "count": table.Count}
	if table.GlobalId == 0 {
		updates["global_id"] = json.GlobalId
	}
	if err := db.Table(table.TableName()).Where("coupon_id = ?", table.CouponId).Updates(updates).Error; err != nil {

		resp := MarkResponse{}
		resp.Body.Message = "Coupon mark failed."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp := MarkResponse{}
	resp.Body.Message = "Coupon mark success."
	resp.Body.Code = int64(errors.Ok)

	c.JSON(http.StatusOK, resp.Body)
}

func makeCouponName(channel, activity uint64) string {

	return fmt.Sprintf("%d_%d", channel, activity)
}

// Binding from JSON
type RequestQuery struct {
	Coupon string `form:"coupon" json:"coupon" binding:"required"`
}

//
// in: body
// swagger:parameters coupon_query
type RequestQueryBodyParams struct {
	//
	// swagger:allOf
	// in: body
	RequestQuery RequestQuery `form:"RequestQuery" json:"RequestQuery" binding:"required"`
}

// A coupon.ResponseQuery is an response message to client.
// swagger:response ResponseQuery
type ResponseQuery struct {
	// in: body
	// description: coupon responsequery
	Body struct {
		Code    int64          `json:"code"`
		Message string         `json:"message"`
		Info    storage.Coupon `json:"info"`
	}
}

// http://10.1.16.69:8080/api/v1/coupon/query?coupon=N419T-J17X-PA4TL
//
// swagger:route GET /coupon/query coupon coupon_query
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
//       200: ResponseQuery
//       400: ResponseQuery
//       500: ResponseQuery
func query(c *gin.Context) {

	var json RequestQuery

	// validation
	if err := c.Bind(&json); err != nil {
		resp := ResponseQuery{}
		resp.Body.Message = "query coupon failed! " + err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusBadRequest, resp.Body)

		return
	}

	var setting = settings.Get("loveauth")
	var salt = setting.GetString("coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(salt)))
	var generator = coupon2.NewEncoding(key, 6, 0, 5, 4, 5)

	channel, activity, index, err := generator.DecodePartsString(json.Coupon)
	if err != nil {
		resp := ResponseQuery{}
		resp.Body.Message = fmt.Sprintf("Coupon decode failed: " + err.Error())
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	var db = storage.AuthDatabase()
	if db == nil {
		resp := ResponseQuery{}
		resp.Body.Message = "Database is not initialized."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusInternalServerError, resp.Body)

		return
	}

	var table = storage.Coupon{}
	table.SetName(makeCouponName(channel, activity))

	table.CouponId = uint32(index)
	if table.QueryCoupon() == nil {
		resp := ResponseQuery{}
		resp.Body.Message = fmt.Sprintf("Coupon query failed, not found: %d, %s", index, json.Coupon)
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp := ResponseQuery{}
	resp.Body.Message = "Coupon query successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Info = table

	c.JSON(http.StatusOK, resp.Body)
}
