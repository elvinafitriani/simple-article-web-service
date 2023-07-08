package main

import (
	"github.com/gin-gonic/gin"

	"article/connection"
	"article/routers"
)

func main() {
	r := gin.Default()
	db := connection.Connection()
	redis := connection.Redis()

	eng := &routers.Routes{
		Db:    db,
		R:     r,
		Redis: redis,
	}

	eng.Routers()

	r.Run()
}
