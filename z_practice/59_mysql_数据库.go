package main

import (
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"database/sql"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/z_go_test?charset=utf8mb4")
	checkErr(err)
	fmt.Println("open mysql success")
	fmt.Println("======================== insert test ========================")
	id := insertDb(db)

	//fmt.Println("======================== update test ========================")
	//updateDb(db, id)

	fmt.Println("======================== select test ========================")
	selectDb(db)

	fmt.Println("======================== delete test ========================")
	deleteDb(db, id)

	db.Close()
}

func deleteDb(db *sql.DB, id int64) {
	stmt, err := db.Prepare("delete from userinfo where uid = ?")
	checkErr(err)
	fmt.Println("prepare sql finish")

	res, err := stmt.Exec(id)
	checkErr(err)
	fmt.Println("exec sql finish")

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("delete affect:", affect)
}

func selectDb(db *sql.DB) {
	rows, err := db.Query("select * from userinfo")
	checkErr(err)
	fmt.Println("query sql finish")
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Printf("uid: %v, username: %v, department: %v, created: %v\n", uid, username, department, created)
	}
}

func updateDb(db *sql.DB, id int64) {
	stmt, err := db.Prepare("update userinfo set username = ? where uid = ?")
	checkErr(err)
	fmt.Println("prepare sql finish")

	res, err := stmt.Exec("AAAAA", id)
	checkErr(err)
	fmt.Println("exec sql finish")

	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println("update affect:", affect)
}

func insertDb(db *sql.DB) int64 { //插入数据
	stmt, err := db.Prepare("insert into userinfo(username, departname, created) values (?,?,?)")
	checkErr(err)
	fmt.Println("prepare sql finish")

	res, err := stmt.Exec("abcd", "研发部", time.Now().Format("2006-01-02 15:04:05"))
	checkErr(err)
	fmt.Println("exec sql finish")

	id, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("last insert id:", id)
	return id
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		panic(err)
	}
}
