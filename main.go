package main

import (
	"log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/maxprime/util"


)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

var (
	release   = util.MustGetEnv("RELEASE", "v0.0.0 (Not set)")
)

func main() {

	// router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/img", "./static/img")
	r.Static("/css", "./static/css")

	// root & health
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", map[string]interface{}{
			"release": release,
		})
	})
	r.GET("/health", healthHandler)

	// prime
	r.GET("/prime", defaultPrimeHandler)
	r.GET("/prime/:max", primeWithArgHandler)


	// port
	port := util.MustGetEnv(portVariableName, defaultPort)
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting: %s \n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}

}


func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}