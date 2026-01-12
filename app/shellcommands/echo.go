package shellcommands

import (
	"fmt"
	"regexp"
	"strings"
)


func ProcessEcho(input string){

	
	cmdEcho := strings.Split(input," ")[0]
	// data := strings.Split(input, " ")
	data := strings.Replace(input,cmdEcho,"",1)
	// fmt.Println(data)
	
	data = strings.TrimSpace(data)
	// fmt.Println(data)
	// fmt.Println(data)
	var newData string 
	if len(data) == 1{
		fmt.Println("")
		return 
	}
	
	// cmdArgs := strings.Join(data[1:]," ")
	cmdArgs := data 
	if !StringHasQuotes(data){
		
		output := PrintStringWithNoQuotes(cmdArgs)
		fmt.Printf("%s\n",output)
		return
	}
	// cmdArgs = PreprocessArgs(data)
	cmdArgs = ""
	for _,val := range strings.Split(data,"\""){
		if strings.TrimSpace(val) != ""{
			cmdArgs += fmt.Sprintf("\"%s\"",val)
			// fmt.Println(val)
		}
	}
	cmdArgs = strings.TrimSpace(cmdArgs)
	
	// fmt.Printf("%s\n",cmdArgs)

	// redb := regexp.MustCompile(`^"*"`)
	strwdb := strings.HasPrefix(cmdArgs,"\"")
	if strwdb{
		// fmt.Printf("Processing the double quotes on :%s\n",cmdArgs)
		newData = HandleDoubleQuotes(cmdArgs)
		newData = strings.TrimSpace(newData)
		fmt.Printf("%s\n",newData)
		return 
	}
	// fmt.Println("Stripped of double Quotes: ",newData)
	strippedData := HandleSingleQuotes(cmdArgs)
	if len(data) > 1{
		// fmt.Printf("%s\n",strings.Join(data[1:]," "))
		fmt.Printf("%s\n",strippedData)
		return 
	}
}

func HandleSingleQuotes(input string) string {
	if !strings.Contains(input,"'"){
		re := regexp.MustCompile(`\s+`)
		return re.ReplaceAllString(input," ")
	}

	removedDoubleQuotes := strings.ReplaceAll(input,"''","")
	removeSingle := strings.Trim(removedDoubleQuotes,"'")

	re := regexp.MustCompile(`'`)
	removeSingle = re.ReplaceAllString(removeSingle,"")

	return removeSingle 
}

func HandleDoubleQuotes(input string) string {
	removeQuotes := strings.ReplaceAll(input,"\"\"","")
	re := regexp.MustCompile(`"`)
	return re.ReplaceAllString(removeQuotes,"")
}

func StringHasQuotes(input string) bool {
	if strings.HasPrefix(input,"'") || strings.HasPrefix(input,"\""){
		return true
	}
	return false
}

func PrintStringWithNoQuotes(input string) string {
	re := regexp.MustCompile(`\s+`)
	results := re.ReplaceAllString(input," ")
	return results
}


