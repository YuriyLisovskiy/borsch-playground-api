/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Settings struct {
	GinMode             string        `json:"gin_mode"`
	ShutdownTimeoutSec  time.Duration `json:"shutdown_timeout_sec"`
	BorschVersions      []string      `json:"borsch_versions"`
	ApiDocumentationUrl string        `json:"api_documentation_url"`
	Database            *Database     `json:"database"`
}

func (s *Settings) PerformChecks() error {
	switch s.GinMode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		break
	default:
		return fmt.Errorf(
			"invalid Gin mode, available values are '%s', '%s', '%s'",
			gin.DebugMode,
			gin.ReleaseMode,
			gin.TestMode,
		)
	}

	return nil
}

func Load() (*Settings, error) {
	localSettings, err := loadFromJson("settings.local.json")
	if err == nil {
		return localSettings, nil
	}

	return loadFromJson("settings.json")
}

func loadFromJson(filename string) (*Settings, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var s Settings
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
