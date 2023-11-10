package handler_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"track/api/router"
	"track/types"

	"github.com/gin-gonic/gin"
)

func trackDomain(r *gin.Engine, domain string) string {
	req, err := http.NewRequest("POST", "/track/"+domain, nil)
	if err != nil {
		panic(err)
	}

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		panic("failed to track domain")
	}

	var resp types.TrackResponse

	json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		panic(err)
	}

	if resp.Success != true {
		panic("failed to track domain")
	}

	if resp.Message != "ok" {
		panic("failed to track domain")
	}

	if resp.Key == "" {
		panic("failed to track domain")
	}

	return string(resp.Key)
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

		var resp types.TrackResponse

		json.Unmarshal(rr.Body.Bytes(), &resp)

		if err != nil {
			t.Fatal(err)
		}

		if resp.Success != true {
			t.Errorf("handler returned wrong success value: got %v want %v",
				resp.Success, true)
		}

		if resp.Message != "ok" {
			t.Errorf("handler returned wrong message: got %v want %v",
				resp.Message, "ok")
		}

		if resp.Key == "" {
			t.Errorf("handler returned empty key")
		}

	})
}

func TestHit(t *testing.T) {
	t.Run("tracks a hit", func(t *testing.T) {
		r := router.SetUpRouter("http://localhost:8080")

		key := trackDomain(r, "test.com")

		hitUrl := "http://test.com/index.html"

		encodedUrl := base64.URLEncoding.EncodeToString([]byte(hitUrl))

		req, err := http.NewRequest("GET", fmt.Sprintf("/hit/%s?k=%s", encodedUrl, key), nil)
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
