package handler

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vdgalyns/link-shortener/internal/repository"
	"github.com/vdgalyns/link-shortener/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_Get(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "link incorrect",
			request: "/",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				response:    "link incorrect",
			},
		},
		{
			name:    "link incorrect #2",
			request: "/abc123",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				response:    "link not found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositories := repository.NewRepository()
			services := service.NewService(repositories)
			handlers := NewHandler(services)

			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.Get)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			r, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			rBody := strings.TrimSuffix(string(r), "\n")
			assert.Equal(t, tt.want.response, rBody)
		})
	}
}

func TestHandler_Post(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		response    string
	}
	tests := []struct {
		name    string
		request string
		body    string
		want    want
	}{
		{
			name:    "link empty",
			request: "/",
			body:    "",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				response:    "link cannot be empty",
			},
		},
		{
			name:    "link incorrect",
			request: "/",
			body:    "http://ya.",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				response:    "link incorrect",
			},
		},
		{
			name:    "link correct",
			request: "/",
			body:    "http://ya.ru",
			want: want{
				statusCode:  201,
				contentType: "text/plain",
				response:    "http://example.com/1b556b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositories := repository.NewRepository()
			services := service.NewService(repositories)
			handlers := NewHandler(services)

			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.Add)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			r, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			rBody := strings.TrimSuffix(string(r), "\n")
			assert.Equal(t, tt.want.response, rBody)
		})
	}
}
