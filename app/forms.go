/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

type CreateJobForm struct {
	LangV      string `json:"lang_v"`
	SourceCode string `json:"source_code"`
}
