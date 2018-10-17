package main

import (
	"fmt"
	. "github.com/waltermblair/logic/logic"
	"log"
	"os"
)

// Creates rabbit client with queue specified by env variable. Creates processor and runs services.
func main() {

	log.Println("RABBIT HOST: ", os.Getenv("RABBIT_HOST"))
	log.Println("THIS QUEUE: " , os.Getenv("THIS_QUEUE"))

	rabbit := NewRabbitClient(fmt.Sprintf("amqp://guest:guest@%s:5672/", os.Getenv("RABBIT_HOST")), os.Getenv("THIS_QUEUE"))
	processor := NewProcessor()
	go rabbit.RunConsumer(processor)
	RunAPI()

}