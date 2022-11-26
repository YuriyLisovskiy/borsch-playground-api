/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package server

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YuriyLisovskiy/borsch-playground-api/internal/jobs"
	"github.com/YuriyLisovskiy/borsch-runner-service/pkg/messages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *Application) getJobHandler(c *gin.Context) {
	job, err := a.jobService.GetJob(c.Param("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.sendJsonError(c, http.StatusNotFound, errors.New("job not found"))
		} else {
			a.sendJsonError(c, http.StatusInternalServerError, err)
		}

		return
	}

	job.OutputUrl = job.GetOutputUrl(c)
	c.JSON(http.StatusOK, job)
}

func (a *Application) getJobOutputHandler(c *gin.Context) {
	jobId := c.Param("id")
	offsetParam := c.DefaultQuery("offset", "-1")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		a.sendJsonError(c, http.StatusBadRequest, errors.New("offset is invalid integer value"))
		return
	}

	limitParam := c.DefaultQuery("limit", "-1")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		a.sendJsonError(c, http.StatusBadRequest, errors.New("limit is invalid integer value"))
		return
	}

	job, err := a.jobService.GetJob(jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.sendJsonError(c, http.StatusNotFound, errors.New("job not found"))
		} else {
			a.sendJsonError(c, http.StatusInternalServerError, err)
		}

		return
	}

	outputs, err := a.jobService.GetJobOutputs(jobId, offset, limit)
	if err != nil {
		a.sendJsonError(c, http.StatusInternalServerError, err)
		return
	}

	formatParam := c.DefaultQuery("format", "json")
	switch strings.ToLower(formatParam) {
	case "json":
		c.JSON(http.StatusOK, gin.H{"status": job.Status, "rows": outputs, "message": job.Message})
	case "txt":
		outputString := ""
		if job.Message != "" {
			outputString = job.Message
		} else {
			outputLen := len(outputs)
			for i, output := range outputs {
				outputString += output.Text
				if i < outputLen-1 {
					outputString += "\n"
				}
			}
		}

		c.String(http.StatusOK, outputString)
	default:
		a.sendJsonError(
			c,
			http.StatusBadRequest,
			errors.New("invalid response format, available values are 'json' an 'txt'"),
		)
	}
}

func (a *Application) createJobHandler(c *gin.Context) {
	var form CreateJobForm
	err := c.ShouldBindJSON(&form)
	if err != nil {
		a.sendJsonError(c, http.StatusBadRequest, err)
		return
	}

	if form.LangVersion == "" {
		a.sendJsonError(c, http.StatusBadRequest, errors.New("language version is not provided"))
		return
	}

	if !stringArrayContains(a.settings.AvailableInterpreters, form.LangVersion) {
		a.sendJsonError(c, http.StatusBadRequest, errors.New("language version does not exist"))
		return
	}

	if len(form.SourceCode) == 0 {
		a.sendJsonError(c, http.StatusBadRequest, errors.New("source code is not provided"))
		return
	}

	job := &jobs.Job{
		ID:            uuid.New().String(),
		SourceCodeB64: base64.StdEncoding.EncodeToString([]byte(form.SourceCode)),
		ExitCode:      nil,
		Status:        jobs.JobStatusAccepted,
	}

	err = a.jobService.CreateJob(job)
	if err != nil {
		a.sendJsonError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"job_id": job.ID, "output_url": job.GetOutputUrl(c)})
	a.publishJob(&form, job)
}

// publishJob pushes the job to the RabbitMQ and update its status.
func (a *Application) publishJob(form *CreateJobForm, job *jobs.Job) {
	jobMessage := messages.JobMessage{
		ID:            job.ID,
		LangVersion:   form.LangVersion,
		SourceCodeB64: job.SourceCodeB64,
		Timeout:       1 * time.Second, // TODO: replace with user.quota
	}
	err := a.amqpJobService.Publish(&jobMessage)
	if err != nil {
		log.Printf("Failed to publish job: %v", err)
		job.Status = jobs.JobStatusRejected
	} else {
		job.Status = jobs.JobStatusQueued
	}

	err = a.jobService.UpdateJob(job)
	if err != nil {
		log.Printf("Failed to update job: %v", err)
	}
}
