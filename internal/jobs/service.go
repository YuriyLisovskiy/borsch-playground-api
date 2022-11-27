package jobs

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YuriyLisovskiy/borsch-playground-api/internal/amqp"
	"github.com/YuriyLisovskiy/borsch-runner-service/pkg/messages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JobService interface {
	GetJobResult(form *GetJobForm) (interface{}, error)
	GetJobOutputsAsJsonResult(form *GetJobOutputsForm) (interface{}, error)
	GetJobOutputsAsTextResult(form *GetJobOutputsForm, writer gin.ResponseWriter) error
	CreateJobResult(form *CreateJobForm) (interface{}, error)
}

type JobServiceImpl struct {
	repository  JobRepository
	amqpService amqp.JobService
}

func NewJobServiceImpl(jobRepository JobRepository, amqpService amqp.JobService) *JobServiceImpl {
	return &JobServiceImpl{
		repository:  jobRepository,
		amqpService: amqpService,
	}
}

func (js *JobServiceImpl) GetJobResult(form *GetJobForm) (interface{}, error) {
	job, err := js.repository.GetJob(form.JobId)
	if err != nil {
		return nil, err
	}

	job.OutputUrl = job.GetOutputUrl(form.RequestHost, form.RequestURI)
	return job, nil
}

func (js *JobServiceImpl) GetJobOutputsAsJsonResult(form *GetJobOutputsForm) (interface{}, error) {
	job, err := js.repository.GetJob(form.JobId)
	if err != nil {
		return nil, err
	}

	outputsChan, err := js.repository.GetJobOutputs(form.JobId, form.Offset, form.Limit)
	if err != nil {
		return nil, err
	}

	var outputs []JobOutputRow
	for output := range outputsChan {
		outputs = append(outputs, output)
	}

	return gin.H{"status": job.Status, "rows": outputs, "message": job.Message}, nil
}

func (js *JobServiceImpl) GetJobOutputsAsTextResult(form *GetJobOutputsForm, writer gin.ResponseWriter) error {
	header := writer.Header()
	header.Set("Transfer-Encoding", "chunked")
	header.Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	flusher := writer.(http.Flusher)
	flusher.Flush()
	job, err := js.repository.GetJob(form.JobId)
	if err != nil {
		return err
	}

	if job.Message != "" {
		_, err := writer.Write([]byte(job.Message))
		if err != nil {
			log.Println(err)
			return nil
		}
	} else {
		outputs, err := js.repository.GetJobOutputs(form.JobId, form.Offset, form.Limit)
		if err != nil {
			return err
		}

		for output := range outputs {
			_, err = writer.Write([]byte(fmt.Sprintf("%s\n", output.Text)))
			if err != nil {
				log.Println(err)
				return nil
			}

			flusher.Flush()
		}
	}

	flusher.Flush()
	return nil
}

func (js *JobServiceImpl) CreateJobResult(form *CreateJobForm) (interface{}, error) {
	job := &Job{
		ID:            uuid.New().String(),
		SourceCodeB64: base64.StdEncoding.EncodeToString([]byte(form.SourceCode)),
		ExitCode:      nil,
		Status:        JobStatusAccepted,
	}

	err := js.repository.CreateJob(job)
	if err != nil {
		return nil, err
	}

	err = js.publishAndUpdateJob(form, job)
	if err != nil {
		log.Println(err)
	}

	return gin.H{"job_id": job.ID, "output_url": job.GetOutputUrl(form.RequestHost, form.RequestURI)}, nil
}

// publishAndUpdateJob pushes the job to the RabbitMQ and update its status.
func (js *JobServiceImpl) publishAndUpdateJob(form *CreateJobForm, job *Job) error {
	jobMessage := messages.JobMessage{
		ID:            job.ID,
		LangVersion:   form.LangVersion,
		SourceCodeB64: job.SourceCodeB64,
		Timeout:       1 * time.Second, // TODO: replace with user.quota
	}
	err := js.amqpService.Publish(&jobMessage)
	if err != nil {
		log.Printf("failed to publish job: %v", err)
		job.Status = JobStatusRejected
	} else {
		job.Status = JobStatusQueued
	}

	err = js.repository.UpdateJob(job)
	if err != nil {
		return fmt.Errorf("failed to update job: %v", err)
	}

	return nil
}
