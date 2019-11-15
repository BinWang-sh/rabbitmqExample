package main

import (
	config "binTest/rabbitmqTest/t1/l1/conf"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)

	for routine := 0; routine < config.CONSUMERCNT; routine++ {
		go func(routineNum int) {
			ch, err := conn.Channel()
			failOnError(err, "Failed to open a channel")
			defer ch.Close()

			q, err := ch.QueueDeclare(
				config.QUEUENAME,
				false,
				false,
				false,
				false,
				nil,
			)

			failOnError(err, "Failed to declare a queue")

			msgs, err := ch.Consume(
				q.Name,
				"MsgConsumer",
				false,
				false,
				false,
				false,
				nil,
			)

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
