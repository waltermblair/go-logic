package brain

func RunDemo(body MessageBody, rabbit RabbitClient) (err error){

//	build message

//  publish message
err = rabbit.Publish(body)

return err

}

