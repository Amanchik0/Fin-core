package events

import (
	"encoding/json"
	"justTest/internal/models/events"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewPublisher(rabbitmqURl string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitmqURl)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{conn, ch}, nil

}

func (p *Publisher) Close() error {
	if p.ch != nil {
		p.ch.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
	return nil
}
func (p *Publisher) PublishTransactionCreated(event events.TransactionCreatedEvent) error {
	q, err := p.ch.QueueDeclare(
		"transaction_created", false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

// TODO почекать вопрос : Сделать все очереди durable и сообщения persistent
// TODO balance_alert
func (p *Publisher) PublishBudgetExceeded(event events.BudgetExceededEvent) error {
	q, err := p.ch.QueueDeclare(
		"budget_exceeded",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func (p *Publisher) PublishLowBalance(event events.LowBalanceEvent) error {
	q, err := p.ch.QueueDeclare(
		"low_balance",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
func (p *Publisher) PublishBudgetWarning(event events.BudgetWarningEvent) error {
	q, err := p.ch.QueueDeclare(
		"budget_warning",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
func (p *Publisher) PublishNotification(event events.NotificationEvent) error {
	q, err := p.ch.QueueDeclare(
		"notification",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
