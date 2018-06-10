package main

import (
	"github.com/NoUseFreak/wait-for-it/plugin"
	"gopkg.in/mgo.v2"
	"os"
	"time"
	"fmt"
)

func main() {
	parameters := plugin.ParseArguments()

	for {
		err := DoTest(parameters)
		if err == nil {
			os.Exit(0)
		}
		fmt.Println(err)
		time.Sleep(1 * time.Second)
	}
}


func DoTest(parameters map[string]string) error {
	sess, err := mgo.DialWithTimeout(fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		parameters["username"],
		parameters["password"],
		parameters["host"],
		parameters["port"],
	), 1 * time.Second)
	if err != nil {
		return err
	}
	defer sess.Close()
	return sess.Ping()
}