package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	Close() error
	BindQueue(key string) (q amqp.Queue, err error)
	Reconnect() error
	HandleConsumedDeliveries(key string, handle func(payload map[string]interface{}) error)
	Publish(key string, message []byte, attempts int32) error
}

type MessageBody struct {
	Data []byte
	Type string
}

type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}

type RabbitMQ struct {
	Ctx          context.Context
	Uri          string
	Exchange     string
	Kind         string
	Prefix       string
	RetryCount   int
	Err          chan error
	conn         *amqp.Connection
	channel      *amqp.Channel
	consumerTags map[string]bool
}

func (c *RabbitMQ) Connect() error {
	var err error
	c.conn, err = amqp.Dial(c.Uri)
	if err != nil {
		return err
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.Err <- errors.New("BunnyConnection Closed")
	}()
	c.consumerTags = make(map[string]bool, 0)
	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}
	if err := c.channel.ExchangeDeclare(c.Exchange, c.Kind, true, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

func (c *RabbitMQ) Close() error {
	for consumerTag := range c.consumerTags {
		if err := c.channel.Cancel(consumerTag, true); err != nil {
			return err
		}
	}

	return c.conn.Close()
}

func (c *RabbitMQ) BindQueue(key string) (q amqp.Queue, err error) {
	if q, err = c.channel.QueueDeclare(fmt.Sprintf("%s:%s", c.Prefix, key), true, false, false, false, nil); err != nil {
		return q, err
	}
	if err := c.channel.QueueBind(q.Name, key, c.Exchange, false, nil); err != nil {
		return q, err
	}
	if err := c.channel.Qos(1, 0, false); err != nil {
		return q, err
	}

	return q, nil
}

func (c *RabbitMQ) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}

	return nil
}

// Consumer Things
func (c *RabbitMQ) consume(key string) (<-chan amqp.Delivery, error) {
	q, err := c.BindQueue(key)
	if err != nil {
		return nil, err
	}

	consumerTag := fmt.Sprintf("ctag:%s", q.Name)
	c.consumerTags[consumerTag] = true
	deliveries, err := c.channel.Consume(q.Name, consumerTag, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (c *RabbitMQ) HandleConsumedDeliveries(key string, handle func(payload map[string]interface{}) error) {
	delivery, err := c.consume(key)
	if err != nil {
		panic(err)
	}

	for {
		go handler(*c, key, delivery, handle)
		if err := <-c.Err; err != nil {
			fmt.Println(fmt.Sprintf("[CONSUMER] RabbitMQ connection closed: %s", err.Error()))

			c.Reconnect()
			deliveries, err := c.consume(key)
			if err != nil {
				panic(err)
			}
			delivery = deliveries
		}
	}
}

// Publisher Things
func (c *RabbitMQ) Publish(key string, message []byte, attempts int32) error {
	if attempts > int32(c.RetryCount) {
		fmt.Println(fmt.Sprintf("[PUBLISHER] Too many attempts: %s", key))
		return nil
	}

	select {
	case err := <-c.Err:
		if err != nil {
			fmt.Println(fmt.Sprintf("[PUBLISHER] RabbitMQ connection closed: %s", err.Error()))
			c.Reconnect()
		}
	default:
	}

	p := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
		Headers: amqp.Table{
			"x-attempts": attempts,
		},
	}

	fmt.Println(fmt.Sprintf("[PUBLISHER] Publishing message: %s - %d", key, attempts))

	if err := c.channel.PublishWithContext(c.Ctx, c.Exchange, key, false, false, p); err != nil {
		fmt.Println(fmt.Sprintf("[PUBLISHER] Error in publishing message: %s", err.Error()))

		c.Reconnect()
		return c.Publish(key, message, attempts+1)
	}

	fmt.Println(fmt.Sprintf("[PUBLISHER] Published message: %s - %d", key, attempts))
	return nil
}

func handler(c RabbitMQ, key string, messages <-chan amqp.Delivery, handle func(payload map[string]interface{}) error) {
	for message := range messages {
		fmt.Println(fmt.Sprintf("[*] Received message: %s", key))
		attempts := int32(1)

		if message.Headers["x-attempts"] != nil {
			attempts = message.Headers["x-attempts"].(int32)
		}

		d, _ := deserialize(message.Body)
		err := handle(d)

		message.Ack(false)

		if err != nil {
			fmt.Println(err.Error())

			if attempts < int32(c.RetryCount) {
				c.Publish(key, message.Body, attempts+int32(1))
			} else {
				fmt.Println(fmt.Sprintf("Too many attempts: %s", key))
			}
		}
	}
}

func deserialize(b []byte) (map[string]interface{}, error) {
	var msg map[string]interface{}
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber() // need for dealing with large number in shop_id
	err := decoder.Decode(&msg)
	return msg, err
}
