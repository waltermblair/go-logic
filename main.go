package main

import (
	. "github.com/waltermblair/brain/brain"
)

func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", "0", "1") // todo - replace with env
	go rabbit.RunConsumer()
	RunAPI(rabbit)

}