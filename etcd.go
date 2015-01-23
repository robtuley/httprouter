package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/robtuley/report"
)

type etcdNode struct {
	Key           string
	Value         string
	Dir           bool
	CreatedIndex  int
	ModifiedIndex int
}

func longPollForJson(keyUrl string) chan etcdNode {
	watchC := make(chan etcdNode)

	go func() {
		for {
			tick := report.Tick()
			nodes, err := doEtcdRequest(keyUrl)
			report.Tock(tick, "etcd.response", report.Data{
				"url":   keyUrl,
				"nodes": len(nodes),
			})
			if err != nil {
				report.Action("etcd.response.error", report.Data{
					"url":   keyUrl,
					"error": err.Error(),
				})
				time.Sleep(time.Duration(5 * time.Second))
				continue
			}

			for _, nd := range nodes {
				watchC <- nd
			}
			time.Sleep(time.Duration(5 * time.Second))
		}
	}()

	return watchC
}

func doEtcdRequest(keyUrl string) ([]etcdNode, error) {
	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	resp, err := client.Get(keyUrl)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("etcd status code %d", resp.StatusCode)
	}

	/*
		dump, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(dump))
	*/

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var msg struct {
		ErrorCode int
		Message   string
		Node      struct {
			Nodes []etcdNode
		}
	}
	if err = decoder.Decode(&msg); err != nil {
		return nil, err
	}
	if msg.ErrorCode != 0 {
		return nil, fmt.Errorf("etcd error %d %s", msg.ErrorCode, msg.Message)
	}

	return msg.Node.Nodes, nil
}
