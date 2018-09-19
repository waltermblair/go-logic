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
	Publish(MessageBody) error
	InitRabbit()
}

type RabbitClientImpl struct {
	URL 	string
	que 	*cony.Queue
	exc    	cony.Exchange
	bnd 	cony.Binding
	cli     *cony.Client
	cns     *cony.Consumer
	pbl     *cony.Publisher
}

func NewRabbitClient(url string) RabbitClient {

	r := RabbitClientImpl{URL: url}
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
		cony.DeclareQueue(r.que),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.bnd),
	})

	// Declare and register a consumer
	cns := cony.NewConsumer(r.que)

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

func (r *RabbitClientImpl) Publish(body MessageBody) error {

	cli := cony.NewClient(
		cony.URL(r.URL),
		cony.Backoff(cony.DefaultBackoff),
	)

	cli.Declare([]cony.Declaration{
		cony.DeclareQueue(r.que),
		cony.DeclareExchange(r.exc),
		cony.DeclareBinding(r.bnd),
	})

	pbl := cony.NewPublisher(r.exc.Name, "pubSub")
	cli.Publish(pbl)

	go func() {
		for cli.Loop() {
			select {
			case err := <-cli.Errors():
				fmt.Println(err)
			}
		}
	}()

	fmt.Println("Client publishing to exchange", r.exc.Name)

	bytes, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("Error unmarshaling MessageBody: %v\n", err)
	}

	go func() {
		fmt.Println("Publishing body: ", body)

		err = pbl.Publish(amqp.Publishing{
			Body: bytes,
		})
		if err != nil {
			fmt.Printf("Client publish error: %v\n", err)
		} else {
			fmt.Println("Done!")
		}
	}()

	return err

}

func (r *RabbitClientImpl) InitRabbit() {

	// Declarations
	// The queue name will be supplied by the AMQP server
	r.que = &cony.Queue{
		AutoDelete: false,
		Name:       "myQueue",
		Durable:	true,
	}
	r.exc = cony.Exchange{
		Name:       "myExc",
		Kind:       "topic",
		AutoDelete: false,
		Durable:	true,
	}
	r.bnd = cony.Binding{
		Queue:    r.que,
		Exchange: r.exc,
		Key:      "pubSub",
	}

}