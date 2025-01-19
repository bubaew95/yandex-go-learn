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

		accessContentTypes := map[string]bool{
			"application/json": true,
			"text/plain":       true,
		}

		contentType := r.Header.Get("content-type")
		_, accessContentType := accessContentTypes[contentType]

		if accessContentType {
			logger.Log.Debug("accessContentTypes")
			acceptEncoding := r.Header.Get("Accept-Encoding")
			isSupportGZip := strings.Contains(acceptEncoding, "gzip")
			if isSupportGZip {
				logger.Log.Debug("Accept-Encoding run")

				cw := compresses.NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			isSendGZip := strings.Contains(contentEncoding, "gzip")
			if isSendGZip {
				logger.Log.Debug("Content-Encoding run")

				cr, err := compresses.NewCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = cr
				defer cr.Close()
			}
		}

		h.ServeHTTP(ow, r)
	})
}
