package main

import (
	"flag"
	"fmt"
)

var (
	port    = flag.Int("port", 3306, "Mysql port")
	address = flag.String("address", "localhost", "Mysql address")
	schema  = flag.String("schema", "", "Mysql database")
)

func main() {
	flag.Parse()

	fmt.Printf("mysql port: %d\n", *port)
	fmt.Printf("mysql address: %s\n", *address)
	fmt.Printf("mysql schema: %s\n", *schema)
}
