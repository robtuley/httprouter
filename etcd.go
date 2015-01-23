package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/robtuley/report"
)

type etcdResp struct {
	Action  string   `json:"action"`
	ErrCode int      `json:"errorCode"`
	ErrMsg  string   `json:"Message"`
	Node    etcdNode `json:"node"`
}

type etcdNode struct {
	Key      string     `json:"key"`
	Value    string     `json:"value"`
	IsDir    bool       `json:"dir"`
	Index    int        `json:"modifiedIndex"`
	Children []etcdNode `json:"nodes"`
}

func longPollForKeyChanges(key string) (chan etcdResp, error) {
	watchC := make(chan etcdResp)

	keyUrl, err := url.Parse("http://127.0.0.1:4001/v2/keys/" + key)
	if err != nil {
		report.Action("etcd.url.error", report.Data{"key": key})
		return nil, err
	}
	params := url.Values{}
	//params.Add("wait", "true")
	params.Add("recursive", "true")
	keyUrl.RawQuery = params.Encode()

	go func() {
		for {
			tick := report.Tick()
			resp, err := doEtcdRequest(keyUrl)
			report.Tock(tick, "etcd.response", report.Data{
				"url": keyUrl,
			})
			if err != nil {
				report.Action("etcd.response.error", report.Data{
					"url":   keyUrl,
					"error": err.Error(),
				})
				time.Sleep(time.Duration(5 * time.Second))
				continue
			}

			watchC <- resp
			time.Sleep(time.Duration(5 * time.Second))
		}
	}()

	return watchC, nil
}

func doEtcdRequest(keyUrl *url.URL) (etcdResp, error) {
	var msg etcdResp
	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	resp, err := client.Get(keyUrl.String())
	if err != nil {
		return msg, err
	}
	if resp.StatusCode != http.StatusOK {
		return msg, fmt.Errorf("etcd status code %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = decoder.Decode(&msg); err != nil {
		return msg, err
	}
	if msg.ErrCode != 0 {
		return msg, fmt.Errorf("etcd error %d %s", msg.ErrCode, msg.ErrMsg)
	}

	return msg, nil
}
