package main

import (
    "net/http"
    "./routes"
)

func main() {
  r := routes.Router()
  http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}
