package shellcommands

import (
	"fmt"
	"os"
	"path/filepath"
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

func CheckCmdStartWithQuotes(input string) bool {
	if strings.HasPrefix(input,"'") || strings.HasPrefix(input,"\""){
		return true
	}
	return false
}

func PrintEcho(cmd Command){
	originalStdErr := os.Stderr
	originalStdOut := os.Stdout
	
	
	if cmd.StderrRedirect.Redirect {
		var filePerms os.FileMode
		if cmd.StdoutRedirect.Append {
			 filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_APPEND)
		} else {
			filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		}
		if _, err := os.Stat(cmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(cmd.StderrRedirect.RedirectLocation), 0700) // Create your file
		}
		if _, err := os.Stat(cmd.StderrRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(cmd.StderrRedirect.RedirectLocation), 0700) // Create your file
		}
		file,err := os.OpenFile(cmd.StderrRedirect.RedirectLocation,int(filePerms),0644)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error opening the file: %s\n",err.Error())
			
		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		os.Stderr = file 
	} 
	if cmd.StdoutRedirect.Redirect {
		var filePerms os.FileMode
		if cmd.StdoutRedirect.Append {
			 filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_APPEND)
		} else {
			filePerms = os.FileMode(os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
		}
		if _, err := os.Stat(cmd.StdoutRedirect.RedirectLocation); os.IsNotExist(err) {
    			os.MkdirAll(filepath.Dir(cmd.StdoutRedirect.RedirectLocation), 0700) // Create your file
		}
		file,err := os.OpenFile(cmd.StdoutRedirect.RedirectLocation, int(filePerms),0644)
		if err != nil {
			fmt.Fprintf(os.Stderr,"Error opening the file: %s\n",err.Error())
			
		}
		// fmt.Printf("Running %s with options: %s\n",baseExec, test.Args)
		defer file.Close()
		os.Stdout = file 
	} 

	fmt.Printf("%s\n",strings.Join(cmd.Args," "))
	
	// cmd.Stderr = os.Stderr
	os.Stdout = originalStdOut
	os.Stderr = originalStdErr
	
}