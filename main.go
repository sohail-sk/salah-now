package main

import (
	"salah-now/db"
	"salah-now/routes"
)

func main() {
	db.InitDB()

	r := routes.SetupRouter()
	r.Run(":7825")
}