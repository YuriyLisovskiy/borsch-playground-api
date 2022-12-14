/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package jobs

import "gorm.io/gorm"

type JobService interface {
	GetJob(id string) (*Job, error)
	CreateJob(job *Job) error
	UpdateJob(job *Job) error
	GetJobOutputs(jobId string, offset, limit int) ([]JobOutputRow, error)
}

type JobServiceImpl struct {
	db *gorm.DB
}

func NewJobServiceImpl(db *gorm.DB) *JobServiceImpl {
	return &JobServiceImpl{db: db}
}

func (js *JobServiceImpl) GetJob(id string) (*Job, error) {
	job := &Job{}
	return job, js.db.First(job, "ID = ?", id).Error
}

func (js *JobServiceImpl) CreateJob(job *Job) error {
	return js.db.Create(job).Error
}

func (js *JobServiceImpl) UpdateJob(job *Job) error {
	return js.db.Save(&job).Error
}

func (js *JobServiceImpl) GetJobOutputs(jobId string, offset, limit int) ([]JobOutputRow, error) {
	_, err := js.GetJob(jobId)
	if err != nil {
		return nil, err
	}

	var outputs []JobOutputRow
	err = js.db.Offset(offset).Limit(limit).Find(&outputs, "job_id = ?", jobId).Error
	return outputs, err
}
