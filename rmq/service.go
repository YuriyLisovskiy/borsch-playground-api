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
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"borsch-playground-api/jobs"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQPJobService interface {
	ConsumeJobResults() error
	PublishJob(job *JobMessage) error
}

const (
	EnvRabbitMQServer      = "RABBITMQ_SERVER"
	EnvRabbitMQJobQueue    = "RABBITMQ_JOB_QUEUE"
	EnvRabbitMQResultQueue = "RABBITMQ_RESULT_QUEUE"
)

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
	mq.jobChannel, mq.jobQueue, err = createQueue(connection, os.Getenv(EnvRabbitMQJobQueue))
	if err != nil {
		return err
	}

	mq.jobResultChannel, mq.jobResultQueue, err = createQueue(connection, os.Getenv(EnvRabbitMQResultQueue))
	if err != nil {
		return err
	}

	return nil
}

func (mq *RabbitMQJobService) CleanUp() {
	logOrNil(mq.connection.Close())
	logOrNil(mq.jobChannel.Close())
	logOrNil(mq.jobResultChannel.Close())
}

func (mq *RabbitMQJobService) ConsumeJobResults() error {
	messages, err := mq.jobResultChannel.Consume(mq.jobResultQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go mq.processMessagesAsync(messages)
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
		"",
		mq.jobQueue.Name,
		false,
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
		job.Outputs = append(job.Outputs, jobs.JobOutputRow{Text: jobResult.Data})
		job.Status = jobs.JobStatusRunning
	case jobResultExit:
		job.ExitCode = new(int)
		*job.ExitCode, err = strconv.Atoi(jobResult.Data)
		job.Status = jobs.JobStatusFinished
	default:
		return fmt.Errorf("invalid type of job result: %s", jobResult.Type)
	}

	return mq.JobService.UpdateJob(job)
}

func (mq *RabbitMQJobService) processMessagesAsync(messages <-chan amqp.Delivery) {
	for d := range messages {
		err := mq.processJobResult(d.Body)
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		err = d.Ack(false)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func createQueue(connection *amqp.Connection, name string) (*amqp.Channel, amqp.Queue, error) {
	if name == "" {
		return nil, amqp.Queue{}, errors.New("RabbitMQ queue is not set")
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %v", err)
	}

	queue, err := channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to set QoS: %v", err)
	}

	return channel, queue, nil
}

func logOrNil(err error) {
	if err != nil {
		log.Println(err)
	}
}
