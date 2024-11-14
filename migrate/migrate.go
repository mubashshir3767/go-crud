package main

import (
	"github.com/mubashshir/go-crud/initializers"
	"github.com/mubashshir/go-crud/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Post{})
}
