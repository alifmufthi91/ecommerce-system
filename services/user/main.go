package main

import "github.com/alifmufthi91/ecommerce-system/services/user/cmd"

// @title 		User Service
// @version 	1.0
// @host 		localhost:8080
// @BasePath 	/api/v1
// @securityDefinitions.apiKey BearerAuth
// @in 			header
// @name 		Authorization
// @schemes 	http https

func main() {
	cmd.Execute()
}
