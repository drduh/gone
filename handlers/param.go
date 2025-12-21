package handlers

import "net/http"

// getRequestParameter returns a parameter from the request
// URL or a form value.
func getRequestParameter(
	r *http.Request, pathLen int, fieldName string) string {
	if pathLen > len(r.URL.Path) {
		return ""
	}
	p := r.URL.Path[pathLen:]
	if p != "" {
		return p
	}
	if queryValue := r.URL.Query().Get(fieldName); queryValue != "" {
		return queryValue
	}
	return r.FormValue(fieldName)
}
