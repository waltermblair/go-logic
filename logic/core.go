package logic

import (
	"errors"
	"log"
	"strconv"
)

type Processor interface {
	GetConfig() Config
	GetOutput() bool
	GetNumReceived() int
	ApplyConfig(Config) (err error)
	ApplyFunction(MessageBody)
	BuildMessage() MessageBody
	Process(MessageBody, RabbitClient) (err error)
}

type ProcessorImpl struct {
	config		Config
	output      bool
	numReceived int
}

func NewProcessor() Processor {
	p := ProcessorImpl{
		Config{}, true, 0,
	}

	return &p
}

func (p *ProcessorImpl) resetInputs() {
	p.output = true
	p.numReceived = 0
}

func (p *ProcessorImpl) GetConfig() (Config) {
	return p.config
}

func (p *ProcessorImpl) GetOutput() (bool) {
	return p.output
}

func (p *ProcessorImpl) GetNumReceived() (int) {
	return p.numReceived
}

func (p *ProcessorImpl) ApplyConfig(cfg Config) (err error) {

	p.config = cfg

	return err

}

func (p *ProcessorImpl) ApplyFunction(body MessageBody) {

	input := body.Input[0]

	switch fn := p.config.Function; fn {
	case "buffer":
		p.output = input
	case "not":
		p.output = !input
	case "and":
		p.output = p.output && input
	case "or":
		if p.numReceived == 0 {
			p.output = input
		} else {
			p.output = p.output || input
		}
	}

}

func (p *ProcessorImpl) BuildMessage() MessageBody {

	if p.config.Status == "down" {
		log.Fatal(errors.New("this component is down"))
	}

	return MessageBody{
		Configs: []Config{},
		Input: []bool{p.output},
	}
}

func (p *ProcessorImpl) Process(body MessageBody, rabbit RabbitClient) (err error){

	log.Println("number of configs in message: ", len(body.Configs))

	//  if there's a config in the message, apply it
	configs := body.Configs

	if len(configs) == 1 {
		p.ApplyConfig(configs[0])
	}

	if body.Input != nil && p.numReceived < p.config.NumInputs {
		p.ApplyFunction(body)
		p.numReceived += 1
	}

	if p.numReceived == p.config.NumInputs {
		//	build and publish one message for each downstream component
		for _, nextQueue := range p.config.NextKeys {
			msg := p.BuildMessage()
			log.Println("sending this message: ", msg, "to queue: ", nextQueue)
			err = rabbit.Publish(msg, strconv.Itoa(nextQueue))
		}
		p.resetInputs()
	}

	return err

}

