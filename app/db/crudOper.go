package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/jinzhu/gorm"
)

//对数据库的增删改查操作

func do() {
	//execSql()
	//create(&UserInfo{UserName: "abcd", Age: 29, Addr: "天通苑"})
	//create1()
	//delete1()

	insData()

	delData1()

	delData2()

	updateData1()

	updateData2()

	query()
}

func query() { //查询
	log.Info("============================ query ============================")
	type queryTest struct {
		gorm.Model
		UserName string `json:"UserName"` //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}
	reCreateTb(&queryTest{})
	dataArr := make([]queryTest, 0, 10)
	dataArr = append(dataArr, queryTest{UserName: "AAA", Age: 20, Addr: "北京0"},
		queryTest{UserName: "BBB", Age: 21, Addr: "北京1"}, queryTest{UserName: "CCC", Age: 22, Addr: "北京2"}, queryTest{UserName: "DDD", Age: 23, Addr: "北京3"},
		queryTest{UserName: "EEE", Age: 24, Addr: "北京4"}, queryTest{UserName: "FFF", Age: 25, Addr: "北京5"}, queryTest{UserName: "HHH", Age: 26, Addr: "北京6"},
		queryTest{UserName: "III", Age: 27, Addr: "北京7"}, queryTest{UserName: "JJJ", Age: 28, Addr: "北京8"}, queryTest{UserName: "KKK", Age: 29, Addr: "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //创建10条数据
		log.WithFields(log.Fields{"data": value}).Debug("*-*-*-*-*- updateData1 insert data finish")
		respones := dbClient.Create(&value)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": value}).Debug("create a record success!")
	}

	var data1 queryTest
	affectRow := dbClient.First(&data1).RowsAffected //获取第一条数据
	log.WithFields(log.Fields{
		"affectRow": affectRow,
		"data1":     data1,
	}).Info("first query finish!")
	data1.UserName = "AAAAAAAAAAAAA"
	affectRow = dbClient.Save(&data1).RowsAffected //更新上边查询出来数据的一个字段
	log.WithFields(log.Fields{
		"affectRow": affectRow,
		"data1":     data1,
	}).Info("save data finish!")


}

func updateData2() { //更新多条数据
	log.Info("============================ updateData2 ============================")
	type updateTest struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}
	reCreateTb(&updateTest{})
	dataArr := make([]updateTest, 0, 10)
	dataArr = append(dataArr, updateTest{UserName: "AAA", Age: 20, Addr: "北京0"},
		updateTest{UserName: "BBB", Age: 21, Addr: "北京1"}, updateTest{UserName: "CCC", Age: 22, Addr: "北京2"}, updateTest{UserName: "DDD", Age: 23, Addr: "北京3"},
		updateTest{UserName: "EEE", Age: 24, Addr: "北京4"}, updateTest{UserName: "FFF", Age: 25, Addr: "北京5"}, updateTest{UserName: "HHH", Age: 26, Addr: "北京6"},
		updateTest{UserName: "III", Age: 27, Addr: "北京7"}, updateTest{UserName: "JJJ", Age: 28, Addr: "北京8"}, updateTest{UserName: "KKK", Age: 29, Addr: "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //创建10条数据
		log.WithFields(log.Fields{"data": value}).Debug("*-*-*-*-*- updateData1 insert data finish")
		respones := dbClient.Create(&value)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": value}).Debug("create a record success!")
	}

	updateCond := updateTest{UserName: "CCC"}
	updateData := updateTest{UserName: "ABCDEF", Age: 33, Addr: "固安县"}
	affectRow := dbClient.Model(&updateTest{}).Where(&updateCond).UpdateColumns(updateData).RowsAffected
	log.WithFields(log.Fields{
		"affectRow":  affectRow,
		"updateData": updateData,
		"updateCond": updateCond,
	}).Info("updateData2 one recode finish!")
}

func updateData1() { //更新一条数据中的多个值
	log.Info("============================ updateData1 ============================")
	type updateTest struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}
	reCreateTb(&updateTest{})
	dataArr := make([]updateTest, 0, 10)
	dataArr = append(dataArr, updateTest{UserName: "AAA", Age: 20, Addr: "北京0"},
		updateTest{UserName: "BBB", Age: 21, Addr: "北京1"}, updateTest{UserName: "CCC", Age: 22, Addr: "北京2"}, updateTest{UserName: "DDD", Age: 23, Addr: "北京3"},
		updateTest{UserName: "EEE", Age: 24, Addr: "北京4"}, updateTest{UserName: "FFF", Age: 25, Addr: "北京5"}, updateTest{UserName: "HHH", Age: 26, Addr: "北京6"},
		updateTest{UserName: "III", Age: 27, Addr: "北京7"}, updateTest{UserName: "JJJ", Age: 28, Addr: "北京8"}, updateTest{UserName: "KKK", Age: 29, Addr: "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //创建10条数据
		log.WithFields(log.Fields{"data": value}).Debug("*-*-*-*-*- updateData1 insert data finish")
		respones := dbClient.Create(&value)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": value}).Debug("create a record success!")
	}

	updateCond1 := updateTest{Id: 6}
	updateData1 := updateTest{UserName: "ABCDEFG", Age: 32, Addr: "廊坊市"}
	//更新主键Id为6的数据为{UserName: "ABCDEFG", Age: 32, Addr: "廊坊市"}
	affectRow := dbClient.Model(&updateCond1).Updates(updateData1).RowsAffected
	log.WithFields(log.Fields{
		"affectRow":   affectRow,
		"updateData1": updateData1,
		"updateCond1": updateCond1,
	}).Info("update one record finish!")

	updateCond2 := updateTest{UserName: "BBB"} ////todo 如果结构中没有Id(即主键)字段，则更新所有数据(不会更新UserName为BBB的数据)
	updateData2 := updateTest{UserName: "ABCDEFG", Age: 32, Addr: "廊坊市"}
	affectRow = dbClient.Model(&updateCond2).Updates(updateData2).RowsAffected
	log.WithFields(log.Fields{
		"affectRow":   affectRow,
		"updateData2": updateData2,
		"updateCond2": updateCond2,
	}).Info("update data without primary key field")
}

func delData2() { //批量删除
	log.Info("============================ delData2 ============================")
	type delTest struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}

	reCreateTb(&delTest{})

	dataArr := make([]delTest, 0, 10)
	dataArr = append(dataArr, delTest{UserName: "AAA", Age: 20, Addr: "北京0"},
		delTest{UserName: "BBB", Age: 21, Addr: "北京1"}, delTest{UserName: "CCC", Age: 22, Addr: "北京2"}, delTest{UserName: "DDD", Age: 23, Addr: "北京3"},
		delTest{UserName: "EEE", Age: 24, Addr: "北京4"}, delTest{UserName: "FFF", Age: 25, Addr: "北京5"}, delTest{UserName: "HHH", Age: 26, Addr: "北京6"},
		delTest{UserName: "III", Age: 27, Addr: "北京7"}, delTest{UserName: "JJJ", Age: 28, Addr: "北京8"}, delTest{UserName: "KKK", Age: 29, Addr: "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //创建10条数据
		log.WithFields(log.Fields{"data": value}).Debug("*-*-*-*-*- delData1 insert data finish")
		respones := dbClient.Create(&value)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": value}).Debug("create a record success!")
	}

	condition := delTest{UserName: "BBB"}
	affectRow := dbClient.Where(condition).Delete(delTest{}).RowsAffected
	log.WithFields(log.Fields{
		"affectRow": affectRow,
		"condition": condition,
	}).Info("delete one data finish")
}

func delData1() { //单条删除数据
	log.Info("============================ delData1 ============================")
	type delTest struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}

	reCreateTb(&delTest{})

	dataArr := make([]delTest, 0, 10)
	dataArr = append(dataArr, delTest{UserName: "AAA", Age: 20, Addr: "北京0"},
		delTest{UserName: "BBB", Age: 21, Addr: "北京1"}, delTest{UserName: "CCC", Age: 22, Addr: "北京2"}, delTest{UserName: "DDD", Age: 23, Addr: "北京3"},
		delTest{UserName: "EEE", Age: 24, Addr: "北京4"}, delTest{UserName: "FFF", Age: 25, Addr: "北京5"}, delTest{UserName: "HHH", Age: 26, Addr: "北京6"},
		delTest{UserName: "III", Age: 27, Addr: "北京7"}, delTest{UserName: "JJJ", Age: 28, Addr: "北京8"}, delTest{UserName: "KKK", Age: 29, Addr: "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //创建10条数据
		log.WithFields(log.Fields{"data": value}).Debug("*-*-*-*-*- delData1 insert data finish")
		respones := dbClient.Create(&value)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": value}).Debug("create a record success!")
	}

	delData1 := delTest{Id: 5} //删除id为5的数据
	respones := dbClient.Delete(&delData1)
	affectRow := respones.RowsAffected
	log.WithFields(log.Fields{
		"affectRow": affectRow, //affectRow=1
		"delData1":  delData1,
	}).Info("delete one data finish")

	delData2 := delTest{UserName: "CCC"} //todo 如果结构中没有Id(即主键)字段，则删除所有数据(不会删除UserName为CCC的数据)
	affectRow = dbClient.Delete(&delData2).RowsAffected
	log.WithFields(log.Fields{
		"affectRow": affectRow, //ffectRow=9
		"delData1":  delData2,
	}).Info("delete data without id field")
}

func insData() { //插入数据
	log.Info("============================ insData ============================")
	type insTest struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}

	reCreateTb(&insTest{})

	type dataMap struct {
		name string
		age  int
		addr string
	}
	dataArr := make([]dataMap, 0, 10) //创建一个有10条数据的数组
	dataArr = append(dataArr,
		dataMap{"AAA", 20, "北京0"}, dataMap{"BBB", 21, "北京1"}, dataMap{"CCC", 22, "北京2"},
		dataMap{"DDD", 23, "北京3"}, dataMap{"EEE", 24, "北京4"}, dataMap{"FFF", 25, "北京5"},
		dataMap{"HHH", 26, "北京6"}, dataMap{"III", 27, "北京7"}, dataMap{"JJJ", 28, "北京8"},
		dataMap{"KKK", 29, "北京9"},
	)
	log.Infof("dataArr is: %v", dataArr)

	for _, value := range dataArr { //循环依次创建数据(create不能一次插入多条数据)
		log.Infof("**** value: %v, name: %v, age: %v, addr: %v", value, value.name, value.age, value.addr)
		tempData := insTest{
			UserName: value.name,
			Age:      value.age,
			Addr:     value.addr,
		}
		respones := dbClient.Create(&tempData)
		affectRow := respones.RowsAffected
		log.WithFields(log.Fields{"affectRow": affectRow, "data": tempData}).Info("create a record success!")
	}
}

func reCreateTb(t interface{}) {
	dbClient.DropTableIfExists(t) //删除表
	log.Info("reCreateTb DropTableIfExists finish!")

	dbClient.AutoMigrate(t) //创建表
	log.Info("reCreateTb CreateTable finish!")
}

//使用 Delete() 方法进行单条删除时，在删除前必须确保传入主键必须存在否则会 删除整张表
func delete1() {
	log.Info("============================ delete1 ============================")
	data := UserInfo1{
		UserName: "AAA",
		Age:      30,
		Addr:     "地址1",
	}
	respones := dbClient.Create(&data)
	affectRow := respones.RowsAffected
	log.WithFields(log.Fields{"affectRow": affectRow, "data": data}).Info("delete1 create a record success!")

	delData := UserInfo1{//删除id为1的数据
		Id: data.Id,
	}
	respones = dbClient.Delete(&delData)
	affectRow = respones.RowsAffected
	log.WithFields(log.Fields{"affectRow": affectRow, "data": delData}).Info("delete1 a record success!")
}

func create1() {
	data := UserInfo{
		UserName: "AAA",
		Age:      30,
		Addr:     "河北",
	}
	respones := dbClient.Omit("Age", "Addr").Create(&data) //Omit插入数据时会忽略 Age Addr 两个字段
	affectRow := respones.RowsAffected
	log.WithFields(log.Fields{"affectRow": affectRow, "data": data}).Info("create1 a record success!")
}

func create(data *UserInfo) {
	respones := dbClient.Create(&data)
	affectRow := respones.RowsAffected
	log.WithFields(log.Fields{"affectRow": affectRow, "data": data}).Info("create a record success!")
}

func execSql() {
	log.Info("==================================================================")
	tbName := "myUserInfo2"
	field1 := "user_name"
	field2 := "age"
	field3 := "addr"
	sqlStr := fmt.Sprintf(`insert into %v (%v, %v, %v) values (?,?,?), (?,?,?)`, tbName, field1, field2, field3)
	log.Infof("execSql sqlStr: %v", sqlStr)

	args := make([]interface{}, 0, 6)
	args = append(args, "aaa", "25", "beijing", "bbb", "30", "shanghai")
	sqlResponse, err := dbClient.DB().Exec(sqlStr, args...)
	if err != nil {
		log.Errorf("sql exec err: %v", err.Error())
		return
	}
	lastInsertId, err := sqlResponse.LastInsertId()
	affectCount, err := sqlResponse.RowsAffected()
	if err != nil {
		log.Errorf("parse sqlResponse err: %v", err.Error())
		return
	}
	log.WithFields(log.Fields{
		"lastInsertId": lastInsertId,
		"affectCount":  affectCount,
	}).Info("sql cmd exec result")
}
