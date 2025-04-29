package handlers

import (
	"bytes"
	"fmt"
	"github.com/bubaew95/yandex-go-learn/internal/core/ports"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleShortenerHandler_AddNewURL() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	mockService := ports.NewMockShortenerService(ctrl)
	mockService.EXPECT().
		GenerateURL(gomock.Any(), "https://example.com", gomock.Any()).
		Return("http://short.url/abc123", nil)

	handler := NewShortenerHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	w := httptest.NewRecorder()

	handler.CreateURL(w, req)

	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(string(body))

	// Output:
	// 201
	// http://short.url/abc123
}
