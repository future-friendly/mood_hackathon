package util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func DecodeRequest(source io.ReadCloser, dest interface{}) error{
	body, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func WriteSuccess(w http.ResponseWriter) {
	resp, err := json.Marshal(map[string]bool{"ok": true})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(200)
	w.Write(resp)
}

type ErrorResponse struct {
	Error string `json:"error"`
}
func WriteError(w http.ResponseWriter, code int, e error) {
	resp, err := json.Marshal(ErrorResponse{Error: e.Error()})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(code)
	w.Write(resp)
}