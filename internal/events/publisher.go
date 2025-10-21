package events

import (
	"encoding/json"

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
func (p *Publisher) PublishTransactionCreated(event TransactionCreatedEvent) error {
	q, err := p.ch.QueueDeclare(
		"transaction_created", true,
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
func (p *Publisher) PublishBudgetExceeded(event BudgetExceededEvent) error {
	q, err := p.ch.QueueDeclare(
		"budget_exceeded",
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

func (p *Publisher) PublishLowBalance(event LowBalanceEvent) error {
	q, err := p.ch.QueueDeclare(
		"low_balance",
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
func (p *Publisher) PublishBudgetWarning(event BudgetWarningEvent) error {
	q, err := p.ch.QueueDeclare(
		"budget_warning",
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
func (p *Publisher) PublishNotification(event NotificationEvent) error {
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
