/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *Application) sendJsonError(c *gin.Context, status int, err error) {
	log.Println(err)
	if status != -1 {
		c.JSON(
			status, gin.H{
				"message":           err.Error(),
				"documentation_url": a.settings.ApiDocumentationUrl,
			},
		)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
	}
}

func stringArrayContains(array []string, item string) bool {
	for _, elem := range array {
		if elem == item {
			return true
		}
	}

	return false
}
