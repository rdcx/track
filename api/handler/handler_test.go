package handler_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"track/api/router"

	"github.com/gin-gonic/gin"
)

func trackDomain(r *gin.Engine, domain string) {
	req, err := http.NewRequest("POST", "/track/"+domain, nil)
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to track domain")
	}
}

func TestTrack(t *testing.T) {
	t.Run("tracks a domain", func(t *testing.T) {
		r := router.SetUpRouter("http://localhost:8080")

		req, err := http.NewRequest("POST", "/track/test.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
	})

	t.Run("returns 400 if domain already exists", func(t *testing.T) {
		r := router.SetUpRouter("http://localhost:8080")

		trackDomain(r, "test.com")

		req2, err := http.NewRequest("POST", "/track/test.com", nil)

		if err != nil {
			t.Fatal(err)
		}

		r2 := httptest.NewRecorder()
		r.ServeHTTP(r2, req2)

		if r2.Code != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				r2.Code, http.StatusBadRequest)
		}
	})
}

func TestHit(t *testing.T) {
	t.Run("tracks a hit", func(t *testing.T) {
		r := router.SetUpRouter("http://localhost:8080")

		trackDomain(r, "test.com")

		hitUrl := "http://test.com/index.html"

		encodedUrl := base64.URLEncoding.EncodeToString([]byte(hitUrl))

		req, err := http.NewRequest("GET", fmt.Sprintf("/hit/%s", encodedUrl), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, http.StatusOK)
		}
	})
}
