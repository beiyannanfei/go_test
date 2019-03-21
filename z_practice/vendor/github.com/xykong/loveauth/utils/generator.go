package utils

import (
	"github.com/xykong/snowflake"
	"github.com/xykong/loveauth/settings"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
	"strings"
)

var idGenerator *snowflake.Node

func initSnowflake() {

	var setting = settings.Get("loveauth")

	var epoch = setting.GetInt64("snowflake.Epoch")
	var nodeId = setting.GetInt64("snowflake.NodeId")

	snowflake.Epoch = epoch
	// max loveauth server is 64
	snowflake.NodeBits = 6
	// max time is 39bit, 17.4 year from 2013
	snowflake.StepBits = 18
	snowflake.Order = snowflake.OrderStepNodeTime

	var err error
	// Create a new Node with a Node number of 1
	idGenerator, err = snowflake.NewNode(nodeId)
	if err != nil {
		fmt.Println(err)

		logrus.WithFields(logrus.Fields{
			"epoch":  epoch,
			"nodeId": nodeId,
		}).Fatal("Snowflake init failed. Please check configs loveauth.json snowflake section.")
		return
	}

	// Generate a snowflake ID.
	id := idGenerator.Generate()

	logrus.WithFields(logrus.Fields{
		"epoch":      epoch,
		"nodeId":     nodeId,
		"Int64  ID":  id,
		"Base2  ID":  id.Base2(),
		"Base64  ID": id.Base64(),
		"ID Time":    id.Time(),
		"ID Node":    id.Node(),
		"ID Step":    id.Step(),
	}).Info("Snowflake init succeed.")
}

func GenerateId() snowflake.ID {
	if idGenerator == nil {
		initSnowflake()
	}

	return idGenerator.Generate()
}

func IdString(id int64) string {

	return snowflake.ID(id).BaseReadable()
}

func GenerateAuthCode() string {

	i := rand.Intn(1000000)

	return generateAuthCodeRaw(i)
}

func generateAuthCodeRaw(code int) string {

	code = code % 1000000

	result := fmt.Sprintf("%06d", code)

	logrus.Infof("GenerateAuthCode: %v", result)

	return result
}

func GeneratePaymentSequence(prefix ... interface{}) string {

	var result []string

	for _, v := range prefix {
		result = append(result, fmt.Sprintf("%v", v))
	}

	//result = append(result, fmt.Sprintf("%v.%v", time.Now().Format(time.RFC3339), GenerateId()))	//防止订单号出现特殊字符
	result = append(result, fmt.Sprintf("%v.%v", time.Now().Format("20060102150405"), GenerateId()))

	return strings.Join(result, "_")
}
