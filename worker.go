package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	failOnErrorWorker(err, "can not connect rabbitmq")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnErrorWorker(err, "fail to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bq",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErrorWorker(err, "fail to declare a queue")
	msg, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErrorWorker(err, "fail to register consume")
	forver := make(chan bool)
	go func() {
		for d := range msg {
			log.Printf("Recived a message %s", d.Body)
			d.Ack(false)
		}
	}()
	log.Printf("[*] Watting for message. To exist press CTRL+C")
	<-forver
	//forver 的用途會讓程式一直保持在go routine的狀態而不會關閉，類似for select
}

func failOnErrorWorker(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", err, msg)
	}
}
