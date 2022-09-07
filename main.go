package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type Alert struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func main() {

	log.Println("Consumer Application")
	conn, err := amqp.Dial("amqps://tmwnunoc:EspITp82kgC9D7sVtaeJ4Cn2lF0QDnU4@shark.rmq.cloudamqp.com/tmwnunoc")
	if err != nil {
		log.Fatalf("Error connecting to rabbitMq %v", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to rabbitmq")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel %v", err)
	}

	defer ch.Close()

	msgs, err := ch.Consume(
		"Arima",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Error publishing queue %v", err)
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var alert Alert
			if err := json.Unmarshal(d.Body, &alert); err != nil {
				log.Fatalf("%v", err)
			}
			/*
				res, err := CreateTweet(alert.Message)
				if err != nil {
					log.Println(err)
				}
				fmt.Printf("Twitter Response: %s \n", res.Tweet.ID)
			*/
			SendMails(alert.Topic, alert.Message, []string{"salemododa2@gmail.com", "robtyler0701@gmail.com"})
			fmt.Printf("Recieved Message: %s \n", alert)
		}
	}()

	<-forever

}
