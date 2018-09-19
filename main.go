package main

import (
	. "github.com/waltermblair/brain/brain"
)

func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/")
	go rabbit.RunConsumer()
	RunAPI(rabbit)

}