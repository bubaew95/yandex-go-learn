package middlewares

import (
	"net/http"
	"strings"

	"github.com/bubaew95/yandex-go-learn/internal/compress"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
)

func GZipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		logger.Log.Info("Content-type: " + r.Header.Get("Content-type"))
		if isContentEncoding(r) {
			cr, err := compress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		if isAcceptEncoding(r) && isAccessContentType(r) {
			cw := compress.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

func isAcceptEncoding(r *http.Request) bool {
	logger.Log.Info("Accept-Encoding isset")

	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}

func isContentEncoding(r *http.Request) bool {
	logger.Log.Info("Content-Encoding isset")

	contentEncoding := r.Header.Get("Content-Encoding")
	return strings.Contains(contentEncoding, "gzip")
}

func isAccessContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-type")
	return contentType == "application/json" || contentType == "text/html"
}
