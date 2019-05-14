package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

var (
	release = getEnv("RELEASE", "v0.0.1")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// static
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/img", "./static/img")
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")

	// routes
	r.GET("/", homeHandler)
	r.GET("/health", healthHandler)
	r.GET("/prime", defaultPrimeHandler)
	r.GET("/prime/:max", primeArgHandler)

	// port
	port := getEnv(portVariableName, defaultPort)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting: %s \n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
