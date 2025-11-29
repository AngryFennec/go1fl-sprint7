package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", v.count)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)

		responseBody := strings.TrimSpace(response.Body.String())

		var cafes []string
		if len(responseBody) != 0 {
			cafes = strings.Split(string(responseBody), ",")
			assert.Equal(t, v.want, len(cafes))
		} else {
			assert.Equal(t, v.want, 0)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"вилка", 1},
		{"кофе", 2},
	}

	for _, v := range request {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)

		responseBody := strings.ToLower(strings.TrimSpace(response.Body.String()))

		var cafes []string
		var count int
		if len(responseBody) != 0 {
			cafes = strings.Split(responseBody, ",")
			for _, cafe := range cafes {
				if strings.Contains(cafe, v.search) {
					count++
				}
			}
			assert.Equal(t, v.want, count)
		} else {
			assert.Equal(t, v.want, 0)
		}
	}
}
