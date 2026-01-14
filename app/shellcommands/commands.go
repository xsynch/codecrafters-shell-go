package shellcommands

import (
	"fmt"
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

func ProcessInput(input string){
	results := ""
	singleQuotesOpen = false 
	dbQuotesOpen = false
	spacePrinted = false 
	escaped = false 
	
	
	
	cmd := strings.Split(input," ")[0]
	cmdArgs := strings.Replace(input,cmd,"",1)
	cmdArgs = strings.TrimSpace(cmdArgs)
	// runeSlice := []rune(input)

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
		fmt.Println(results)
	

	
	
	
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