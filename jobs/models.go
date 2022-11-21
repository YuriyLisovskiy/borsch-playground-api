/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package jobs

import (
	"fmt"
	"strings"

	"borsch-playground-api/common"
	"github.com/gin-gonic/gin"
)

type JobOutputRow struct {
	common.Model

	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Text  string `json:"text"`
	JobID string `json:"job_id"`
}

type JobStatus string

const (
	JobStatusAccepted JobStatus = "accepted"
	JobStatusRejected JobStatus = "rejected"
	JobStatusQueued   JobStatus = "queued"
	JobStatusRunning  JobStatus = "running"
	JobStatusFinished JobStatus = "finished"
)

type Job struct {
	common.Model

	Code      string         `json:"source_code"`
	Outputs   []JobOutputRow `json:"-" gorm:"foreignKey:JobID"`
	ExitCode  *int           `json:"exit_code"`
	OutputUrl string         `json:"output_url" gorm:"-:all"`
	Status    JobStatus      `json:"status"`
}

func (m *Job) GetOutputUrl(c *gin.Context) string {
	rURI := c.Request.RequestURI
	if !strings.HasSuffix(rURI, m.ID) {
		rURI += m.ID
	}

	return fmt.Sprintf("%s://%s%s/output", "http", c.Request.Host, rURI)
}
