package merakiclient

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/marian-craciunescu/merakibeat/config"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/beat"
)

//api/v1/scanning/receiver/

const (
	moduleName = "scanReceiver"
	listenURL  = "/api/v1/scanning/receiver/"
)

type ScanReceiver struct {
	secret     string
	validator  string
	version    string
	mux        *http.ServeMux
	hostPort   string
	beatClient beat.Client
	certPATH   string
	keyPATH    string
	logger     *logp.Logger
}

func NewScanReceiver(config config.Config, bc beat.Client) *ScanReceiver {

	sr := ScanReceiver{
		secret:     config.ScanSecret,
		validator:  config.ScanValidator,
		certPATH:   config.ServerCert,
		keyPATH:    config.SerkerKey,
		version:    "2.0",
		mux:        http.NewServeMux(),
		hostPort:   fmt.Sprintf(":%d", config.ScanPort),
		beatClient: bc,
		logger:     logp.NewLogger(moduleName),
	}
	sr.mux.HandleFunc(listenURL, sr.handleReceive)

	return &sr
}

func (sr *ScanReceiver) handleReceive(w http.ResponseWriter, r *http.Request) {
	sr.logger.Debugf("Entering handleReceive")
	switch r.Method {
	case http.MethodGet:
		sr.handleReceiveValidation(w, r)
		return
	case http.MethodPost:
		sr.handleReceiveData(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		return
	}
}
func (sr *ScanReceiver) handleReceiveValidation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintln(w, sr.validator)
	if err != nil {
		sr.logger.Errorf("Error writing response in validation err=%s", err.Error())
	}
	return
}

func (sr *ScanReceiver) handleReceiveData(w http.ResponseWriter, r *http.Request) {
	sr.logger.Debugf("Entering scanreciever\n")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sr.logger.Error("Error reading body %s\n", err.Error())
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Body %s", string(body[:]))
	var scanData ScanData
	err = json.Unmarshal(body, &scanData)
	if err != nil {
		sr.logger.Error("Error un marshalling json body %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if scanData.Secret != sr.secret {
		sr.logger.Error("Error Ivalid secret \n")
		http.Error(w, "Invalid Secrect", http.StatusMethodNotAllowed)
		return
	}
	sr.logger.Debugf("Publishing scan data %+v\n", scanData)
	mapstrArr, err := scanData.GetMapStr("MerakiScanEvent", map[string]string{})
	for _, mapStr := range mapstrArr {
		seenTime, _ := mapStr.GetValue("client.seenTime")
		seenTimeStr, _ := seenTime.(string)
		ts, err := time.Parse("2006-01-02T15:04:05.999999999", seenTimeStr)
		sr.logger.Debugf("Timestamp %s %+v", seenTimeStr, ts)
		if err != nil {
			ts = time.Now()
		}
		sr.beatClient.Publish(beat.Event{
			Timestamp: ts,
			Fields:    mapStr,
		})
		sr.logger.Debugf("Published event %+v\n", mapStr)
	}
	return
}

func (sr *ScanReceiver) Run() {
	sr.logger.Fatal(
		http.ListenAndServeTLS(
			sr.hostPort, sr.certPATH, sr.keyPATH, sr.mux),
	)
}
