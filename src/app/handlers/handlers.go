package handlers

import (
  "fmt"
  "net/http"
  "encoding/json"
  "text/template"
  "github.com/gorilla/mux"
  "lib/dbLogger"
  "app/sessions"
  "lib/userManagement"
)

func success(w http.ResponseWriter) {
  w.WriteHeader(http.StatusOK)
  w.Header().Set("Access-Control-Allow-Origin", "*")
}

func ClientHeartbeats(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  heartbeats := dbLogger.GetHeartbeats(vars["clientID"])
  json, _ := json.Marshal(heartbeats)
  success(w)
  fmt.Fprintf(w, "%v", string(json))
}

func Index(w http.ResponseWriter, r *http.Request) {
  heartbeats := dbLogger.ListHeartbeats()
  json, _ := json.Marshal(heartbeats)
  success(w)
  fmt.Fprintf(w, "%v", string(json))
}
func Heartbeats(w http.ResponseWriter, r *http.Request) {
  heartbeats := dbLogger.ListHeartbeats()
  json, _ := json.Marshal(heartbeats)
  success(w)
  fmt.Fprintf(w, "%v", string(json))
}

func Login(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  session, err := sessions.Store.Get(r, "session")
  loginTemplate, _ := template.ParseFiles("app/src/login.html")
  if err != nil {
    loginTemplate.Execute(w, nil)
  } else {
    if r.Method == "POST" {
      if userManagement.Verify(vars["email"], vars["password"]) {
        session.Values["loggedin"] = "true"
        session.Save(r, w)
        http.Redirect(w, r, "/", 302)
      }
    } else if r.Method == "GET" {
      loginTemplate.Execute(w, nil)
    }
  }
}

func Signup(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  signupTemplate, _ := template.ParseFiles("app/src/signup.html")
  if r.Method == "POST" {
    userManagement.CreateUser(vars["email"], vars["password"])
    http.Redirect(w, r, "/", 302)
  } else if r.Method == "GET" {
    signupTemplate.Execute(w, nil)
  }
}
