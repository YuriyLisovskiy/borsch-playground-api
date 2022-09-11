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

	"github.com/YuriyLisovskiy/borsch-playground-service/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type JobInfoHandler struct {
	jobId string
	db    *gorm.DB
}

func (h *JobInfoHandler) Write(out string) {
	var job models.JobDbModel
	err := h.db.First(&job, "ID = ?", h.jobId).Error
	if err != nil {
		log.Println(err)
		return
	}

	job.Outputs = append(job.Outputs, models.JobOutputRowDbModel{Text: out})
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

	var job models.JobDbModel
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

func txtHandler(handler func(*gin.Context) (int, string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		status, str, err := handler(c)
		if err != nil {
			log.Println(err)
			if status != -1 {
				c.String(http.StatusInternalServerError, err.Error())
			} else {
				c.String(http.StatusInternalServerError, "Internal error")
			}
		} else {
			c.String(status, str)
		}
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
