package logic

import (
	"fmt"
	"log"
	"strconv"
	"errors"
)

var config Config

func applyConfig(cfg Config) (err error) {

	config = cfg
	return err

}

// todo - multiple inputs1
func applyFunction(body MessageBody) bool {

	if config.Status == "down" {
		log.Fatal(errors.New("this component is down"))
	}

	input := body.Input[0]
	var output bool

	switch fn := config.Function; fn {
	case "buffer":
		output = input
	case "not":
		output = !input
	}

	return output

}

func buildMessage(body MessageBody) MessageBody {

	output := applyFunction(body)

	return MessageBody{
		Configs: []Config{},
		Input: []bool{output},
	}
}

func Process(body MessageBody, rabbit RabbitClient) (err error){

	fmt.Println("number of configs in message: ", len(body.Configs))
	fmt.Println("number of nextKeys: ", len(config.NextKeys))

	if len(body.Configs) == 1 {
		fmt.Println("received config message, applying config...")
		applyConfig(body.Configs[0])
	}

	//	build and publish each message
	for _, nextQueue := range config.NextKeys {

		msg := buildMessage(body)

		fmt.Println("sending this message: ", msg, "to queue: ", nextQueue)

		err = rabbit.Publish(msg, strconv.Itoa(nextQueue))

	}

	return err

}

