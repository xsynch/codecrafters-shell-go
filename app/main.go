package main

import (
	"fmt"
	"log"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	fmt.Printf("$ ")
	var userInput string
	n, err := fmt.Scanln(&userInput)
	if err != nil {
		log.Fatal(err)
	}
	if n > 1{
		log.Fatal("Please include only one command right now")
	}
	fmt.Printf("%s: command not found\n",userInput)
	
}
