package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteJSON tests valid and invalid encoding.
func TestWriteJSON(t *testing.T) {
	t.Run("validate status and body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		payload := map[string]string{"hello": "world"}

		writeJSON(rr, http.StatusOK, payload)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d",
				http.StatusOK, rr.Code)
		}
		if ct := rr.Header().Get(
			"Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}

		var got map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response body: %v", err)
		}
		if got["hello"] != "world" {
			t.Errorf("expected hello=world, got %v", got)
		}
	})

	t.Run("validate 500 with invalid data", func(t *testing.T) {
		rr := httptest.NewRecorder()

		// force marshal error
		writeJSON(rr, http.StatusOK, math.NaN())

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected %d, got %d",
				http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("validate struct encoding", func(t *testing.T) {
		type User struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		rr := httptest.NewRecorder()
		writeJSON(rr, http.StatusOK,
			User{ID: 1, Name: "Bob"})

		if rr.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d",
				http.StatusOK, rr.Code)
		}

		var got User
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if got.ID != 1 || got.Name != "Bob" {
			t.Errorf("unexpected body: %+v", got)
		}
	})

	t.Run("validate slice encoding", func(t *testing.T) {
		rr := httptest.NewRecorder()
		writeJSON(rr, http.StatusOK, []int{1, 2, 3})

		var got []int
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if len(got) != 3 ||
			got[0] != 1 || got[1] != 2 || got[2] != 3 {
			t.Errorf("got invalid slice: %v", got)
		}
	})

	t.Run("validate nil encoding", func(t *testing.T) {
		rr := httptest.NewRecorder()
		writeJSON(rr, http.StatusOK, nil)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, got %d",
				http.StatusOK, rr.Code)
		}
		if body := rr.Body.String(); body != "null\n" {
			t.Errorf("expected null, got %q", body)
		}
	})
}
