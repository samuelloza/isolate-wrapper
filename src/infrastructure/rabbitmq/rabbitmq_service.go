package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/samuelloza/isolate-wrapper/src/domain"
	"github.com/streadway/amqp"
)

type RabbitService struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitService(amqpURL string) (*RabbitService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitService{Conn: conn, Channel: ch}, nil
}

type EvaluationMessage struct {
	Input      domain.EvaluationInput `json:"input"`
	ReplyQueue string                 `json:"reply_queue"`
}

func (r *RabbitService) Listen(queueName string, handler func(domain.EvaluationInput) error) error {
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
