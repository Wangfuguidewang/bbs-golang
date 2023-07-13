package main

import (
	"bbs-go/model"
	"bbs-go/routes"
)

func main() {
	model.Initdb()
	routes.InitRouter()
}
