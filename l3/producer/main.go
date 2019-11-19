package main

import (
	config "binTest/rabbitmqTest/t1/l3/conf"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"fanout_exchange", //exchange name
		"fanout",          //exchange kind
		true,              //durable
		false,             //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	msgs := os.Args[1:]
	msgNum := len(msgs)

	for cnt := 0; cnt < msgNum; cnt++ {
		msgBody := msgs[cnt]
		err = ch.Publish(
			"fanout_exchange", //exchange
			"",                //routing key
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msgBody),
			})

		log.Printf(" [x] Sent %s", msgBody)
	}
	failOnError(err, "Failed to publish a message")

}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
