package app

import (
	"fmt"
	"net/http"
)

func GetURL(urls map[string]string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if http.MethodGet != req.Method {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		id := req.PathValue(`id`)

		url, ok := urls[id]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(url)
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}
