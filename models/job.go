/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package models

type JobOutputRowDbModel struct {
	Model

	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Text  string `json:"text"`
	JobID string `json:"job_id"`
}

type JobDbModel struct {
	Model

	Code       string                `json:"code"`
	Outputs    []JobOutputRowDbModel `json:"-" gorm:"foreignKey:JobID"`
	ExitCode   *int                  `json:"exit_code"`
	OutputsUrl string                `json:"outputs_url" gorm:"-:all"`
}
