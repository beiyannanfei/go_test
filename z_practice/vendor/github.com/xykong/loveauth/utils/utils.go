package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func I(a ...interface{}) []interface{} {
	return a
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GenerateNonceStr() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	count := len(chars)
	rand.Seed(time.Now().UnixNano())
	var nonceStr = ""
	for i := 0; i < 32; i++ {
		nonceStr = fmt.Sprintf("%s%c", nonceStr, chars[rand.Intn(count)])
	}

	return nonceStr
}

func LoadTemplates(filename string, list interface{}) error {

	t := reflect.TypeOf(list)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Slice {
		return errors.New("input should be slice, not " + t.Kind().String())
	}

	sl := reflect.Indirect(reflect.ValueOf(list))
	typeOfT := sl.Type().Elem()

	templatePath := settings.GetString("lovepay", "template.path")
	excelFileName := path.Join(templatePath, filename)

	basename := filepath.Base(filename)
	basename = strings.TrimSuffix(basename, filepath.Ext(basename))

	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"pwd":     I(os.Getwd()),
			"os.Args": os.Args,
			"path":    excelFileName,
			"list":    reflect.TypeOf(list),
		}).Error("loadTemplates open file failed.")
		return err
	}

	for _, sheet := range xlFile.Sheets {

		if len(xlFile.Sheet) > 1 &&
			strings.ToLower(sheet.Name) != "data" &&
			strings.ToLower(sheet.Name) != basename {

			continue
		}

		keys := make([]string, 0, sheet.MaxCol)

		for idxRow, row := range sheet.Rows {

			var data []string
			ptr := reflect.New(typeOfT).Interface()
			val := reflect.Indirect(reflect.ValueOf(ptr))

			for idxCell, cell := range row.Cells {

				text := cell.String()
				// read first line of column name
				if idxRow == 0 {
					if len(text) > 0 {
						keys = append(keys, text)
					}
					continue
				}

				// skip null row in excel
				if idxCell == 0 {
					if text == "" {
						break
					}
				}

				field, ok := val.Type().FieldByName(keys[idxCell])
				if !ok {
					continue
				}

				switch field.Type.Kind() {
				case reflect.Int:
					i, err := strconv.Atoi(text)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":   err,
							"column":  keys[idxCell],
							"idxCell": idxCell,
							"idxRow":  idxRow,
							"text":    text,
						}).Info("loadTemplates read number failed.")
					}
					data = append(data, fmt.Sprintf(`"%v":%v`, keys[idxCell], i))
				case reflect.String:
					data = append(data, fmt.Sprintf(`"%v":"%v"`, keys[idxCell], text))
				case reflect.Struct:
					t, _ := cell.GetTime(false)
					data = append(data, fmt.Sprintf(`"%v":"%v"`, keys[idxCell], t.Format(time.RFC3339)))
				case reflect.Int64:
					i, err := strconv.Atoi(text)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":   err,
							"column":  keys[idxCell],
							"idxCell": idxCell,
							"idxRow":  idxRow,
							"text":    text,
						}).Info("loadTemplates read number failed.")
					}
					data = append(data, fmt.Sprintf(`"%v":%v`, keys[idxCell], i))
				case reflect.Float64:
					i, err := strconv.ParseFloat(text, 64)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"error":   err,
							"column":  keys[idxCell],
							"idxCell": idxCell,
							"idxRow":  idxRow,
							"text":    text,
						}).Info("loadTemplates read float64 failed.")
					}
					data = append(data, fmt.Sprintf(`"%v":%v`, keys[idxCell], i))
				default:
					logrus.WithFields(logrus.Fields{
						"case":  "default",
						"field": "field",
					}).Info("loadTemplates")
				}
			}

			if len(data) == 0 {
				continue
			}

			jsonString := fmt.Sprintf("{%v}", strings.Join(data, ","))

			//var shopItemTemplate ShopItemTemplate
			err = json.Unmarshal([]byte(jsonString), ptr)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":      err,
					"jsonString": jsonString,
				}).Error("loadTemplates parse template file failed.")
				return err
			}

			//logrus.WithFields(logrus.Fields{
			//	"json":    jsonString,
			//	"element": ptr,
			//	"typeOfT": typeOfT,
			//}).Info("loadTemplates.")

			s := reflect.ValueOf(ptr).Elem()
			sl = reflect.Append(sl, s)

			dstPtrValue := reflect.ValueOf(list)
			dstValue := reflect.Indirect(dstPtrValue)
			dstValue.Set(sl)
		}
	}

	//logrus.Warnf("loadTemplates: %v", list)

	return nil
}