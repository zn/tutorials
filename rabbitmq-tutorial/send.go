package main

import(
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s", msg, err)
	}
}

func main(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil)

	failOnError(err, "Failed to declare a queue")

	body := "Hello world!"

	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		Body: []byte(body),
		ContentType: "text/plain",
	})
	failOnError(err, "Failed to publish a message")
}