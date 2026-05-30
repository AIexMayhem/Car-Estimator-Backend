package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ParseJson(r io.Reader, v any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can't read json: %v", err)
	}

	if err = json.Unmarshal(bytes, v); err != nil {
		return fmt.Errorf("can't unpack json: %v", err)
	}

	return nil
}

func RenderJson(w http.ResponseWriter, v any) {
	bytes, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Can't render JSON from object:", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
