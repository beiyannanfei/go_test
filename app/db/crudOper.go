package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

//对数据库的增删改查操作

func do() {
	//execSql()
	//create(&UserInfo{UserName: "abcd", Age: 29, Addr: "天通苑"})
	//create1()
	//delete1()

	insData()
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
