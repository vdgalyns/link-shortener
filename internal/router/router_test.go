package router

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdgalyns/link-shortener/internal/config"
	"github.com/vdgalyns/link-shortener/internal/handler"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/service"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (int, string) {
	var b io.Reader
	if body != "" {
		b = bytes.NewBuffer([]byte(body))
	}
	req, err := http.NewRequest(method, ts.URL+path, b)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func testRequestWithJSON(t *testing.T, ts *httptest.Server, method, path string, body interface{}) (int, []byte) {
	var b bytes.Buffer
	jsonEncoder := json.NewEncoder(&b)
	err := jsonEncoder.Encode(body)
	require.NoError(t, err)
	req, err := http.NewRequest(method, ts.URL+path, &b)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, respBody
}

func NewTestServer() (*httptest.Server, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	repositories := repository.NewRepository(cfg)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services, cfg)

	r := NewRouter(handlers)
	ts := httptest.NewServer(r)
	return ts, nil
}

func TestGet(t *testing.T) {
	ts, err := NewTestServer()
	if err != nil {
		log.Fatal(err)
	}
	defer ts.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		responseCode int
		responseBody string
	}{
		{
			name:         "NotExistLink",
			method:       http.MethodGet,
			path:         "/abc123",
			responseCode: http.StatusBadRequest,
			responseBody: "link not found",
		},
		{
			name:         "IncorrectLink",
			method:       http.MethodGet,
			path:         "/abc123456789",
			responseCode: http.StatusBadRequest,
			responseBody: "link incorrect",
		},
		{
			name:         "NotExistRoute",
			method:       http.MethodGet,
			path:         "/",
			responseCode: http.StatusMethodNotAllowed,
			responseBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, tt.method, tt.path, "")
			assert.Equal(t, tt.responseCode, statusCode)
			assert.Equal(t, tt.responseBody, strings.TrimSuffix(body, "\n"))
		})
	}
}

func TestAdd(t *testing.T) {
	ts, err := NewTestServer()
	if err != nil {
		log.Fatal(err)
	}
	defer ts.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		requestBody  string
		responseCode int
		responseBody string
	}{
		{
			name:         "BodyEmpty",
			method:       http.MethodPost,
			path:         "/",
			requestBody:  "",
			responseCode: http.StatusBadRequest,
			responseBody: "link cannot be empty",
		},
		{
			name:         "IncorrectLink",
			method:       http.MethodPost,
			path:         "/",
			requestBody:  "http:",
			responseCode: http.StatusBadRequest,
			responseBody: "link incorrect",
		},
		{
			name:         "CorrectLink",
			method:       http.MethodPost,
			path:         "/",
			requestBody:  "http://ya.ru",
			responseCode: http.StatusCreated,
			responseBody: ts.URL + "/1b556b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequest(t, ts, tt.method, tt.path, tt.requestBody)
			assert.Equal(t, tt.responseCode, statusCode)
			assert.Equal(t, tt.responseBody, strings.TrimSuffix(body, "\n"))
		})
	}
}

func TestAddWithJSON(t *testing.T) {
	ts, err := NewTestServer()
	if err != nil {
		log.Fatal(err)
	}
	defer ts.Close()

	type RequestBody struct {
		URL string `json:"url"`
	}
	type ResponseBody struct {
		Result string `json:"result"`
	}

	tests := []struct {
		name         string
		method       string
		path         string
		requestBody  RequestBody
		responseCode int
		responseBody interface{}
	}{
		{
			name:         "EmptyLink",
			method:       http.MethodPost,
			path:         "/api/shorten",
			requestBody:  RequestBody{},
			responseCode: http.StatusBadRequest,
			responseBody: "link cannot be empty",
		},
		{
			name:         "IncorrectLink",
			method:       http.MethodPost,
			path:         "/api/shorten",
			requestBody:  RequestBody{URL: "http"},
			responseCode: http.StatusBadRequest,
			responseBody: "link incorrect",
		},
		{
			name:         "CorrectLink",
			method:       http.MethodPost,
			path:         "/api/shorten",
			requestBody:  RequestBody{URL: "http://ya.ru"},
			responseCode: http.StatusCreated,
			responseBody: ResponseBody{Result: ts.URL + "/1b556b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := testRequestWithJSON(t, ts, tt.method, tt.path, tt.requestBody)
			assert.Equal(t, tt.responseCode, statusCode)
			if statusCode != http.StatusCreated {
				assert.Equal(t, tt.responseBody, strings.TrimSuffix(string(body), "\n"))
				return
			}
			responseBody := ResponseBody{}
			err := json.Unmarshal(body, &responseBody)
			require.NoError(t, err)
			assert.Equal(t, tt.responseBody, responseBody)
		})
	}
}
