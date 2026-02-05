package main

import (
	// "bufio"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/shellcommands"
	"github.com/codecrafters-io/shell-starter-go/app/helpers"
	"github.com/chzyer/readline"
)


var _ = fmt.Print

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),	
	readline.PcItem("type"),
	
)


func main() {
	
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		Listener: &helpers.BellListener{Completer: completer},
		// FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()


	userExit,_ := regexp.Compile("^exit$")
	userEcho,_ := regexp.Compile("^echo ")
	userType,_ := regexp.Compile("^type ")
	
	for {
		/*fmt.Print("$ ")
		reader := bufio.NewReader((os.Stdin))
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}*/
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		

		line = strings.Trim(line,"\n")
		originalLine := line 
		
		// userCmd := shellcommands.BuildCommand(line)
		// fmt.Println(userCmd)
		
		// fmt.Printf("Line entered: %s\n",strings.Trim(line,"\n"))

		// var userInput string
		// _, err := fmt.Scanln(&userInput)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		userCmdString := shellcommands.BuildCommand(originalLine)
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
			shellcommands.PrintEcho(userCmdString)
			// shellcommands.ProcessInput(line, redirectLocation)	
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
		if !executeProgram(originalLine,redirectLocation){
		
			fmt.Printf("%s: command not found\n",line)
		}
	}
	
}



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
	// fmt.Printf("Sending over: %s\n",progName)
	userCmd := shellcommands.BuildCommand(progName)
	// fmt.Printf("Looking up command: %s\n",userCmd.Name)
	_, err := exec.LookPath(userCmd.Name)
	if err != nil {
		// fmt.Errorf(err.Error())
		return false 
	}
	// fmt.Println(userCmd)
	test.Args = resultsTest
	cmd := exec.Command(userCmd.Name,userCmd.Args...)
	cmd.Stderr = os.Stderr
	
	
	if userCmd.StderrRedirect.Redirect {
		var filePerms os.FileMode
		if userCmd.StderrRedirect.Append {
			 filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_APPEND)
		} else {
			filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		}
		if _, err := os.Stat(userCmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(userCmd.StderrRedirect.RedirectLocation), 0700) // Create your file
		}
		// file,err := os.OpenFile(userCmd.StderrRedirect.RedirectLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0644)
		file,err := os.OpenFile(userCmd.StderrRedirect.RedirectLocation, int(filePerms),0644)
		if err != nil {
			fmt.Fprintf(cmd.Stderr,"Error opening the file: %s\n",err.Error())
			return true 
		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		cmd.Stderr = file 
	} else {
		cmd.Stderr = os.Stderr		

	}
	if userCmd.StdoutRedirect.Redirect {
		var filePerms os.FileMode
		if userCmd.StdoutRedirect.Append {
			 filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_APPEND)
		} else {
			filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		}
		if _, err := os.Stat(userCmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(userCmd.StdoutRedirect.RedirectLocation), 0700) // Create your file
		}
		if _, err := os.Stat(userCmd.StdoutRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(userCmd.StdoutRedirect.RedirectLocation), 0700) // Create your file
		}
		file,err := os.OpenFile(userCmd.StdoutRedirect.RedirectLocation, int(filePerms),0644)
		if err != nil {
			fmt.Fprintf(cmd.Stderr,"Error opening the file: %s\n",err.Error())
			return true 
		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		cmd.Stdout = file 
	} else {
		cmd.Stdout = os.Stdout

	}
	// cmd.Stderr = os.Stderr
	err = cmd.Run()

	
	return true 

	
	
}

