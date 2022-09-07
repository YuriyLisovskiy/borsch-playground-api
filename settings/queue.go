/*
 * Borsch Playground API
 *
 * Copyright (C) 2022 Yuriy Lisovskiy - All Rights Reserved
 * You may use, distribute and modify this code under the
 * terms of the MIT license.
 */

package settings

import "github.com/YuriyLisovskiy/borsch-playground-service/core"

type Queue struct {
	Workers  int `json:"workers"`
	Capacity int `json:"capacity"`
}

func (q *Queue) Create() (*core.Queue, error) {
	jobQueue := core.NewQueue(q.Capacity)
	err := jobQueue.Open(q.Workers)
	if err != nil {
		return nil, err
	}

	return jobQueue, nil
}
