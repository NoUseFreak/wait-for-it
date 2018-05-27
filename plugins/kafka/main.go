package main

import (
	"github.com/Shopify/sarama"
	"github.com/NoUseFreak/wait-for-it/plugin"
	"os"
	"time"
	"strings"
)

func main() {
	parameters := plugin.ParseArguments()
	for {
		success := DoKafkaTest(parameters)
		if success {
			os.Exit(0)
		}
		time.Sleep(1 * time.Second)
	}
}

func DoKafkaTest(parameters map[string]string) bool {
	defer func() {
		recover()
	}()
	brokers := strings.Split(parameters["brokers"], ",")

	config := sarama.NewConfig()
	consumer, _ := sarama.NewConsumer(brokers, config)
	defer consumer.Close()

	return true
}