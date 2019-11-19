package main

import (
	config "binTest/rabbitmqTest/t1/l4/conf"
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
		config.EXCHANGENAME, //exchange name
		"direct",            //exchange kind
		true,                //durable
		false,               //autodelete
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	if len(os.Args) < 3 {
		fmt.Println("Arguments error(ex:producer info/debug/warn/error msg1 msg2 msg3")
		return
	}

	routingKey := os.Args[1]

	validKey := false
	for _, item := range config.RoutingKeys {
		if routingKey == item {
			validKey = true
			break
		}
	}

	if validKey == false {
		fmt.Println("Arguments error, no valid routing key specified.(ex:producer info/debug/warn/error msg1 msg2 msg3")
		return
	}
	msgs := os.Args[2:]

	msgNum := len(msgs)

	for cnt := 0; cnt < msgNum; cnt++ {
		msgBody := msgs[cnt]
		err = ch.Publish(
			config.EXCHANGENAME, //exchange
			routingKey,          //routing key
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
