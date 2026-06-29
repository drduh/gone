package handlers

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteJSON tests valid and invalid encoding.
func TestWriteJSON(t *testing.T) {
	type File struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name     string
		code     int
		data     any
		wantCode int
		wantBody string
		wantType string
	}{
		{
			name:     "map encoding",
			code:     http.StatusOK,
			data:     map[string]string{"hello": "world"},
			wantCode: http.StatusOK,
			wantBody: `{"hello":"world"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "struct encoding",
			code:     http.StatusOK,
			data:     File{ID: 1, Name: "testFile"},
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"name":"testFile"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "slice encoding",
			code:     http.StatusOK,
			data:     []int{1, 2, 3},
			wantCode: http.StatusOK,
			wantBody: `[1,2,3]` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "nil encoding",
			code:     http.StatusOK,
			data:     nil,
			wantCode: http.StatusOK,
			wantBody: "null\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "non-200 status code",
			code:     http.StatusCreated,
			data:     File{ID: 1, Name: "testFile"},
			wantCode: http.StatusCreated,
			wantBody: `{"id":1,"name":"testFile"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "unicode content",
			code:     http.StatusOK,
			data:     map[string]string{"msg": "héllo wörld"},
			wantCode: http.StatusOK,
			wantBody: `{"msg":"héllo wörld"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "marshal failure",
			code:     http.StatusOK,
			data:     math.NaN(),
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"failed to encode response"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "empty struct",
			code:     http.StatusOK,
			data:     File{},
			wantCode: http.StatusOK,
			wantBody: `{"id":0,"name":""}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "empty slice",
			code:     http.StatusOK,
			data:     []int{},
			wantCode: http.StatusOK,
			wantBody: `[]` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "nil slice",
			code:     http.StatusOK,
			data:     []int(nil),
			wantCode: http.StatusOK,
			wantBody: "null\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name:     "html escaping",
			code:     http.StatusOK,
			data:     map[string]string{"msg": "<script>alert('xss')</script>"},
			wantCode: http.StatusOK,
			wantBody: `{"msg":"\u003cscript\u003ealert('xss')\u003c/script\u003e"}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
		{
			name: "nested struct",
			code: http.StatusOK,
			data: struct {
				File File `json:"file"`
			}{File: File{ID: 1, Name: "testFile"}},
			wantCode: http.StatusOK,
			wantBody: `{"file":{"id":1,"name":"testFile"}}` + "\n",
			wantType: "application/json; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeJSON(rr, tt.code, tt.data)

			if rr.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d",
					tt.wantCode, rr.Code)
			}
			if ct := rr.Header().Get("Content-Type"); ct != tt.wantType {
				t.Errorf("expected Content-Type %q, got %q",
					tt.wantType, ct)
			}
			if body := rr.Body.String(); body != tt.wantBody {
				t.Errorf("want - got: \n%q\n%q",
					tt.wantBody, body)
			}
		})
	}
}
