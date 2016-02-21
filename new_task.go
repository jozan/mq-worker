package main

import (
	amqp "github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, ch, q := queueConnect()
	defer conn.Close()
	defer ch.Close()

	body := bodyFrom(os.Args)
	err := ch.Publish(
		"tussi_exchange", // exchange
		q.Name,           // routing key
		true,             // mandatory
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
