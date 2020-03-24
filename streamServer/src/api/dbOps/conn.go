package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	//TODO many service try to connect to mysql servre, if some of them are failed, we need to identify which client has been terminated and try to figure out the reason.
	dbConn, err = sql.Open("mysql", "root:anderson@tcp(127.0.0.1:3306)/stream_server?charset=utf8")
	if err != nil {
		panic("wrong here:" + err.Error())
	}
}
