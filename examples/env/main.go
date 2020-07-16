package main

import (
	"fmt"
	"os"
)

func main() {
	// Setting env variable
	// os.Setenv("MYSQL_PASSWORD", "hakty11")

	// Reading an env variable
	pass := os.Getenv("MYSQL_PASSWORD") // secret
	fmt.Println(pass)
}

// Docker container
// Pass the environemnt while starting the container
