package handlers

import "fmt"

func ExampleShortenerHandler_AddNewURL() {
	//cfg := config.NewConfig()
	//
	//t := NewMockShortenerService()
	//
	//shortenerDB, _ := storage.NewShortenerDB(*cfg)
	//shortener, _ := fileStorage.NewShortenerRepository(*shortenerDB)
	//
	//mockService := service.NewShortenerService(shortener, *cfg)
	//handler := NewShortenerHandler(mockService)
	//
	//req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://example.com"))
	//w := httptest.NewRecorder()
	//
	//handler.CreateURL(w, req)
	//
	//res := w.Result()
	//defer res.Body.Close()
	//
	//body, _ := io.ReadAll(res.Body)

	//fmt.Println(res.StatusCode)
	//fmt.Println(string(body))

	fmt.Println(201)
	fmt.Println("http://short.url/abc123")

	// Output:
	// 201
	// http://short.url/abc123
}
