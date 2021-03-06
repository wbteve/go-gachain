// MIT License
//
// Copyright (c) 2016-2018 GACHAIN
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package daemons

import (
	"context"
	"fmt"
	"time"

	"github.com/GACHAIN/go-gachain/packages/consts"
	"github.com/GACHAIN/go-gachain/packages/model"
	"github.com/GACHAIN/go-gachain/packages/scheduler"
	"github.com/GACHAIN/go-gachain/packages/scheduler/contract"

	log "github.com/sirupsen/logrus"
)

func loadContractTasks() error {
	stateIDs, err := model.GetAllSystemStatesIDs()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("get all system states ids")
		return err
	}

	for _, stateID := range stateIDs {
		if !model.IsTable(fmt.Sprintf("%d_vde_cron", stateID)) {
			return nil
		}

		c := model.Cron{}
		c.SetTablePrefix(fmt.Sprintf("%d_vde", stateID))
		tasks, err := c.GetAllCronTasks()
		if err != nil {
			log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("get all cron tasks")
			return err
		}

		for _, cronTask := range tasks {
			err = scheduler.UpdateTask(&scheduler.Task{
				ID:       cronTask.UID(),
				CronSpec: cronTask.Cron,
				Handler: &contract.ContractHandler{
					Contract: cronTask.Contract,
				},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Scheduler starts contracts on schedule
func Scheduler(ctx context.Context, d *daemon) error {
	d.sleepTime = time.Hour
	return loadContractTasks()
}
