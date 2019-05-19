package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type BoardCtl struct{}

func NewBoardCtl() (*BoardCtl, error) {
	log.Println("BoardCtl: Initialized")
	return &BoardCtl{}, nil
}

type result struct {
	Result      string  `json:"result"`
	Message     string  `json:"message"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Brightness  float64 `json:"brightness"`
}

func (b *BoardCtl) Get() (*Env, error) {
	log.Printf("BoardCtl: Getting data")
	out, err := exec.Command(
		"boardctl",
		"env",
	).Output()
	log.Printf("BoardCtl: Output: %s", string(out))

	// outputをパースしてみる
	var res result
	jerr := json.Unmarshal(out, &res)

	// パースできなければ
	if jerr != nil {

		if err != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf(string(out))
		}

	}

	if res.Result != "ok" {
		return nil, fmt.Errorf(res.Message)
	}

	return &Env{
		Temperature: res.Temperature,
		Humidity:    res.Humidity,
		Brightness:  res.Brightness,
	}, nil
}
