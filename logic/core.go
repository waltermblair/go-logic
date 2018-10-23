package logic

import (
	"errors"
	"log"
	"strconv"
)

type Processor interface {
	GetConfig() Config
	ApplyConfig(Config) (err error)
	ApplyFunction(MessageBody) bool
	BuildMessage(MessageBody) MessageBody
	Process(MessageBody, RabbitClient) (err error)
}

type ProcessorImpl struct {
	config		Config
}

func NewProcessor() Processor {
	p := ProcessorImpl{
		Config{},
	}

	return &p
}

func (p *ProcessorImpl) GetConfig() (Config) {
	return p.config
}

func (p *ProcessorImpl) ApplyConfig(cfg Config) (err error) {

	p.config = cfg

	return err

}

func (p *ProcessorImpl) ApplyFunction(body MessageBody) bool {

	input := body.Input[0]
	var output bool

	switch fn := p.config.Function; fn {
	case "buffer":
		output = input
	case "not":
		output = !input
	}

	return output

}

func (p *ProcessorImpl) BuildMessage(body MessageBody) MessageBody {

	if p.config.Status == "down" {
		log.Fatal(errors.New("this component is down"))
	}

	output := p.ApplyFunction(body)

	return MessageBody{
		Configs: []Config{},
		Input: []bool{output},
	}
}

func (p *ProcessorImpl) Process(body MessageBody, rabbit RabbitClient) (err error){

	log.Println("number of configs in message: ", len(body.Configs))

	//  if there's a config in the message, apply it
	configs := body.Configs

	if len(configs) == 1 {
		p.ApplyConfig(configs[0])
	}

	//	build and publish one message for each downstream component
	for _, nextQueue := range p.config.NextKeys {

		msg := p.BuildMessage(body)

		log.Println("sending this message: ", msg, "to queue: ", nextQueue)

		err = rabbit.Publish(msg, strconv.Itoa(nextQueue))

	}

	return err

}

