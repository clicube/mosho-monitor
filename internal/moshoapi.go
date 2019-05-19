package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type MoshoApi struct {
	client    *http.Client
	basicUser string
	basicPass string
}

const (
	baseUrl = "https://asia-northeast1-mosho-166cd.cloudfunctions.net"
)

func NewMoshoApi() (*MoshoApi, error) {
	client := http.DefaultClient

	txtbyte, err := ioutil.ReadFile("basic.txt")
	if err != nil {
		log.Printf("Error: Failed to read credential file: %+v", err)
		return nil, err
	}
	txt := strings.TrimRight(string(txtbyte), "\n")
	cred := strings.Split(string(txt), ":")

	return &MoshoApi{
		client:    client,
		basicUser: cred[0],
		basicPass: cred[1],
	}, nil
}

func (m *MoshoApi) Send(env *Env) error {
	log.Printf("MoshoApi: Sending env: %+v", env)

	time := time.Now().Unix()
	postEntity := jsonEnv{
		Temperature: env.Temperature,
		Humidity:    env.Humidity,
		Brightness:  env.Brightness,
		Time:        time,
	}

	postJson, err := json.Marshal(postEntity)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseUrl+"/api/v1/envs", bytes.NewBuffer(postJson))
	if err != nil {
		return err
	}

	req.SetBasicAuth(m.basicUser, m.basicPass)
	req.Header.Set("Content-Type", "application/json")

	res, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("MoshoApi: Status: %s", res.Status)
	}
	log.Printf("MoshoApi: Status: %s", res.Status)

	return nil
}

type jsonEnv struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Brightness  float64 `json:"brightness"`
	Time        int64   `json:"time"`
}
