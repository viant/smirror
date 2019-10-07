package smirror

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"smirror/mon"
)

//StorageMonitor cloud function entry point
func StorageMonitor(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 0 {
		defer func() {
			_ = r.Body.Close()
		}()
	}
	err := checkStorage(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func checkStorage(writer http.ResponseWriter, httpRequest *http.Request) error {
	request := &mon.Request{}
	if err := json.NewDecoder(httpRequest.Body).Decode(&request); err != nil {
		return errors.Wrapf(err, "failed to decode %T", request)
	}
	service, err := mon.NewFromEnv(ConfigEnvKey)
	if err != nil {
		return err
	}
	response := service.Check(context.Background(), request)
	if err =  json.NewEncoder(writer).Encode(response);err != nil {
		log.Fatal(err)
	}
	return nil
}
