package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/streadway/amqp"
)

type RabbitService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitService(amqpURL string) (*RabbitService, error) {
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for {
		conn, err = amqp.DialConfig(amqpURL, amqp.Config{
			Heartbeat: 5 * time.Second,
			Locale:    "en_US",
		})
		if err != nil {
			log.Printf("RabbitMQ connection failed: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			log.Printf("Failed to open channel: %v. Retrying in 5s...", err)
			_ = conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	return &RabbitService{Conn: conn, Channel: ch}, nil
}

type EvaluationMessage struct {
	Input      domain.EvaluationInput `json:"input"`
	ReplyQueue string                 `json:"reply_queue"`
}

func (r *RabbitService) Listen(queueName string, handler func(domain.EvaluationInput) error) error {
	err := r.Channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := r.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("[*] Listening on queue: %s", queueName)

	for d := range msgs {
		var msg domain.EvaluationInput

		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("[!] Failed to parse message: %v", err)
			d.Nack(false, false)
			continue
		}

		if err := handler(msg); err != nil {
			log.Printf("[!] Error handling message: %v", err)
			d.Nack(false, false)
			continue
		}

		d.Ack(false)
	}

	return nil
}

func (r *RabbitService) Publish(queueName string, result domain.EvaluationResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to serialize result: %w", err)
	}

	return r.Channel.Publish(
		"", queueName, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}

func (r *RabbitService) Close() {
	_ = r.Channel.Close()
	_ = r.Conn.Close()
}
