package sessions

import (
    "net/http"
    "github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("secret-password"))

func IsLoggedIn(r *http.Request) bool {
    session, _ := Store.Get(r, "session")
    if session.Values["loggedin"] == "true" {
        return true
    }
    return false
}
