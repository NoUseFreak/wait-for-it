package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"os"
	"strings"
	"time"
)

func main() {
	parameters := map[string]string{}
	for _, pair := range strings.Split(os.Args[1], ",") {
		pair := strings.Split(pair, "=")
		parameters[pair[0]] = pair[1]
	}

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