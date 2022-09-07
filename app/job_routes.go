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
	"fmt"
	"net/http"
	"strconv"

	"github.com/YuriyLisovskiy/borsch-playground-service/core"
	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (a *Application) addJobRoutes(rg *gin.RouterGroup) {
	jobRouter := rg.Group("/jobs")

	jobRouter.GET("/:id", jsonHandler(a.getJobHandler))
	jobRouter.GET("/:id/outputs", jsonHandler(a.getJobOutputsHandler))
	jobRouter.POST("/", jsonHandler(a.createJobHandler))
}

func (a *Application) getJobHandler(c *gin.Context) (int, interface{}, error) {
	var job models.JobDbModel
	err := a.db.First(&job, "ID = ?", c.Param("id")).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("job not found")
		}

		return -1, nil, err
	}

	job.OutputsUrl = fmt.Sprintf("%s://%s%s/outputs", "http", c.Request.Host, c.Request.RequestURI)
	return http.StatusOK, job, nil
}

func (a *Application) getJobOutputsHandler(c *gin.Context) (int, interface{}, error) {
	jobId := c.Param("id")
	qOffset := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(qOffset)
	if err != nil {
		jObj := map[string]interface{}{
			"message": "Validation Failed",
			"error": map[string]string{
				"resource": "Job Outputs",
				"value":    qOffset,
			},
			"documentation_url": "TODO:",
		}
		return http.StatusBadRequest, jObj, nil
	}

	err = a.db.Model(&models.JobDbModel{}).First(nil, "ID = ?", jobId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http.StatusNotFound, nil, errors.New("job not found")
		}

		return -1, nil, err
	}

	var jobOutputs []models.JobOutputRowDbModel
	err = a.db.Offset(offset).Find(&jobOutputs, "job_id = ?", jobId).Error
	if err != nil {
		return -1, nil, err
	}

	return http.StatusOK, jobOutputs, nil
}

func (a *Application) createJobHandler(c *gin.Context) (int, interface{}, error) {
	var form CreateJobForm
	err := c.ShouldBindJSON(&form)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	if !stringArrayContains(a.settings.BorschVersions, form.LangV) {
		return http.StatusBadRequest, nil, errors.New("Language version does not exist")
	}

	if len(form.SourceCode) == 0 {
		return http.StatusBadRequest, nil, errors.New("Code is not provided or empty")
	}

	job := &models.JobDbModel{
		Model: models.Model{
			ID: uuid.New().String(),
		},
		Code:     form.SourceCode,
		Outputs:  []models.JobOutputRowDbModel{},
		ExitCode: nil,
	}

	err = a.db.Create(job).Error
	if err != nil {
		return -1, nil, err
	}

	err = a.enqueueJob(job, form.LangV)
	if err != nil {
		switch err {
		case core.ErrQueueIsFull:
			return http.StatusServiceUnavailable, nil, errors.New("Service temporary is unavailable, try again later")
		case core.ErrQueueIsUnavailable:
			fallthrough
		default:
			return -1, nil, err
		}
	}

	return http.StatusCreated, gin.H{"job_id": job.ID}, nil
}
