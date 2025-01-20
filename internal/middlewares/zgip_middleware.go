package middlewares

import (
	"net/http"
	"strings"

	"github.com/bubaew95/yandex-go-learn/internal/compresses"
	"github.com/bubaew95/yandex-go-learn/internal/logger"
)

func GZipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		logger.Log.Info("accessContentTypes")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isSupportGZip := strings.Contains(acceptEncoding, "gzip")
		if isSupportGZip {
			logger.Log.Info("Accept-Encoding run")

			cw := compresses.NewCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isSendGZip := strings.Contains(contentEncoding, "gzip")
		if isSendGZip {
			logger.Log.Info("Content-Encoding run")

			cr, err := compresses.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
