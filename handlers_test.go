package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type handler = func(http.ResponseWriter, *http.Request)

func passRequest(r *http.Request, h handler) *httptest.ResponseRecorder {
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	handler.ServeHTTP(responseRecorder, r)
	return responseRecorder
}

func TestMainHandlerBasic(t *testing.T) {
	u := url.URL{Path: "/cafe?"}
	values := url.Values{}
	values.Add("city", "moscow")
	values.Add("count", "1")
	u.RawQuery = values.Encode()

	req := httptest.NewRequest(http.MethodGet, u.String(), nil)
	resp := passRequest(req, mainHandle)

	require.Equal(t, http.StatusOK, resp.Code)
	require.NotEmpty(t, resp.Body)
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	u := url.URL{Path: "/cafe?"}
	city := "moscow"
	count := strconv.Itoa(len(cafeList[city]) + 1)
	values := url.Values{}
	values.Add("city", city)
	values.Add("count", count)
	u.RawQuery = values.Encode()

	req := httptest.NewRequest(http.MethodGet, u.String(), nil)
	resp := passRequest(req, mainHandle)
	assert.Equal(t, http.StatusOK, resp.Code)

	returnedCafeList := strings.Split(resp.Body.String(), ",")
	expectedCafeListLen := len(cafeList[city])
	assert.Len(t, returnedCafeList, expectedCafeListLen)
}

func TestMainHandlerWhenCityNotFound(t *testing.T) {
	u := url.URL{Path: "/cafe?"}
	city := "NotPresent"
	count := "1"
	values := url.Values{}
	values.Add("city", city)
	values.Add("count", count)
	u.RawQuery = values.Encode()

	_, found := cafeList[city]
	require.False(t, found, "test data broken, %q must absent in cafeList", city)

	req := httptest.NewRequest(http.MethodGet, u.String(), nil)
	resp := passRequest(req, mainHandle)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "wrong city value", resp.Body.String())
}
