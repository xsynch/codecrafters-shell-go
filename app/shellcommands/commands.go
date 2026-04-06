package shellcommands

import (
	"fmt"
	"io"

	"log"

	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type StdOutRedirect struct {
	Redirect         bool
	Append           bool
	RedirectLocation string
}

type StdErrRedirect struct {
	Redirect         bool
	Append           bool
	RedirectLocation string
}

type Command struct {
	Name           string
	Args           []string
	StdoutRedirect StdOutRedirect
	StderrRedirect StdErrRedirect
}

var Commands = []string{"exit", "echo", "exit", "type"}
var singleQuotesOpen bool
var dbQuotesOpen bool
var escaped bool
var spacePrinted bool
var SpecialChars = []string{"\"", "\\", "$", "`", "\n"}

var _ = log.Printf

func NewCommand() Command {
	return Command{}
}

func PreprocessArgs(commandArgs string) string {
	results := ""
	cmdArgs := strings.Split(commandArgs, " ")
	for _, val := range cmdArgs {
		if strings.TrimSpace(val) != "" {
			results += " " + val
		}
	}
	// fmt.Println(strings.TrimSpace(results))
	return results
}

func ProcessInput(input string, redirectionLocation string) {
	results := ""
	singleQuotesOpen = false
	dbQuotesOpen = false
	spacePrinted = false
	escaped = false

	cmd := strings.Split(input, " ")[0]
	cmdArgs := strings.Replace(input, cmd, "", 1)
	cmdArgs = strings.TrimSpace(cmdArgs)

	for i := 0; i < len(cmdArgs); i++ {
		if cmdArgs[i] == '\\' && !singleQuotesOpen && !dbQuotesOpen {
			escaped = true
			results += string(cmdArgs[i+1])
			i += 1
			continue
		}
		if cmdArgs[i] != ' ' {
			spacePrinted = false
		}
		if cmdArgs[i] == '"' && !dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = true

			// fmt.Println("Starting double Quotes")
			// spacePrinted = false
			continue
		}
		if cmdArgs[i] == '\'' && !singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = true
			// fmt.Println("Starting Single Quotes")
			// spacePrinted = false
			continue
		}
		if cmdArgs[i] == '"' && dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = false
			// fmt.Println("Ending double quotes")
			continue
		}
		if cmdArgs[i] == '\'' && singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = false
			// fmt.Println("Ending Single Quotes")
			continue
		}
		if cmdArgs[i] == '\\' && dbQuotesOpen && slices.Contains(SpecialChars, string(cmdArgs[i])) {
			results += string(cmdArgs[i+1])
			i += 1
			continue
		}
		//if neither quote is open then need to check if a space has already been printed
		//if space has already been printed then don't need to print another space
		if cmdArgs[i] == ' ' && !singleQuotesOpen && !dbQuotesOpen && !spacePrinted {
			// fmt.Println("space found")
			results += string(cmdArgs[i])
			spacePrinted = true
			continue
		}
		if cmdArgs[i] == ' ' && spacePrinted {
			continue
		}
		results += string(cmdArgs[i])

	}
	if len(redirectionLocation) == 0 {
		fmt.Println(results)
	} else {
		if _, err := os.Stat(redirectionLocation); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(redirectionLocation), 0700) // Create your file
		}

		originalStdout := os.Stdout
		file, err := os.OpenFile(redirectionLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening the file: %s\n", err.Error())

		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		os.Stdout = file
		fmt.Fprintf(file, "%s\n", results)
		os.Stdout = originalStdout

	}

}

func CmdHelper(input string) []string {
	strResults := ""
	results := []string{}
	singleQuotesOpen = false
	dbQuotesOpen = false
	spacePrinted = false
	// fmt.Printf("input: %s with length: %d\n",input,len(input))

	for i := 0; i < len(input); i++ {

		if input[i] == '"' && !dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = true
			continue
		}
		if input[i] == '"' && dbQuotesOpen && !singleQuotesOpen {
			// results = append(results, strResults)
			dbQuotesOpen = false
			continue
		}
		if input[i] == '\'' && !singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = true
			// strResults += string(input[i])
			continue
		}
		if input[i] == '\'' && singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = false
			// results = append(results, strResults)
			continue
		}
		if input[i] == ' ' && !singleQuotesOpen && !dbQuotesOpen {
			// fmt.Println("space found")

			// fmt.Printf("Adding: %s to results\n",strResults)
			results = append(results, strResults)
			strResults = ""
			continue
		}
		if input[i] == '\\' && dbQuotesOpen && slices.Contains(SpecialChars, string(input[i+1])) {
			i += 1
			strResults += string(input[i])
			continue
		}
		strResults += string(input[i])

	}

	if strings.TrimSpace(strResults) != "" {
		results = append(results, strResults)
	}

	return results
}

func CheckForRedirect(input string) (bool, string, string) {
	singleQuotesOpen = false
	dbQuotesOpen = false
	newLine := ""

	// fmt.Println("Checking for redirect in string: ",input)

	for i := 0; i < len(input); i++ {
		currentChar := input[i]
		if currentChar == '"' && !dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = true
			continue
		}
		if currentChar == '"' && dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = false
			continue
		}
		if currentChar == '\'' && !dbQuotesOpen && !singleQuotesOpen {
			singleQuotesOpen = true
			continue
		}
		if currentChar == '\'' && !dbQuotesOpen && singleQuotesOpen {
			singleQuotesOpen = false
			continue
		}
		// if (currentChar == '1'  || currentChar == '>') && !dbQuotesOpen && !singleQuotesOpen{
		if currentChar == '>' && !dbQuotesOpen && !singleQuotesOpen {
			redirectionLocation := ""

			// if currentChar == '1' && input[i+1] == '>' && input[i-1] != '-'{
			if input[i-1] == '1' {
				redirectionLocation = strings.Split(input[i:], ">")[1]
				newLine = strings.Replace(input, input[i-1:], "", -1)
				// fmt.Printf("%s with: %d\n",strings.Split(input[i:],">"),len(strings.Split(input[i:],">")))
				// fmt.Printf("redirecting the output to: %s\n",input[i:])
			} else {
				// fmt.Printf("on char: %v in list: %s with: %d\n",string(input[i]),strings.Split(input[i:],">"),len(strings.Split(input[i:],">")))
				redirectionLocation = strings.Split(input[i:], ">")[1]
				newLine = strings.Replace(input, input[i:], "", -1)

			}
			if input[i-1] == '2' {

			}
			// newLine = strings.Replace(input,input[i:],"",-1)
			redirectionLocation = strings.TrimSpace(redirectionLocation)
			// fmt.Printf("Redirect Location: %s\n",redirectionLocation)
			// fmt.Printf("Trying to remove %s from input\n",input[i:])
			return true, newLine, redirectionLocation

		}
	}

	return false, "", ""
}

func BuildCommand(cmdLine string) (Command, string) {
	moreCommands := ""
	fullCmd := NewCommand()
	baseCmd := ""
	escaped = false
	redirectOpen := false
	_ = redirectOpen

	errRedirectOpen := false
	_ = errRedirectOpen
	stdoutRedirectOpen := false
	_ = stdoutRedirectOpen

	fullCmd.StderrRedirect.Redirect = false
	fullCmd.StdoutRedirect.Redirect = false
	cmdLine = strings.TrimSpace(cmdLine)

	// if CheckCmdStartWithQuotes(cmdLine){
	strResults := ""
	results := []string{}
	singleQuotesOpen = false
	dbQuotesOpen = false
	spacePrinted = false

	// newLine := ""
	// fmt.Printf("%s is the sent command line:\n",cmdLine)

	for i := 0; i < len(cmdLine); i++ {

		currentChar := cmdLine[i]
		// fmt.Println(string(currentChar))
		if escaped {

			strResults += string(currentChar)

			// fmt.Printf("Current stResults: %s\n",strResults)
			escaped = false
			continue
		}
		if cmdLine[i] == '"' && !dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = true
			continue
		}
		if cmdLine[i] == '"' && dbQuotesOpen && !singleQuotesOpen {

			dbQuotesOpen = false
			continue
		}
		if cmdLine[i] == '\'' && !singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = true
			// strResults += string(input[i])
			continue
		}
		if cmdLine[i] == '\'' && singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = false

			continue
		}
		if currentChar == '>' && !dbQuotesOpen && !singleQuotesOpen {
			// redirectionLocation := ""

			switch cmdLine[i-1] {
			case '1':
				fullCmd.StdoutRedirect.Redirect = true
				// fmt.Printf("%s is now in results\n",strResults)
				strResults = strResults[:len(strResults)-1]
				stdoutRedirectOpen = true
				// fmt.Printf("%s is now in results\n",strResults)

			case '2':
				// fmt.Printf("in stderr case\n")
				errRedirectOpen = true
				fullCmd.StderrRedirect.Redirect = true
				if len(strResults)-1 > 0 {
					strResults = strResults[:len(strResults)-1]
				} else {
					strResults = ""
				}
			case '>':
				if fullCmd.StdoutRedirect.Redirect {
					fullCmd.StdoutRedirect.Append = true
				} else if fullCmd.StderrRedirect.Redirect {
					fullCmd.StderrRedirect.Append = true
				}
			default:
				fullCmd.StdoutRedirect.Redirect = true
				stdoutRedirectOpen = true
			}

			continue
			// redirectionLocation = strings.TrimSpace(redirectionLocation)

		}

		if cmdLine[i] == ' ' && !singleQuotesOpen && !dbQuotesOpen {
			// fmt.Printf("Current args: %s with length: %d\n",fullCmd.Args,len(fullCmd.Args))
			if len(baseCmd) == 0 {
				// fmt.Printf("full cmd: %s\n",strResults)
				// log.Printf("Setting command name: %s\n",strResults)
				baseCmd = strResults
				fullCmd.Name = baseCmd

			} else if fullCmd.StdoutRedirect.Redirect && stdoutRedirectOpen {
				// fmt.Printf("adding to stdout: %s\n",strResults)
				if len(strResults) > 0 {
					fullCmd.StdoutRedirect.RedirectLocation = strResults
					stdoutRedirectOpen = false
				}

			} else if fullCmd.StderrRedirect.Redirect && errRedirectOpen {
				// fmt.Printf("adding to stderr: %s\n",strResults)
				if len(strResults) > 0 {
					fullCmd.StdoutRedirect.RedirectLocation = strResults
					errRedirectOpen = false
				}

			} else {

				// if strResults != " "{
				// strResults = strings.TrimSpace(strResults)
				// fmt.Printf("Adding %s to strResults\n",strResults)
				if len(strResults) > 0 {
					fullCmd.Args = append(fullCmd.Args, strResults)
					results = append(results, strResults)
					redirectOpen = false
				}
				// }

			}
			strResults = ""
			// fmt.Printf("Current args: %s with length: %d\n",fullCmd.Args,len(fullCmd.Args))
			continue
		}

		if cmdLine[i] == '\\' && !singleQuotesOpen {
			if dbQuotesOpen && slices.Contains(SpecialChars, string(cmdLine[i+1])) {
				escaped = true
				// continue
			} else if dbQuotesOpen && !slices.Contains(SpecialChars, string(cmdLine[i+1])) {
				escaped = false
				strResults += string(currentChar)
			} else if !dbQuotesOpen {
				escaped = true
			}

			continue
		}
		if cmdLine[i] == '|' && !singleQuotesOpen && !dbQuotesOpen {
			// fmt.Printf("We may have a new command: %s\n",strings.Trim(cmdLine[i+1:],""))
			// log.Printf("fullcmd after | found: %v and strResults: %v\n",fullCmd,strResults)
			moreCommands = strings.TrimSpace(cmdLine[i+1:])
			// BuildCommand(strings.Trim(cmdLine[i+1:],""))
			break
		}

		strResults += string(cmdLine[i])

		// }

	}
	if strings.TrimSpace(strResults) != "" && fullCmd.Name != "" {
		if errRedirectOpen {
			fullCmd.StderrRedirect.RedirectLocation = strResults
		} else if stdoutRedirectOpen {
			fullCmd.StdoutRedirect.RedirectLocation = strResults
		} else {
			// fmt.Printf("Left over string: %s\n",strResults)
			// fmt.Printf("Last args: %s with length: %d\n",fullCmd.Args,len(fullCmd.Args))
			fullCmd.Args = append(fullCmd.Args, strResults)
			results = append(results, strResults)
		}
	} else if strings.TrimSpace(strResults) != "" && fullCmd.Name == "" {
		fullCmd.Name = strResults
	}

	// log.Printf("more commands: %s\n",moreCommands)
	// log.Printf("Command being returned: %v\n",fullCmd)
	return fullCmd, moreCommands

}

func RunCommands(commands []Command) {
	allCmds := make([]*exec.Cmd, 0, len(commands))
	// echoResults := ""
	// typeResults := ""
	// builtinMap := map[int]string{}

	for _, userCmd := range commands {
		// log.Printf("Command Name: %s\n",userCmd.Name)
		cmd := exec.Command(userCmd.Name, userCmd.Args...)


		cmd.Stderr = os.Stderr

		if userCmd.StderrRedirect.Redirect {
			var filePerms os.FileMode
			if userCmd.StderrRedirect.Append {
				filePerms = os.FileMode(os.O_WRONLY | os.O_CREATE | os.O_APPEND)
			} else {
				filePerms = os.FileMode(os.O_WRONLY | os.O_CREATE | os.O_TRUNC)
			}
			if _, err := os.Stat(userCmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(userCmd.StderrRedirect.RedirectLocation), 0700)
			}

			file, err := os.OpenFile(userCmd.StderrRedirect.RedirectLocation, int(filePerms), 0644)
			if err != nil {
				fmt.Fprintf(cmd.Stderr, "Error opening the file: %s\n", err.Error())
				return
			}
			// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
			defer file.Close()
			cmd.Stderr = file
		} else {
			cmd.Stderr = os.Stderr

		}
		// if idx != len(commands)-1 {
		if userCmd.StdoutRedirect.Redirect {
			var filePerms os.FileMode
			if userCmd.StdoutRedirect.Append {
				filePerms = os.FileMode(os.O_WRONLY | os.O_CREATE | os.O_APPEND)
			} else {
				filePerms = os.FileMode(os.O_WRONLY | os.O_CREATE | os.O_TRUNC)
			}
			if _, err := os.Stat(userCmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(userCmd.StdoutRedirect.RedirectLocation), 0700) // Create your file
			}
			if _, err := os.Stat(userCmd.StdoutRedirect.RedirectLocation); os.IsNotExist(err) {
				os.MkdirAll(filepath.Dir(userCmd.StdoutRedirect.RedirectLocation), 0700) // Create your file
			}
			file, err := os.OpenFile(userCmd.StdoutRedirect.RedirectLocation, int(filePerms), 0644)
			if err != nil {
				fmt.Fprintf(cmd.Stderr, "Error opening the file: %s\n", err.Error())
				return
			}
			// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
			defer file.Close()
			cmd.Stdout = file
		} else {
			cmd.Stdout = os.Stdout
		}

		// log.Printf("Adding %s to list with args: %s\n", cmd.Path, cmd.Args)
		allCmds = append(allCmds, cmd)
		// _ = cmd.Run()
	}
	// log.Printf("all commands: %v\n",allCmds)
	for i := 0; i < len(commands)-1; i++ {
		
		r, w, _ := os.Pipe()		
		// log.Printf("Setting %s to w\n", allCmds[i].Path)
		if commands[i].Name == "echo" {
			if _, err := io.WriteString(w,ReturnEcho(commands[i])); err != nil {
				log.Printf("%s\n",err.Error())
				return 
			}
			allCmds[i+1].Stdin = r
			w.Close()
			continue
			// log.Printf("%s",ReturnEcho(commands[i]))
		} else if commands[i].Name == "type"{
			resultsType, _ := ReturnSearchPath(commands[i].Args[0])
			// log.Printf("value from return %s",resultsType)
			if _, err := io.WriteString(w,resultsType); err != nil {
				log.Printf("%s\n",err.Error())
				return 
			}	
			allCmds[i+1].Stdin = r
			w.Close()
			// r.Close()
			continue		
			// log.Printf("%s",resultsType)
			

		} else {
			allCmds[i].Stdout = w
		}
			allCmds[i+1].Stdin = r
		// }

	    if err := allCmds[i].Start(); err != nil {
        	w.Close()
            r.Close()
        }
		w.Close()
		// r.Close()
		
	}
	switch commands[len(commands) - 1].Name {
	case "echo":
		PrintEcho(commands[len(commands) - 1])
		return 
	case "type":
		path, found := ReturnSearchPath(commands[len(commands) - 1].Args[0])
		if !found{
			log.Printf("Command not found\n")		
		} else {
			fmt.Printf("%s\n",path)
		}
		return 
	default:
		allCmds[len(commands)-1].Stdout = os.Stdout
		// log.Printf("stdout: %v\n",allCmds[len(commands)-1].Stdout)
		err := allCmds[len(allCmds)-1].Start()
		if err != nil {
        	log.Printf("Error starting the command: %v\n",err)
    	}

		
	}




	for i := len(allCmds) - 1; i >= 0; i-- {
		if err := allCmds[i].Wait(); err != nil {
			fmt.Printf("waiting for %s: %s", allCmds[i].Path, err.Error())
		}

	}


}


