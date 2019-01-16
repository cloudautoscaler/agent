package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	validToken := "avalidtoken"
	validLoad := 0.13
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{}
		err := json.NewDecoder(r.Body).Decode(&data)
		sendError := func(msg string, e int) {
			t.Log(msg)
			w.WriteHeader(e)
		}
		if err != nil {
			sendError("invalid json", http.StatusBadRequest)
			return
		}
		load, ok := data["load"].(float64)
		if !ok {
			sendError("expected a valid \"load\" key as float64", http.StatusBadRequest)
			return
		}
		if token, ok := data["token"].(string); !ok || token != validToken {
			sendError("expected a valid \"token\" key as string", http.StatusUnauthorized)
			return
		}
		t.Log("load received:", load)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()
	// Test the happy path.
	os.Setenv(envToken, validToken)
	os.Setenv(envURI, server.URL)
	err := send(validLoad)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	// Invalid token must return a specific error.
	os.Setenv(envToken, "invalidtoken")
	err = send(validLoad)
	if err != errNoAuth {
		t.Fatal("invalid error:", err)
	}
}
