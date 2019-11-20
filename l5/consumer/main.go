package main

import (
	config "binTest/rabbitmqTest/t1/l5/conf"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

/*


./consumer "#" info.payment.* *.log debug.payment.#

*/

func main() {

	conn, err := amqp.Dial(config.RMQADDR)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	forever := make(chan bool)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		config.EXCHANGENAME, //exchange name
		"topic",             //exchange kind
		true,                //durable
		false,               //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	if len(os.Args) < 2 {
		log.Println(len(os.Args))
		log.Println(`"Arguments error(Example: ./consumer "#" info.payment.* *.log debug.payment.#"`)
		return
	}

	topics := os.Args[1:]
	topicsCnt := len(topics)

	for routing := 0; routing < topicsCnt; routing++ {
		go func(routingNum int) {

			q, err := ch.QueueDeclare(
				"",
				false, //durable
				false, //delete when unused
				true,  //exclusive
				false, //no-wait
				nil,   //arguments
			)

			failOnError(err, "Failed to declare a queue")

			err = ch.QueueBind(
				q.Name,
				topics[routingNum],
				config.EXCHANGENAME,
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

			failOnError(err, "Failed to register a consumer")

			for msg := range msgs {
				log.Printf("In %s consume a message: %s\n", topics[routingNum], msg.Body)
			}

		}(routing)
	}

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
	}
}
