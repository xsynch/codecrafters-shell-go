package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)


var _ = fmt.Print

func main() {
	userExit,_ := regexp.Compile("^exit$")
	userEcho,_ := regexp.Compile("^echo ")
	for {
		fmt.Print("$ ")
		reader := bufio.NewReader((os.Stdin))
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		line = strings.Trim(line,"\n")
		// fmt.Printf("Line entered: %s\n",strings.Trim(line,"\n"))

		// var userInput string
		// _, err := fmt.Scanln(&userInput)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		if userExit.MatchString(line){
			os.Exit(0)
		}
		if userEcho.MatchString(line){
			processEcho(line)
			continue 
		}
		
		fmt.Printf("%s: command not found\n",line)
	}
	
}


func processEcho(input string){
	data := strings.Split(input, " ")
	if len(data) == 1{
		fmt.Println("")
		return 
	}
	if len(data) > 1{
		fmt.Printf("%s\n",strings.Join(data[1:]," "))
		return 
	}
}