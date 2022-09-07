/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package settings

type Runner struct {
	Image     string `json:"image"`
	TagSuffix string `json:"tag_suffix"`
	Shell     string `json:"shell"`
	Command   string `json:"command"`
}
