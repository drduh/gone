package handlers

import (
	"encoding/json"
	"net/http"
)

// deny serves a JSON response for disallowed requests.
func deny(w http.ResponseWriter, code int, reason string) {
	writeJSON(w, code, errorJSON(reason))
}

// errorJSON returns an error string map containing the string.
func errorJSON(s string) map[string]string {
	return map[string]string{"error": s}
}

// writeJSON serves a JSON response with data.
func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set(
		"Content-Type", "application/json; charset=utf-8")

	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"failed to encode response"}`,
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(code)
	_, _ = w.Write(append(buf, '\n'))
}
