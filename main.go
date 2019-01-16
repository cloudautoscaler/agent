// This agent sends information about server load to AutoScaler.cloud backend.
// It just push and does not pull any information from the service.
// It only supports GNU/Linux systems.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cloudautoscaler/agent/backoff"
)

const (
	// 0: 5min
	// 1: 10min
	// 2: 15min
	loadInterval = 1
	// envToken is the environment variable name
	// expected to store the token for the server where the agent is running.
	envToken = "CAS_TOKEN"
	envURI   = "CAS_URI"
	URIPath  = "/agent/v1/load"
)

var errInvalid = errors.New("invalid loadavg info")
var errInvalidServerResp = errors.New("invalid server response")
var errNoAuth = errors.New("unauthorized request")

func main() {
	for range time.Tick(time.Second * 60) {
		load, err := sampleLoad()
		if err != nil {
			log.WithField("cause", err).Fatal("unable to get load, is this a Linux system?")
		}
		err = backoff.On(
			func() error {
				return send(load)
			},
			5,
		)
		ctx := log.WithField("cause", err)
		switch err {
		case nil:
			continue
		case errNoAuth:
			ctx.Fatal("unrecoverable error")
		default:
			ctx.Error("unable to send data")
		}
	}
}

func sampleLoad() (float64, error) {
	content, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, err
	}
	bb := bytes.Fields(content)
	if len(bb) < 3 {
		return 0, errInvalid
	}
	return strconv.ParseFloat(string(bb[loadInterval]), 64)

}

func send(load float64) error {
	token := os.Getenv(envToken)
	uri := os.Getenv(envURI)
	if token == "" {
		return fmt.Errorf("variable %q not found", envToken)
	}
	if uri == "" {
		return fmt.Errorf("variable %q not found", envURI)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	body := strings.NewReader(fmt.Sprintf(`{"load": %.2f}`, load))
	url := fmt.Sprintf("%s%s/%s", uri, URIPath, token)
	resp, err := client.Post(url, "application/json", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusAccepted:
		break
	case http.StatusUnauthorized:
		return errNoAuth
	default:
		return errInvalidServerResp
	}
	return nil
}
