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
	"encoding/json"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
)

func ConsumeLPDing() {
	for {
		L := redi.PopAllLPDing()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendLPDingList(L)
	}
}

func SendLPDingList(L []*model.LPDing) {
	for _, lpding := range L {
		LPDingWorkerChan <- 1
		go SendLPDing(lpding)
	}
}

func SendLPDing(lpding *model.LPDing) {
	defer func() {
		<-LPDingWorkerChan
	}()

	tos := strings.Split(lpding.Tos, ",")

	url := g.Config().Api.LPDing
	client_id := g.Config().Api.LPClientId
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Header("X-Requested-With", "XMLHttpRequest")
	r.Header("Accept-Encoding", "identity")
	r.Header("Content-Type", "application/x-www-form-urlencoded")
	data := struct {
		AlarmCode string   `json:"alarmCode"`
		Content   string   `json:"content"`
		AT        []string `json:"_at"`
	}{
		AlarmCode: g.Config().Api.LPAlarmCode,
		Content:   lpding.Content,
		AT:        tos,
	}

	datajb, err := json.Marshal(data)
	if err != nil {
		log.Errorf("send lpding fail, receiver:%s, subject:%s, cotent:%s, error:%v", lpding.Tos, lpding.Subject, lpding.Content, err)
		return
	}
	body := "client_id=" + client_id + "&currentUserId=0&data=" + string(datajb)
	r.Body(body)

	resp, err := r.String()
	if err != nil {
		log.Errorf("send lpding fail, receiver:%s, subject:%s, cotent:%s, error:%v", lpding.Tos, lpding.Subject, lpding.Content, err)
	}

	log.Debugf("send lpding:%v, resp:%v, url:%s", lpding, resp, url)
}
