package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var DB *sql.DB
var NotFounInfo = errors.New("No Data")

const (
	userName = "root"
	passwd   = "123456"
	address  = "127.0.0.1:3306"
)

type UserInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 连接mysql
func ConnectDB() (dbErr error) {
	dbPath := fmt.Sprintf("%s:%s@tcp(%s)/test?parseTime=true", userName, passwd, address)
	DB, dbErr = sql.Open("mysql", dbPath)
	if dbErr != nil {
		return errors.Wrapf(dbErr, "connect %s err", dbPath)
	}
	if dbErr = DB.Ping(); dbErr != nil {
		return errors.Wrapf(dbErr, "ping %s err", dbPath)
	}
	return nil
}

// 查询
func GetUserByName(username string) (data *UserInfo, err error) {
	s := "SELECT * FROM user_info WHERE name = ? "
	err = DB.QueryRow(s, username).Scan(data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NotFounInfo
		}
		return nil, errors.Wrapf(err, "get user err:%s,%s", s, username)
	}
	return
}

func main() {
	err := ConnectDB()
	if err != nil {
		fmt.Printf("original error:%+v\n", errors.Cause(err))
		fmt.Printf("stack:\n%+v\n", err)
		return
	}
	u, err := GetUserByName("mujibin")
	if err != nil {
		if errors.Is(err, NotFounInfo) {
			fmt.Println("no")
			return
		}
		fmt.Printf("original error:%+v\n", errors.Cause(err))
		fmt.Printf("stack:\n%+v\n", err)
		return
	}
	fmt.Println(u.ID, u.Age)
}
