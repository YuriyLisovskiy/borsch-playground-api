/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

import (
	"log"
	"net/http"

	"github.com/YuriyLisovskiy/borsch-playground-api/jobs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobInfoHandler struct {
	jobId string
	db    *gorm.DB
}

func (h *JobInfoHandler) Write(out string) {
	var job jobs.Job
	err := h.db.First(&job, "ID = ?", h.jobId).Error
	if err != nil {
		log.Println(err)
		return
	}

	job.Outputs = append(job.Outputs, jobs.JobOutputRowDbModel{Text: out})
	h.db.Model(&job).Update("Outputs", job.Outputs)
}

func (h *JobInfoHandler) OnError(err error) {
	// TODO:
	log.Printf("[JOB ERROR]: %v\n", err)
}

func (h *JobInfoHandler) OnExit(exitCode int, exitErr error) {
	if exitErr != nil {
		h.OnError(exitErr)
	}

	var job jobs.Job
	err := h.db.First(&job, "ID = ?", h.jobId).Error
	if err != nil {
		h.OnError(err)
		return
	}

	job.ExitCode = new(int)
	*job.ExitCode = exitCode
	h.db.Model(&job).Update("ExitCode", job.ExitCode)
}

func jsonHandler(handler func(*gin.Context) (int, interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		status, obj, err := handler(c)
		if err != nil {
			log.Println(err)
			if status != -1 {
				c.JSON(status, gin.H{"message": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
		} else {
			c.JSON(status, obj)
		}
	}
}

func (a *Application) sendJsonError(c *gin.Context, status int, err error) {
	log.Println(err)
	if status != -1 {
		c.JSON(
			status, gin.H{
				"message":           err.Error(),
				"documentation_url": a.settings.ApiDocumentationUrl,
			},
		)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
	}
}

func (a *Application) sendTxtError(c *gin.Context, status int, err error) {
	log.Println(err)
	if status != -1 {
		c.String(status, err.Error())
	} else {
		c.String(http.StatusInternalServerError, "internal error")
	}
}

func stringArrayContains(array []string, item string) bool {
	for _, elem := range array {
		if elem == item {
			return true
		}
	}

	return false
}
