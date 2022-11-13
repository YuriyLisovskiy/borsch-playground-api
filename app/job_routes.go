/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-playground-api/common"
	"github.com/YuriyLisovskiy/borsch-playground-api/jobs"
	rmq "github.com/YuriyLisovskiy/borsch-playground-api/rabbitmq"
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
			a.sendJsonError(c, http.StatusNotFound, err)
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
		a.sendJsonError(c, http.StatusNotFound, err)
		return
	}

	formatParam := c.DefaultQuery("format", "json")
	switch strings.ToLower(formatParam) {
	case "json":
		c.JSON(http.StatusOK, gin.H{"exit_code": job.ExitCode, "rows": outputs})
	case "txt":
		outputLen := len(outputs)
		outputString := ""
		for i, output := range outputs {
			outputString += output.Text
			if i < outputLen-1 {
				outputString += "\n"
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

func (a *Application) createJobHandler(c *gin.Context) (int, interface{}, error) {
	var form CreateJobForm
	err := c.ShouldBindJSON(&form)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if !stringArrayContains(a.settings.BorschVersions, form.LangV) {
		return http.StatusBadRequest, nil, errors.New("language version does not exist")
	}

	if len(form.SourceCode) == 0 {
		return http.StatusBadRequest, nil, errors.New("code is not provided or empty")
	}

	job := &jobs.Job{
		Model: common.Model{
			ID: uuid.New().String(),
		},
		Code:     form.SourceCode,
		Outputs:  []jobs.JobOutputRowDbModel{},
		ExitCode: nil,
	}

	err = a.jobService.CreateJob(job)
	if err != nil {
		return -1, nil, err
	}

	jobMessage := rmq.JobMessage{
		ID:          job.ID,
		LangVersion: form.LangV,
		SourceCode:  form.SourceCode,
	}
	err = a.amqpJobService.PublishJob(&jobMessage)
	if err != nil {
		return -1, nil, err
	}

	return http.StatusCreated, gin.H{"job_id": job.ID, "output_url": job.GetOutputUrl(c)}, nil
}
