package internal

import (
	"log"
	"time"
)

type Env struct {
	Temperature float64
	Humidity    float64
	Brightness  float64
}

func Start() error {

	errch := make(chan error)

	board,err := NewBoardCtl()
	if err != nil {
		return nil
	}

	datadog,err := NewDatadog()
	if err != nil {
		return nil
	}

	moshoapi,err := NewMoshoApi()
	if err != nil {
		return nil
	}

	go func() {

		ticker := time.NewTicker(time.Minute * 5)
		defer ticker.Stop()

		for {
			err := doCheck(board, datadog, moshoapi)
			if err != nil {
				errch <- err
				return
			}
			select {
			case <-ticker.C:
			}
		}

	}()

	select {
	case err := <-errch:
		return err
	}
}

func doCheck(board *BoardCtl, datadog *Datadog, moshoapi *MoshoApi) error {
	env, err := board.Get()
	if err != nil {
		log.Printf("Monitor: Error: board.Get(): %s", err)
		return nil
	}

	err = datadog.Send(env)
	if err != nil {
		log.Printf("Monitor: datadog.Send(env): %s", err)
	}

	err = moshoapi.Send(env)
	if err != nil {
		log.Printf("Monitor: moshoapi.Send(env): %s", err)
	}
	return nil
}
