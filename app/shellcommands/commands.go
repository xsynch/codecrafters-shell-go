package shellcommands

import (
	"fmt"
	"strings"
)

var Commands = []string{"exit","echo","exit","type"}
var singleQuotesOpen bool 
var dbQuotesOpen bool 

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
	spacePrinted := false 
	
	
	cmd := strings.Split(input," ")[0]
	cmdArgs := strings.Replace(input,cmd,"",1)
	cmdArgs = strings.TrimSpace(cmdArgs)
	// runeSlice := []rune(input)
	
	for _,val := range cmdArgs {	
		if val != ' '	{
			spacePrinted = false
		}
		if val == '"' && !dbQuotesOpen && !singleQuotesOpen{
			dbQuotesOpen = true 
			// fmt.Println("Starting double Quotes")
			// spacePrinted = false
			continue 
		}
		if val == '\'' && !singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = true
			// fmt.Println("Starting Single Quotes")
			// spacePrinted = false
			continue
		}
		if val == '"' && dbQuotesOpen && !singleQuotesOpen {
			dbQuotesOpen = false
			// fmt.Println("Ending double quotes")
			continue 
		}
		if val == '\'' && singleQuotesOpen && !dbQuotesOpen {
			singleQuotesOpen = false
			// fmt.Println("Ending Single Quotes")
			continue
		}
		//if neither quote is open then need to check if a space has already been printed
		//if space has already been printed then don't need to print another space
		if val == ' ' && !singleQuotesOpen && !dbQuotesOpen && !spacePrinted{
			// fmt.Println("space found")
			results += string(val)
			spacePrinted = true
			continue 
		} 
		if val == ' ' && spacePrinted{
			continue 
		}
		results += string(val)

		
		

	}
	fmt.Println(results)
	
	
}