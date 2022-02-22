package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type user struct {
	Id   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

var db *sqlx.DB

// 连接数据库
func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

// 查询单挑数据
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id = ?;"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Println("Get row failed", err)
		return
	}
	fmt.Println(u.Id, u.Name, u.Age)
}

// 查询多条数据
func queryMultiRowDemo() {
	sqlStr := "select * from USER where id > ?;"
	var u []user
	err := db.Select(&u, sqlStr, 1)
	if err != nil {
		return
	}
	for key, val := range u {
		fmt.Println(key, val)
	}
}

// 插入数据
func insertDemo() {
	sqlStr := "insert into USER (name, age) values (?, ?);"
	exec, err := db.Exec(sqlStr, "小白", 99)
	if err != nil {
		fmt.Println("exec failed", err)
		return
	}
	theId, err := exec.LastInsertId()
	if err != nil {
		fmt.Println("get lastInsertId failed!", err)
		return
	}
	fmt.Printf("插入数据的id=%d\n", theId)
}

// 更新数据
func updateDemo() {
	sqlStr := "update user set NAME = ? where id = ?;"
	exec, err := db.Exec(sqlStr, "小黑", 20)
	if err != nil {
		fmt.Println("exec failed", err)
		return
	}
	affectRow, err := exec.RowsAffected()
	if err != nil {
		fmt.Println("affected failed", err)
		return
	}
	fmt.Println("受影响的行数:", affectRow)
}

// 删除数据
func deleteDemo() {
	sqlStr := "delete from USER where id = ?;"
	exec, err := db.Exec(sqlStr, 21)
	if err != nil {
		fmt.Println("delete failed", err)
		return
	}
	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		fmt.Println("rows affected failed", err)
		return
	}
	fmt.Println("受影响的行数:", rowsAffected)
}

// 插入数据
func insertUserDemo() {
	sqlStr := "insert into USER (name, age) values (:name, :age);"
	result, err := db.NamedExec(sqlStr, map[string]interface{}{
		"name": "沙河",
		"age":  20,
	})
	if err != nil {
		fmt.Println("exec failed", err)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	fmt.Println(rowAffected)
}

// 查询数据
func namedQuery() {
	sqlStr := "select * from USER where name = :name ;"
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{
		"name": "苏旭",
	})
	if err != nil {
		fmt.Println("query failed", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Println("scan failed", err)
			return
		}
		fmt.Println(u)
	}
}

// 事物操作
func transactionDeme() (err error) {
	tx, err := db.Beginx()
	if err != nil {
		fmt.Println("begin failed", err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "update user set age = 200 where id = ?"
	exec, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		fmt.Println("exec failed", err)
		return
	}
	row, err := exec.RowsAffected()
	if err != nil {
		return
	}
	if row != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "update user set age = 200 where id = ?"
	exec, err = tx.Exec(sqlStr2, 30)
	if err != nil {
		fmt.Println("exec failed", err)
		return
	}
	row, err = exec.RowsAffected()
	if err != nil {
		return
	}
	if row != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

// BatchInsertUser 批量插入
func BatchInsertUser(users []interface{}) error {
	sqlStr := "insert into user (name, age) values (?, ?, ?);"
	query, args, _ := sqlx.In(sqlStr, users...) // 如果arg实现了driver.Value, sqlx.In()会通过调用Value来展开
	fmt.Println(query)
	fmt.Println(args)
	_, err := db.Exec(query, args)
	return err
}

// BatchInsertUser2 批量插入
func BatchInsertUser2(users []*user) error {
	_, err := db.NamedExec("insert into user (name, age) values (:name, :age);", users)
	return err
}

// QueryByIDs 根据给定的id查询
func QueryByIDs(ids []int) (users []user, err error) {
	query, args, err := sqlx.In("select id, name, age from user where id = (?) ;", ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&users, query, args...)
	return
}

func main() {
	err := initDB()
	if err != nil {
		panic(err)
	}
	// 查询单挑数据
	//queryRowDemo()
	// 查询多条数据
	//queryMultiRowDemo()
	// 插入数据
	//insertDemo()
	//更新数据
	//updateDemo()
	// 删除数据
	//deleteDemo()
	// 插入数据
	//insertUserDemo()
	// 查询
	//namedQuery()
	// 事务
	//transactionDeme()
	// 批量插入
	//BatchInsertUser([])
}
