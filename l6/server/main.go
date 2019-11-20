package main

import (
	config "binTest/rabbitmqTest/t1/l6/conf"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Qos(
		config.SERVERINSTANCESCNT,
		0,
		false,
	)

	forever := make(chan bool)

	for routine := 0; routine < config.SERVERINSTANCESCNT; routine++ {
		go func(routineNum int) {

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
				"",
				false, // auto ack
				false,
				false,
				false,
				nil,
			)

			for msg := range msgs {
				log.Printf("In %d start consuming message: %s\n", routineNum, msg.Body)

				bookName := queryBookID(string(msg.Body))

				err = ch.Publish(
					"",
					msg.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: msg.CorrelationId,
						Body:          []byte(bookName),
					})

				if err != nil {
					fmt.Println("Failed to reply msg to client")
				} else {
					fmt.Println("Response to client:", bookName)
				}
				msg.Ack(false)
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

func queryBookID(bookID string) string {
	bookName := "QUERIED_" + bookID
	time.Sleep(time.Duration(rand.Intn(9)) * time.Second)

	return bookName
}
