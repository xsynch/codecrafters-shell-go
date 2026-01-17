package shellcommands

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var Commands = []string{"exit","echo","exit","type"}
var singleQuotesOpen bool 
var dbQuotesOpen bool 
var escaped bool
var spacePrinted bool 
var SpecialChars = []string{"\"","\\","$","`"}

func PreprocessArgs(commandArgs string) string {
	results := ""
	cmdArgs := strings.Split(commandArgs, " ")
	for _, val := range cmdArgs{
		if strings.TrimSpace(val) != ""{
			results += " " + val
		}
	}
	// fmt.Println(strings.TrimSpace(results))
	return results
}

func ProcessInput(input string, redirectionLocation string){
	results := ""
	singleQuotesOpen = false 
	dbQuotesOpen = false
	spacePrinted = false 
	escaped = false 
	
	
	
	cmd := strings.Split(input," ")[0]
	cmdArgs := strings.Replace(input,cmd,"",1)
	cmdArgs = strings.TrimSpace(cmdArgs)
	

	for i :=0; i < len(cmdArgs); i++{
		if cmdArgs[i] == '\\' && !singleQuotesOpen && !dbQuotesOpen{
			escaped = true
			results += string(cmdArgs[i+1])
			i += 1
			continue
		}
		if cmdArgs[i] != ' '	{
			spacePrinted = false
		}
		if cmdArgs[i] == '"' && !dbQuotesOpen && !singleQuotesOpen{
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
		if cmdArgs[i] == '\\' && dbQuotesOpen &&slices.Contains(SpecialChars,string(cmdArgs[i])){
			results += string(cmdArgs[i+1])
			i+=1
			continue
		}
		//if neither quote is open then need to check if a space has already been printed
		//if space has already been printed then don't need to print another space
		if cmdArgs[i] == ' ' && !singleQuotesOpen && !dbQuotesOpen && !spacePrinted{
			// fmt.Println("space found")
			results += string(cmdArgs[i])
			spacePrinted = true
			continue 
		} 
		if cmdArgs[i] == ' ' && spacePrinted{
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
			file,err := os.OpenFile(redirectionLocation, os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0644)
			if err != nil {
				fmt.Fprintf(os.Stderr,"Error opening the file: %s\n",err.Error())
			
			}
			// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
			defer file.Close()
			os.Stdout = file
			fmt.Fprintf(file,"%s\n",results)
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
	
	for i := 0; i < len(input); i++{

		if input[i] == '"' && !dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = true 
			continue 
		}
		if input[i] == '"' && dbQuotesOpen && !singleQuotesOpen{
			// results = append(results, strResults)
			dbQuotesOpen = false 
			continue
		}
		if input[i] == '\'' && !singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = true 
			// strResults += string(input[i])
			continue 
		}
		if input[i] == '\'' && singleQuotesOpen && !dbQuotesOpen{
			singleQuotesOpen = false 
			// results = append(results, strResults)
			continue
		}
		if input[i] == ' ' && !singleQuotesOpen && !dbQuotesOpen {
			// fmt.Println("space found")
			// strResults += string(input[i])
			// spacePrinted = true
			// fmt.Printf("Adding: %s to results\n",strResults)
			results = append(results, strResults)
			strResults = ""
			continue 
		} 
		if input[i] == '\\' && dbQuotesOpen &&slices.Contains(SpecialChars,string(input[i+1])) {
			i += 1
			strResults += string(input[i])
			continue 
		}
		strResults += string(input[i])


		
		

		

	}
	
					
			if strings.TrimSpace(strResults) != "" {
				results = append(results,strResults)
			}
		
		
	
	
	return results
}

func CheckForRedirect(input string) (bool,string,string) {
	singleQuotesOpen = false 
	dbQuotesOpen = false 
	newLine := ""
	// fmt.Println("Checking for redirect in string: ",input)
	
	for i := 0; i < len(input); i++ {
		currentChar := input[i]
		if currentChar == '"' && !dbQuotesOpen && !singleQuotesOpen{
			dbQuotesOpen = true 
			continue 
		}
		if currentChar == '"' && dbQuotesOpen && !singleQuotesOpen{
			dbQuotesOpen = false
			continue 
		}
		if currentChar == '\'' && !dbQuotesOpen && !singleQuotesOpen{
			singleQuotesOpen = true 
			continue 
		}
		if currentChar == '\'' && !dbQuotesOpen && singleQuotesOpen{
			singleQuotesOpen = false 
			continue 
		}
		// if (currentChar == '1'  || currentChar == '>') && !dbQuotesOpen && !singleQuotesOpen{
		if currentChar == '>' && !dbQuotesOpen && !singleQuotesOpen{
			redirectionLocation := ""
			
			
			// if currentChar == '1' && input[i+1] == '>' && input[i-1] != '-'{
			if input[i-1] == '1'{
				redirectionLocation = strings.Split(input[i:],">")[1]
				newLine = strings.Replace(input,input[i-1:],"",-1)
				// fmt.Printf("%s with: %d\n",strings.Split(input[i:],">"),len(strings.Split(input[i:],">")))
				// fmt.Printf("redirecting the output to: %s\n",input[i:])
			} else {
				// fmt.Printf("on char: %v in list: %s with: %d\n",string(input[i]),strings.Split(input[i:],">"),len(strings.Split(input[i:],">")))
				redirectionLocation = strings.Split(input[i:],">")[1]
				newLine = strings.Replace(input,input[i:],"",-1)

			}
			// newLine = strings.Replace(input,input[i:],"",-1)
			redirectionLocation = strings.TrimSpace(redirectionLocation)
			// fmt.Printf("Redirect Location: %s\n",redirectionLocation)
			// fmt.Printf("Trying to remove %s from input\n",input[i:])
			return true, newLine, redirectionLocation

		}
	}


	return false,"","" 
}