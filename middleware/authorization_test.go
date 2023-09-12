package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorizationMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	tt := []struct {
		name           string
		token          string
		noteID         string
		expectedStatus int
	}{
		{"No Token", "", "", http.StatusUnauthorized},
		{"Admin Token", "admintoken123", "1", http.StatusOK},
		{"Admin Token with no noteID", "admintoken123", "", http.StatusOK},
		{"User Token with access", "usertoken123", "1", http.StatusOK},
		{"User Token with no access", "usertoken123", "2", http.StatusForbidden},
		{"User Token with no noteID", "usertoken123", "", http.StatusBadRequest},
		{"Invalid Token", "invalidtoken123", "1", http.StatusForbidden},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/notes", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.token != "" {
				req.Header.Set("Authorization", tc.token)
			}

			if tc.noteID != "" {
				q := req.URL.Query()
				q.Add("note_id", tc.noteID)
				req.URL.RawQuery = q.Encode()
			}

			rr := httptest.NewRecorder()

			middleware := AuthorizationMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, rr.Code)
			}
		})
	}
}
