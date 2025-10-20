package events

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewConsumer(rabbitmqURL string) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{conn, ch}, nil

}
func (c *Consumer) ConsumeTransactionCreated() error {
	q, err := c.ch.QueueDeclare(
		"transaction_created",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	msgs, err := c.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		var event TransactionCreatedEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}
		log.Printf("Event received: %v", event)
		// Здесь можно добавить логику обработки:
		// - Проверить бюджет
		// - Отправить уведомление
		// - Обновить аналитику
	}
	return nil
}

func (c *Consumer) BalanceAlertConsumer() error {
	q, err := c.ch.QueueDeclare(
		"balance_alert",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		var event LowBalanceEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}
		log.Printf("Event received: %v", event)

	}
	return nil
}
