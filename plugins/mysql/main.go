package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"../../plugin"
	"database/sql"
	"os"
	"time"
)

func main() {
	parameters := plugin.ParseArguments()

	dsn := fmt.Sprintf(
		"%s:%s@%s/%s",
		parameters["username"],
		parameters["password"],
		parameters["hostname"],
		parameters["database"],
	)

	for {
		db, _ := sql.Open("mysql", dsn)
		defer db.Close()
		err := db.Ping()
		if err == nil {
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
}
