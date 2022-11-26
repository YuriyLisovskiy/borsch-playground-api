/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/YuriyLisovskiy/borsch-playground-api/internal/jobs"
	msg "github.com/YuriyLisovskiy/borsch-runner-service/pkg/messages"
	amqp "github.com/rabbitmq/amqp091-go"
)

type JobService interface {
	ConsumeResults() error
	Publish(job *msg.JobMessage) error
}

const (
	EnvRabbitMQServer      = "RABBITMQ_SERVER"
	EnvRabbitMQJobQueue    = "RABBITMQ_JOB_QUEUE"
	EnvRabbitMQResultQueue = "RABBITMQ_RESULT_QUEUE"
)

type RabbitMQJobService struct {
	Server     string
	JobService jobs.JobRepository

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

func (mq *RabbitMQJobService) ConsumeResults() error {
	messages, err := mq.jobResultChannel.Consume(mq.jobResultQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go mq.processMessagesAsync(messages)
	return nil
}

func (mq *RabbitMQJobService) Publish(job *msg.JobMessage) error {
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
	jobResult := msg.JobResultMessage{}
	err := json.Unmarshal(data, &jobResult)
	if err != nil {
		return err
	}

	job, err := mq.JobService.GetJob(jobResult.ID)
	if err != nil {
		return err
	}

	if jobResult.Type == msg.JobResultLog {
		log.Printf("[%s] LOG\n", jobResult.ID)
		return mq.JobService.CreateOutput(
			&jobs.JobOutputRow{
				Text:  jobResult.Data,
				JobID: jobResult.ID,
			},
		)
	}

	switch jobResult.Type {
	case msg.JobResultStart:
		log.Printf("[%s] STARTED\n", jobResult.ID)
		job.Status = jobs.JobStatusRunning
	case msg.JobResultExit:
		log.Printf("[%s] EXIT\n", jobResult.ID)
		job.ExitCode = new(int)
		*job.ExitCode = *jobResult.ExitCode
		job.Status = jobs.JobStatusFinished
		job.FinishedAt = new(time.Time)
		*job.FinishedAt = time.Now().UTC()
	case msg.JobResultError:
		log.Printf("[%s] ERROR\n", jobResult.ID)
		job.Message = jobResult.Data
		job.Status = jobs.JobStatusFinished
		job.FinishedAt = new(time.Time)
		*job.FinishedAt = time.Now().UTC()
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
