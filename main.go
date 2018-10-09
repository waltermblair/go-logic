package main

import (
	. "github.com/waltermblair/logic/logic"
	"os"
)

// Creates rabbit client with queue specified by env variable. Creates processor and runs services.
func main() {

	rabbit := NewRabbitClient("amqp://guest:guest@localhost:5672/", os.Getenv("THIS_QUEUE"))
	processor := NewProcessor()
	go rabbit.RunConsumer(processor)
	RunAPI()

}