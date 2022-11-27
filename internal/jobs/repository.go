/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package jobs

import (
	"log"

	"gorm.io/gorm"
)

type JobRepository interface {
	GetJob(id string) (*Job, error)
	CreateJob(job *Job) error
	UpdateJob(job *Job) error
	GetJobOutputs(jobId string, offset, limit int) (chan JobOutputRow, error)
	CreateOutput(output *JobOutputRow) error
}

type JobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepositoryImpl(db *gorm.DB) *JobRepositoryImpl {
	return &JobRepositoryImpl{db: db}
}

func (js *JobRepositoryImpl) GetJob(id string) (*Job, error) {
	job := &Job{}
	return job, js.db.First(job, "ID = ?", id).Error
}

func (js *JobRepositoryImpl) CreateJob(job *Job) error {
	return js.db.Create(job).Error
}

func (js *JobRepositoryImpl) UpdateJob(job *Job) error {
	return js.db.Save(job).Error
}

// func (js *JobRepositoryImpl) GetJobOutputs(jobId string, offset, limit int) ([]JobOutputRow, error) {
// 	_, err := js.GetJob(jobId)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var outputs []JobOutputRow
// 	err = js.db.Offset(offset).Limit(limit).Find(&outputs, "job_id = ?", jobId).Error
// 	return outputs, err
// }

func (js *JobRepositoryImpl) GetJobOutputs(jobId string, offset, limit int) (chan JobOutputRow, error) {
	tx := js.db.Model(&JobOutputRow{}).Offset(offset).Limit(limit).Where("job_id = ?", jobId)
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	jobChan := make(chan JobOutputRow)
	go func() {
		defer close(jobChan)
		for rows.Next() {
			jobOutput := JobOutputRow{}
			err = js.db.ScanRows(rows, &jobOutput)
			if err != nil {
				log.Println(err)
			} else {
				jobChan <- jobOutput
			}
		}
	}()

	return jobChan, nil
}

func (js *JobRepositoryImpl) CreateOutput(output *JobOutputRow) error {
	return js.db.Create(output).Error
}
