package main

import (
	config "binTest/rabbitmqTest/t1/l3/conf"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)

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

	for routine := 0; routine < config.CONSUMERCNT; routine++ {
		go func(routineNum int) {

			q, err := ch.QueueDeclare(
				"",
				false, //durable
				false,
				true,
				false,
				nil,
			)

			failOnError(err, "Failed to declare a queue")

			err = ch.QueueBind(
				q.Name,
				"",
				"fanout_exchange",
				false,
				nil,
			)
			failOnError(err, "Failed to bind exchange")

			msgs, err := ch.Consume(
				q.Name,
				"",
				true, //Auto Ack
				false,
				false,
				false,
				nil,
			)

			if err != nil {
				log.Fatal(err)
			}

			for msg := range msgs {
				log.Printf("In %d consume a message: %s\n", routineNum, msg.Body)
			}

		}(routine)
	}

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
