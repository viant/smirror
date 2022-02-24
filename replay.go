package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/viant/smirror/replay"
	"log"
	"net/http"
)

//StorageReplay cloud function entry point
func StorageReplay(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		defer func() {
			_ = r.Body.Close()
		}()
	}
	err := replayUnprocessed(w, r)
	if err != nil {
		log.Print(err)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func replayUnprocessed(writer http.ResponseWriter, httpRequest *http.Request) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	request := &replay.Request{}
	if err = json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return errors.Wrapf(err, "failed to decode %T", request)
	}
	service := replay.Singleton()
	response := service.Replay(context.Background(), request)
	if data, err := json.Marshal(response); err == nil {
		fmt.Printf("%s\n", data)
	}
	if err = json.NewEncoder(writer).Encode(response); err != nil {
		return err
	}
	return err
}
