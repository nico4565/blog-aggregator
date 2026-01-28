package main

import (
	"fmt"
	"log"

	"github.com/nico4565/blog-aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Error! Couldn't read the config file: %v\n", err)
	}

	fmt.Printf("config file is:%v\n", c)

	c.SetUser("test_user")
	c, err = config.Read()
	if err != nil {
		log.Fatalf("Error!Couldn't read the config file: %v\n", err)
	}
	fmt.Printf("config file is now:%v\n", c)
}
