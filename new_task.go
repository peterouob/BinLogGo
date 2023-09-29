package main

import (
	"github.com/peterouob/BinLogGo/binlog"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	go binlog.StartbinLog()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "fail to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bq",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "fail to declare a queue")
	//var body interface{} = <-binlog.BinLogDataChan

	//var msg = []byte(fmt.Sprintf("%v", body))

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			//DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(<-binlog.BinLogChan),
		})
	failOnError(err, "fail to publish a message ")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", err, msg)
	}
}
