package main

import (
	"github.com/gocql/gocql"
	"github.com/NoUseFreak/wait-for-it/plugin"
	"os"
	"time"
	"strings"
)

func main() {
	parameters := plugin.ParseArguments()

	hosts := strings.Split(parameters["hosts"], ",")

	cluster := gocql.NewCluster(hosts...)
	if keyspace, ok := parameters["keyspace"]; ok {
		cluster.Keyspace = keyspace
	}

	for {
		session, err := cluster.CreateSession()
		defer session.Close()
		if err == nil {
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
}
