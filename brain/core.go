package brain

import (
	"fmt"
	"strconv"
)

// todo - add function and nextKey from DB
func fetchComponentConfig(config Config) Config {
	return config
}

func selectInput(body MessageBody, config Config) bool {

	var input bool

	if config.Status == "up" {
		input = body.Input[0] // todo - how to select correct input
	}

	return input

}

func buildMessage(body MessageBody, config Config) MessageBody {

	config = fetchComponentConfig(config)
	input := selectInput(body, config)

	return MessageBody{
		Configs: []Config{config},
		Input: []bool{input},
	}
}

func RunDemo(body MessageBody, rabbit RabbitClient) (err error){

	configs := body.Configs
	fmt.Println("number of messages to send: ", len(configs))

	//	build and publish each message
	for _, config := range configs {

		msg := buildMessage(body, config)

		// determine routing key
		nextQueue := strconv.Itoa(config.ID)

		fmt.Println("sending this message: ", msg)

		err = rabbit.Publish(msg, nextQueue)

	}

	return err

}

