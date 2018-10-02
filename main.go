package main

import (
	. "github.com/waltermblair/logic/logic"
)

func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", "1") // TODO - replace with env from docker-compose
	processor := NewProcessor()
	go rabbit.RunConsumer(processor)
	RunAPI()

}