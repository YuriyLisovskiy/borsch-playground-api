/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Settings struct {
	GinMode               string
	ShutdownTimeoutSec    time.Duration
	AvailableInterpreters []string
	ApiDocumentationUrl   string
}

func LoadSettingsFromEnv() (*Settings, error) {
	shutdownTimeout, err := getShutdownTimeoutFromEnv()
	if err != nil {
		return nil, err
	}

	ginMode, err := getGinModeFromEnv()
	if err != nil {
		return nil, err
	}

	s := &Settings{
		GinMode:               ginMode,
		ShutdownTimeoutSec:    shutdownTimeout,
		AvailableInterpreters: strings.Split(os.Getenv("AVAILABLE_INTERPRETERS"), ","),
		ApiDocumentationUrl:   os.Getenv("API_DOC_URI"),
	}

	return s, nil
}

func getShutdownTimeoutFromEnv() (time.Duration, error) {
	shutdownTimeoutSec, err := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT_SEC"))
	if err != nil {
		return -1, err
	}

	if shutdownTimeoutSec < 0 {
		return -1, fmt.Errorf("shutdown timeout should be non-negative integer, received %d", shutdownTimeoutSec)
	}

	return time.Duration(shutdownTimeoutSec), nil
}

func getGinModeFromEnv() (string, error) {
	ginMode := os.Getenv("GIN_MODE")
	switch ginMode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return ginMode, nil
	default:
		return "", fmt.Errorf(
			"invalid Gin mode, available values are '%s', '%s', '%s'",
			gin.DebugMode,
			gin.ReleaseMode,
			gin.TestMode,
		)
	}
}
