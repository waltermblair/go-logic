package brain

import (
	"encoding/json"
	"fmt"
	"github.com/assembla/cony"
	"github.com/streadway/amqp"
	"log"
)

type RabbitClient interface {
	RunConsumer()
	Publish(MessageBody, string) error
	InitRabbit()
}

type RabbitClientImpl struct {
	URL       string
	exc       cony.Exchange
	cli       *cony.Client
	cns       *cony.Consumer
	cnsQue    *cony.Queue
	cnsBnd    cony.Binding
	pbl       *cony.Publisher
	pblQue    *cony.Queue
	pblBnd    cony.Binding
	thisQueue string
	nextQueue string
}

func NewRabbitClient(url string, thisQueue string, nextQueue string) RabbitClient {

	r := RabbitClientImpl{
		URL:       url,
		thisQueue: thisQueue,
		nextQueue: nextQueue,
	}
	r.InitRabbit()

	fmt.Println("Initialized rabbit client at ", r.URL)

	return &r

}

func (r *RabbitClientImpl) RunConsumer() {

	cli := cony.NewClient(
		cony.URL(r.URL),
		cony.Backoff(cony.DefaultBackoff),
	)

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(r.cnsQue),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.cnsBnd),
	})

	// Declare and register a consumer
	cns := cony.NewConsumer(r.cnsQue)

	cli.Consume(cns)

	for cli.Loop() {
		select {
		case msg := <-cns.Deliveries():
			log.Printf("Received body: %q\n", msg.Body)
			msg.Ack(false)
		case err := <-cns.Errors():
			fmt.Printf("Consumer error: %v\n", err)
		case err := <-cli.Errors():
			fmt.Printf("Client error: %v\n", err)
		}
	}

}

func (r *RabbitClientImpl) Publish(body MessageBody, nextQueue string) error {

	cli := cony.NewClient(
		cony.URL(r.URL),
		cony.Backoff(cony.DefaultBackoff),
	)

	r.pblQue = &cony.Queue{
		AutoDelete: false,
		Name:       nextQueue,
		Durable:	true,
	}
	r.pblBnd = cony.Binding{
		Queue:    r.pblQue,
		Exchange: r.exc,
		Key:      nextQueue,
	}

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(r.pblQue),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.pblBnd),
	})

	pbl := cony.NewPublisher(r.exc.Name, nextQueue)
	cli.Publish(pbl)

	go func() {
		for cli.Loop() {
			select {
			case err := <-cli.Errors():
				fmt.Println(err)
			}
		}
	}()

	bytes, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error unmarshaling MessageBody: %v\n", err)
	}

	go func() {
		err = pbl.Publish(amqp.Publishing{
			Body: bytes,
		})
		if err != nil {
			fmt.Printf("Client publish error: %v\n", err)
		}
	}()

	return err

}

func (r *RabbitClientImpl) InitRabbit() {

	r.exc = cony.Exchange{
		Name:       "myExc",
		Kind:       "topic",
		AutoDelete: false,
		Durable:	true,
	}
	r.cnsQue = &cony.Queue{
		AutoDelete: false,
		Name:       r.thisQueue,
		Durable:	true,
	}
	r.cnsBnd = cony.Binding{
		Queue:    r.cnsQue,
		Exchange: r.exc,
		Key:      r.thisQueue,
	}

}