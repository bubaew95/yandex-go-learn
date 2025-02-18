package middleware

import "net/http"

func CookieMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:     "user_id",
			Value:    "112",
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
		h.ServeHTTP(w, r)
	})
}
