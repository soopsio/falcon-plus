// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
)

func ConsumeDing() {
	for {
		L := redi.PopAllDing()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendDingList(L)
	}
}

func SendDingList(L []*model.Ding) {
	for _, ding := range L {
		DingWorkerChan <- 1
		go SendDing(ding)
	}
}

func SendDing(ding *model.Ding) {
	defer func() {
		<-DingWorkerChan
	}()

	url := g.Config().Api.Ding
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Param("tos", ding.Tos)
	r.Param("content", ding.Content)

	resp, err := r.String()
	if err != nil {
		log.Errorf("send ding fail, receiver:%s, subject:%s, cotent:%s, error:%v", ding.Tos, ding.Subject, ding.Content, err)
	}

	log.Debugf("send ding:%v, resp:%v, url:%s", ding, resp, url)
}
