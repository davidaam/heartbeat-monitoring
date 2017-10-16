package main

import (
    "net/http"
    "./routes"
)

func main() {

  http.Handle("/", routes.Router())
  http.ListenAndServe(":8080", nil)
}
