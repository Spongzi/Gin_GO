package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // init()
)

type user struct {
	id   int
	age  int
	name string
}

var db *sql.DB

func initMySql() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/demo"
	// 去初始化全局的db对象，而不是新声明一个
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接
	err = db.Ping()
	if err != nil {
		fmt.Println("连接失败了", err)
		return err
	}
	// 最大连接数
	db.SetMaxOpenConns(200)
	// 最大空闲连接数
	db.SetMaxIdleConns(10)
	return nil
}

// 查询单条数据
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id = ?"
	var u user
	// 非常重要，确保QueryRow之后调用Scan方法，否则持有的数据库链接不会被释放
	row := db.QueryRow(sqlStr, 1)
	err := row.Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Println("scan failed!", err)
		return
	}
	fmt.Printf("获取的数据: id = %v, name = %v, age = %v\n", u.id, u.name, u.age)
}

// 查找多条数据
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?;"
	rows, err := db.Query(sqlStr, 1)
	if err != nil {
		fmt.Println("select failed!", err)
	}
	// 一定要注意关闭！！！！
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("close failed!", err)
		}
	}(rows)
	// 循环读取查询到的内容
	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Println("err failed", err)
		}
		fmt.Printf("获取的数据: id = %v, name = %v, age = %v\n", u.id, u.name, u.age)
	}
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?, ?)"
	exec, err := db.Exec(sqlStr, "苏旭4", 30)
	if err != nil {
		fmt.Println("exec sql failed", err)
		return
	}
	theID, err := exec.LastInsertId()
	if err != nil {
		fmt.Println("get lastInsertId failed", err)
		return
	}
	fmt.Printf("insert success, the id is %v\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age = ? where id = ?"
	result, err := db.Exec(sqlStr, 10, 1)
	if err != nil {
		fmt.Println("exec failed!", err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("affected failed!", err)
		return
	}
	fmt.Printf("update success, the update id = %v\n", rowsAffected)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	result, err := db.Exec(sqlStr, 7)
	if err != nil {
		fmt.Println("exec failed!", err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	fmt.Printf("delete success, delete id = %v\n", rowsAffected)
}

func main() {
	err := initMySql()
	if err != nil {
		fmt.Println(err)
	}
	// Close 释放掉数据库连接相关的资源
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	// 单个查询
	queryRowDemo()
	// 插入数据
	insertRowDemo()
	// 更新数据
	updateRowDemo()
	// 删除数据
	deleteRowDemo()
	// 多个查询
	queryMultiRowDemo()
}

// 10个视频！！！
