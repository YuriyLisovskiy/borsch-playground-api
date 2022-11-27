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
	"time"
)

type JobOutputRow struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	JobID     string    `json:"job_id"`
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
	ID            string     `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"-"`
	FinishedAt    *time.Time `json:"finished_at"`
	Status        JobStatus  `json:"status"`
	SourceCodeB64 string     `json:"source_code_b64"`
	ExitCode      *int       `json:"exit_code"`
	Message       string     `json:"message"`

	OutputUrl string `json:"output_url" gorm:"-:all"`
}

func (m *Job) GetOutputUrl(host, uri string) string {
	if !strings.HasSuffix(uri, m.ID) {
		uri += m.ID
	}

	return fmt.Sprintf("%s://%s%s/output", "http", host, uri)
}
