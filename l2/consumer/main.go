package main

import (
	config "binTest/rabbitmqTest/t1/l2/conf"
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
				true, //durable
				false,
				false,
				false,
				nil,
			)

			failOnError(err, "Failed to declare a queue")

			err = ch.Qos(
				1,     // prefetch count
				0,     // prefetch size
				false, // global
			)

			failOnError(err, "Failed to set QoS")

			msgs, err := ch.Consume(
				q.Name,
				"MsgWorkConsumer",
				false, //Auto Ack
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
				log.Printf("Done")
				msg.Ack(false) //Ack
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
