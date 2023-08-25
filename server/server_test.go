package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock "github.com/takumakume/sbomreport-to-dependencytrack/mock"
)

func Test_healthzFunc(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthzFunc())

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func Test_uploadFunc(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := mock.NewMockUploader(ctrl)

	testCases := []struct {
		name         string
		method       string
		body         io.Reader
		mock         bool
		mockErr      error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "success",
			method:       http.MethodPost,
			body:         strings.NewReader("test body"),
			mock:         true,
			mockErr:      nil,
			expectedCode: http.StatusOK,
			expectedBody: "ok",
		},
		{
			name:         "request body is empty",
			method:       http.MethodPost,
			body:         nil,
			mock:         false,
			mockErr:      nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "request body is empty\n",
		},
		{
			name:         "method not allowed",
			method:       http.MethodGet,
			body:         nil,
			mock:         false,
			mockErr:      nil,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "method not allowed\n",
		},
		{
			name:         "uploader.Run() error",
			method:       http.MethodPost,
			body:         strings.NewReader("test body"),
			mock:         true,
			mockErr:      errors.New("mock error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: "mock error\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock uploader
			if tc.mock {
				mockUploader.EXPECT().Run(ctx, gomock.Any()).Return(tc.mockErr)
			}

			// Set up request
			req, err := http.NewRequest(tc.method, "/", tc.body)
			if err != nil {
				t.Fatal(err)
			}

			// Run function
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(uploadFunc(ctx, mockUploader))
			handler.ServeHTTP(rr, req)

			// Check response
			if status := rr.Code; status != tc.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedCode)
			}

			expected := tc.expectedBody
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
			}
		})
	}
}
