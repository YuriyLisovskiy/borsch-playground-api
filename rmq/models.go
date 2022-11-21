/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package rmq

type JobMessage struct {
	ID            string `json:"id"`
	LangVersion   string `json:"lang_version"`
	SourceCodeB64 string `json:"source_code_b64"`
}

type jobResultType string

const (
	jobResultLog  jobResultType = "log"
	jobResultExit               = "exit"
)

type JobResultMessage struct {
	ID   string        `json:"id"`
	Type jobResultType `json:"type"`
	Data string        `json:"data"`
}
