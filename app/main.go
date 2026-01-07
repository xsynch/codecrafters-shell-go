package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/shellcommands"
)


var _ = fmt.Print


func main() {
	userExit,_ := regexp.Compile("^exit$")
	userEcho,_ := regexp.Compile("^echo ")
	userType,_ := regexp.Compile("^type ")
	
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
			shellcommands.ProcessEcho(line)
			continue 
		}
		if userType.MatchString(line){
			line = checkCommands(line)
			continue 
			// if line == ""{
			// 	continue
			// }

		}
		if !executeProgram(line){
		
			fmt.Printf("%s: command not found\n",line)
		}
	}
	
}


// func processEcho(input string){
// 	data := strings.Split(input, " ")
// 	if len(data) == 1{
// 		fmt.Println("")
// 		return 
// 	}
// 	if len(data) > 1{
// 		fmt.Printf("%s\n",strings.Join(data[1:]," "))
// 		return 
// 	}
// }

func checkCommands(inputs string) string{
	cmd := strings.Split(inputs," ")
	if len(cmd) == 1{
		fmt.Printf("%s requires a command to check\n",cmd[0])
		return ""
	}
	usercmd := cmd[1]
	if slices.Contains(shellcommands.Commands,usercmd){
		fmt.Printf("%s is a shell builtin\n",usercmd)
		
		return ""
	}else {
		if !shellcommands.SearchPath(usercmd){
			fmt.Printf("%s: not found\n",usercmd)
			return usercmd
		}
		return usercmd

	}
	
}

func executeProgram(progName string) bool {
	resultsTest := []string{}
	baseExec := strings.Split(progName, " ")[0]
	// args := strings.Split(progName," ")[1:]
	args2 := strings.Replace(progName,baseExec,"",1)
	args2 = strings.TrimSpace(args2)
	// fmt.Println(args2)
	_, err := exec.LookPath(baseExec)
	if err != nil {
		
		return false 
	}
	if !strings.Contains(args2,"'"){

		resultsTest = strings.Split(args2," ")
	}else {

		for _, val := range strings.Split(args2,"'") {
			if strings.TrimSpace(val) != "" {
				// fmt.Println(val)
				resultsTest = append(resultsTest, fmt.Sprintf("%s",val))
			}
		}
	}
	
	
	// fmt.Printf("results: %s\n",resultsTest)
	// for _,val := range args {
	// 	re := regexp.MustCompile(`'`)
	// 	resultsTest = append(resultsTest,re.ReplaceAllString(val,"\""))
	// }
	// fmt.Printf("new results: %s\n",resultsTest)
	var test exec.Cmd
	test.Args = resultsTest
	cmd := exec.Command(baseExec,test.Args...)
	
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	 return true 
	
}