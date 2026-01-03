package main

import (
	"fmt"
	"log"
	"os"
)


var _ = fmt.Print

func main() {
	
	for {
		fmt.Printf("$ ")
		var userInput string
		n, err := fmt.Scanln(&userInput)
		if err != nil {
			log.Fatal(err)
		}
		if n > 1{
			log.Fatal("Please include only one command right now")
		}
		if userInput == "exit"{
			os.Exit(0)
		}
		fmt.Printf("%s: command not found\n",userInput)
	}
	
}
