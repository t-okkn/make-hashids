package main

import (
	_ "net/http"
	_ "github.com/gin-gonic/gin"
)

func main() {
	// http.ListenAndServe(":8080", SetupRouter())
	SetupRouter().Run()
}

