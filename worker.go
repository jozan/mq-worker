package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	conn, ch, q := queueConnect()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	f, err := os.OpenFile("./backup.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	failOnError(err, "Failed to create file")
	defer f.Close()

	w := csv.NewWriter(f)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
			d.Ack(false)
			writeToFile(w, d.Body)
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func writeToFile(w *csv.Writer, data []byte) {
	str := []string{string(data)} // Casting like a true magicia
	err := w.Write(str)           // Write to buffer
	failOnError(err, "Error writing to csv")

	w.Flush() // Ensure all buffered operations are applied to disk the underlying writer

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
