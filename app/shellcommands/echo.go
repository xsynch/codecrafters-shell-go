package shellcommands

import (
	"fmt"
	"regexp"
	"strings"
)


func ProcessEcho(input string){
	data := strings.Split(input, " ")
	if len(data) == 1{
		fmt.Println("")
		return 
	}
	strippedData := HandleSingleQuotes(strings.Join(data[1:]," "))
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