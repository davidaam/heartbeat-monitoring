package middleware

import (
  "net/http"
  "../sessions"
)

func AuthMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      if sessions.IsLoggedIn(r) || r.URL.Path == "/login" {
        h.ServeHTTP(w, r)
      } else {
        http.Redirect(w, r, "/login", 302)
      }
    })
}
