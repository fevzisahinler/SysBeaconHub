package main 

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Hata: .env dosyası yüklenemedi: %v", err)
		return
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE_NAME")


	// Connect RabbitMQ
	conn, err:= amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to establish RabbitMQ connection: %v", err)
		return
	}
	defer conn.Close()

	// Open Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
		return
	}
	defer ch.Close()

	// Subscribe to Queue
	q, err := ch.QueueDeclare(
		rabbitMQQueue,
		false, // Refers to the durability of the queue. When set to False, the queue exists only as long as the RabbitMQ server is running. When set to True, the queue is persistent and is maintained even when the RabbitMQ server restarts.
		false, // Specifies that this queue belongs only to the originating connection.
		false, // Indicates that this queue belongs only to the connection that created it and is not automatically deleted only when the connection is closed.
		false, // It is used to add special extra settings to the arguments of this queue. This example does not use any extra arguments
		nil,  // "nil" indicates that you do not specify extra settings here either.
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to queue: %v", err)
		return
	}
	// Receive messages
	msgs, err := ch.Consume(
		q.Name, // Specifies which queue to retrieve messages from. The name of the queue created above is used.
		"",     // Specifies the client name. If this is left blank, the client name is automatically assigned.
		true,   // Specifies whether messages are automatically deleted from the queue. When set to True, the message is automatically deleted after it is retrieved from the queue.
		false,  // Specifies whether you want to use the ack (acknowledge) function, which ensures that when messages are forwarded to another client, that client also receives the same message.
		false,	// Specifies whether to use the automatically reject function when receiving messages. It is set to false because messages are processed manually.
		false,  // For extra settings nil is used.
		nil,	// "nil" indicates that you do not specify extra settings here either.
	)
	if err != nil {
		log.Fatalf("Error receiving messages: %v", err)
		return
	}

	
	fmt.Printf("Subscribed to the RabbitMQ queue. Listening to messages has started. Press CTRL+C to exit\n")
	for d := range msgs {
		fmt.Printf("Channel: %s, Message: %s\n", q.Name, d.Body)
	}

}