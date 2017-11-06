package main

import (
    "net/http"
    "app/routes"
)

func main() {
  r := routes.Router()
  http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}
