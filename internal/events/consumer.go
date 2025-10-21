package events

import (
	"encoding/json"
	"justTest/internal/services"
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn                *amqp.Connection
	ch                  *amqp.Channel
	notificationService *services.NotificationService
	budgetService       *services.BudgetService // надо добавить сервис
}

func NewConsumer(rabbitmqURL string, notificationService *services.NotificationService,
	budgetService *services.BudgetService) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{conn, ch, notificationService, budgetService}, nil

}
func (c *Consumer) Close() error {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
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
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	log.Println("[Consumer]  Starting to consume transaction_created events...")

	for d := range msgs {
		var event TransactionCreatedEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Printf(" [Consumer] Error unmarshaling event: %v", err)
			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Received transaction: ID=%d, UserID=%s, Amount=%.2f",
			event.TransactionID, event.UserID, event.Amount)
		// надо добавить транзакцию
		if err := c.budgetService.CheckBudgetAfterTransaction(event); err != nil {
			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Transaction processed successfully")
		d.Ack(false)

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

func (c *Consumer) ConsumeBudgetExceeded() error {
	q, err := c.ch.QueueDeclare(
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
	msgs, err := c.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	log.Println("[Consumer]  Starting to consume budget_exceeded events...")

	for d := range msgs {
		var event BudgetExceededEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Printf(" [Consumer] Error unmarshaling event: %v", err)
			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Budget exceeded: UserID=%s, BudgetID=%d, Excess=%.2f",
			event.UserID, event.BudgetID, event.ExcessAmount)

		if err := c.notificationService.HandleBudgetExceeded(event); err != nil {
			log.Printf("[Consumer]  Error handling budget exceeded: %v", err)

			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Budget exceeded notification created")
		d.Ack(false)
	}
	return nil

}
func (c *Consumer) ConsumeBudgetWarning() error {
	q, err := c.ch.QueueDeclare(
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

	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("[Consumer]  Starting to consume budget_warning events...")

	for d := range msgs {
		var event BudgetWarningEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("[Consumer]  Error unmarshaling budget warning event: %v", err)
			d.Nack(false, false)
			continue
		}

		log.Printf("[Consumer] 📥 Budget warning: UserID=%s, BudgetID=%d, Used=%.0f%%",
			event.UserID, event.BudgetID, event.WarningPercent)

		if err := c.notificationService.HandleBudgetWarning(event); err != nil {
			log.Printf("[Consumer]  Error handling budget warning: %v", err)
			d.Nack(false, true)
			continue
		}
		log.Printf("[Consumer]  Budget warning notification created")
		d.Ack(false)
	}
	return nil
}

func (c *Consumer) ConsumeLowBalance() error {
	q, err := c.ch.QueueDeclare(
		"low_balance",
		true,
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
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	log.Println("[Consumer]  Starting to consume low_balance events...")
	for d := range msgs {
		var event LowBalanceEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Printf("[Consumer] Error unmarshaling event: %v", err)
			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Received low_balance event: %v", event)
		if err := c.notificationService.HandleLowBalance(event); err != nil {
			log.Printf("[Consumer]  Error handling low_balance event: %v", err)
			d.Nack(false, false)
			continue
		}
		log.Printf("[Consumer]  Low balance notification created")
		d.Ack(false)

	}
	return nil

}
func (c *Consumer) ConsumeAll() {
	go func() {
		if err := c.ConsumeTransactionCreated(); err != nil {
			log.Printf("[Consumer] Error consuming transaction created event: %v", err)
		}

	}()

	go func() {
		if err := c.ConsumeBudgetExceeded(); err != nil {
			log.Printf("[Consumer] Error consuming budget exceeded event: %v", err)
		}

	}()
	go func() {
		if err := c.ConsumeBudgetWarning(); err != nil {
			log.Printf("[Consumer] Error consuming budget warning event: %v", err)
		}

	}()
	go func() {
		if err := c.ConsumeLowBalance(); err != nil {
			log.Printf("[Consumer] Error consuming low balance event: %v", err)
		}

	}()
	log.Println("[Consumer]  Starting to consume all events...")

}
