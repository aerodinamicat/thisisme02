package main

import "fmt"

func main() {
	port := "12345"
	jwtSecret := "mysecretphrase"
	dbHost := "appdb"
	dbSchema := "thisisme"
	dbUser := "postgres"
	dbPassword := "mysecretpassword"
	dbPort := "5432"

	fmt.Println(port + " " + jwtSecret + " " + dbHost + " " + dbSchema + " " + dbUser + " " + dbPassword + " " + dbPort)
}
