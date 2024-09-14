package main

import (
	"ssl-checker/auth"
	"ssl-checker/cache"
	"ssl-checker/domains"
	"ssl-checker/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	go cache.UpdateCache()
	auth.LoadUsers()
	domains_utils.LoadDomains()

	r := gin.Default()
	routes.InitRoutes(r)

	r.Run(":8082")
}
