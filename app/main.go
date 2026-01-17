package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
		redir,newLine,redirectLocation := shellcommands.CheckForRedirect(line)	
		if redir {
			// fmt.Println("Redirecting stdout to ",location)
			line = newLine
			// fmt.Println("Updated line: ",line)
		}

		if userExit.MatchString(line){
			os.Exit(0)
		}
		if userEcho.MatchString(line){	

			shellcommands.ProcessInput(line, redirectLocation)	
			// shellcommands.ProcessEcho(line)
			continue 
		}
		if userType.MatchString(line){			
			line = checkCommands(line)
			continue 
			// if line == ""{
			// 	continue
			// }

		}
		if !executeProgram(line,redirectLocation){
		
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

func executeProgram(progName string, redirectLocation string) bool {
	resultsTest := []string{}
	var test exec.Cmd

	cmdAndArgs := shellcommands.CmdHelper(progName)
	// fmt.Printf("program: %s and args: %s\n",cmdAndArgs[0],cmdAndArgs[1:])
	// return true 

	// baseExec := strings.Split(progName, " ")[0]
	baseExec := cmdAndArgs[0]
	
	// args2 := strings.Replace(progName,baseExec,"",1)
	// args2 = strings.TrimSpace(args2)

	args2 := cmdAndArgs[1:]
	// args2 = strings.TrimSpace(args2)

	
	_, err := exec.LookPath(baseExec)
	if err != nil {
		// fmt.Errorf(err.Error())
		return false 
	}
	// fmt.Printf("Executing program: %s\n",progName)
	// args2 = shellcommands.PreprocessArgs(args2)
	
	if shellcommands.StringHasQuotes(fmt.Sprintf("%s",args2)){		
		

		resultsTest = shellcommands.CmdHelper(fmt.Sprintf("%s",args2))


	} else {
		resultsTest = args2// strings.Split(args2," ")
	}

	
	test.Args = resultsTest
	cmd := exec.Command(baseExec,test.Args...)
	cmd.Stderr = os.Stderr
	
	if len(redirectLocation) == 0 {
		cmd.Stdout = os.Stdout
	} else {
		
		if _, err := os.Stat(redirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(redirectLocation), 0700) // Create your file
		}
		file,err := os.OpenFile(redirectLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0644)
		if err != nil {
			fmt.Fprintf(cmd.Stderr,"Error opening the file: %s\n",err.Error())
			return true 
		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		cmd.Stdout = file 

	}
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	// if err != nil {
	// 	fmt.Fprintf(cmd.Stderr,"Command Error: %s\n",err)
	// }

	 return true 
	
}