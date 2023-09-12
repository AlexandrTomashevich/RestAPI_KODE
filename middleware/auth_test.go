package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	tt := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"No Token", "", http.StatusUnauthorized},
		{"Invalid Token", "invalidtoken123", http.StatusUnauthorized},
		{"User Token", "usertoken123", http.StatusOK},
		{"Admin Token", "admintoken123", http.StatusOK},
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

			rr := httptest.NewRecorder()

			middleware := AuthorizationMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, rr.Code)
			}
		})
	}
}
