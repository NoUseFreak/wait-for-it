package main

import (
	"github.com/NoUseFreak/wait-for-it/plugin"
	"github.com/go-redis/redis"
	"os"
	"time"
)

func main() {
	parameters := plugin.ParseArguments()

	for {
		client := redis.NewClient(&redis.Options{
			Addr:     parameters["hostname"] + ":" + parameters["port"],
			Password: parameters["password"],
			DB:       0,
		})

		_, err := client.Ping().Result()
		if err == nil {
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
}
