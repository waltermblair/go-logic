package main

import (
	. "github.com/waltermblair/logic/logic"
)

func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", "1") // todo - replace with env
	go rabbit.RunConsumer()
	RunAPI(rabbit)

}