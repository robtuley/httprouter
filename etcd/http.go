package etcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/robtuley/report"
)

type Resp struct {
	Action  string `json:"action"`
	ErrCode int    `json:"errorCode"`
	ErrMsg  string `json:"Message"`
	Node    Node   `json:"node"`
}

type Node struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsDir    bool   `json:"dir"`
	Index    int    `json:"modifiedIndex"`
	Children []Node `json:"nodes"`
}

func Watch(key string) (chan Resp, chan report.Data) {
	watchC := make(chan Resp)
	errorC := make(chan report.Data)

	go func() {

		// Publish full data via recursive GET

		keyUrl, err := key2Url(key, map[string]string{"recursive": "true"})
		if err != nil {
			errorC <- report.Data{"error": err.Error(), "key": key}
			close(watchC)
			return
		}

		for {
			resp, err := doRequest(keyUrl)
			if err != nil {
				errorC <- report.Data{"error": err.Error(), "url": keyUrl.String()}
				time.Sleep(time.Second * 2)
				continue // re-try
			}
			watchC <- resp
			break
		}

		// Listen to changes via HTTP long poll & publish

		keyUrl, err = key2Url(key, map[string]string{"recursive": "true", "wait": "true"})
		errorC <- report.Data{"url": keyUrl.String()}
		if err != nil {
			errorC <- report.Data{"error": err.Error(), "key": key}
			close(watchC)
			return
		}

		for {
			resp, err := doRequest(keyUrl)
			if err != nil {
				errorC <- report.Data{"error": err.Error(), "url": keyUrl.String()}
				time.Sleep(time.Second * 2)
				continue
			}
			watchC <- resp
			time.Sleep(time.Second)
		}

	}()

	return watchC, errorC
}

func doRequest(keyUrl *url.URL) (Resp, error) {
	var msg Resp
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

func key2Url(key string, opts map[string]string) (*url.URL, error) {
	keyUrl, err := url.Parse("http://127.0.0.1:4001/v2/keys/" + key)
	if err == nil {
		params := url.Values{}
		for k, v := range opts {
			params.Add(k, v)
		}
		keyUrl.RawQuery = params.Encode()
	}
	return keyUrl, err
}
