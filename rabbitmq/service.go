/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/YuriyLisovskiy/borsch-playground-api/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPJobService interface {
	ConsumeJobResults() error
	PublishJob(job *JobMessage) error
}

type RabbitMQJobService struct {
	Server     string
	JobService jobs.JobService

	connection       *amqp.Connection
	jobChannel       *amqp.Channel
	jobResultChannel *amqp.Channel
	jobQueue         amqp.Queue
	jobResultQueue   amqp.Queue
}

func (mq *RabbitMQJobService) Setup() error {
	connection, err := amqp.Dial(mq.Server)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	mq.connection = connection
	// defer connection.Close()
	mq.jobChannel, mq.jobQueue, err = createQueue(connection, os.Getenv("RABBITMQ_JOB_QUEUE"))
	if err != nil {
		return err
	}

	mq.jobResultChannel, mq.jobResultQueue, err = createQueue(connection, os.Getenv("RABBITMQ_RESULT_QUEUE"))
	if err != nil {
		return err
	}

	return nil
}

func (mq *RabbitMQJobService) CleanUp() {
	defer mq.connection.Close()
	defer mq.jobChannel.Close()
	defer mq.jobResultChannel.Close()
}

func (mq *RabbitMQJobService) ConsumeJobResults() error {
	messages, err := mq.jobResultChannel.Consume(
		mq.jobResultQueue.Name, // queue
		"",                     // consumer
		false,                  // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for d := range messages {
			err = mq.processJobResult(d.Body)
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			err = d.Ack(false)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}()
	return nil
}

func (mq *RabbitMQJobService) PublishJob(job *JobMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = mq.jobChannel.PublishWithContext(
		ctx,
		"",               // exchange
		mq.jobQueue.Name, // routing key
		false,            // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	return nil
}

func (mq *RabbitMQJobService) processJobResult(data []byte) error {
	jobResult := JobResultMessage{}
	err := json.Unmarshal(data, &jobResult)
	if err != nil {
		return err
	}

	job, err := mq.JobService.GetJob(jobResult.ID)
	if err != nil {
		return err
	}

	switch jobResult.Type {
	case jobResultLog:
		job.Outputs = append(job.Outputs, jobs.JobOutputRowDbModel{Text: jobResult.Data})
	case jobResultExit:
		job.ExitCode = new(int)
		*job.ExitCode, err = strconv.Atoi(jobResult.Data)
	default:
		return nil
	}

	return mq.JobService.UpdateJob(job)
}

func createQueue(connection *amqp.Connection, name string) (*amqp.Channel, amqp.Queue, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %v", err)
	}

	// defer channel.Close()

	queue, err := channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to set QoS: %v", err)
	}

	return channel, queue, nil
}
