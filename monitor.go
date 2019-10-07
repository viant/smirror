package smirror

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"smirror/mon"
)

//StorageMonitor cloud function entry point
func StorageMonitor(w http.ResponseWriter, r *http.Request) {
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
	return json.NewEncoder(writer).Encode(response)
}
