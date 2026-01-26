package shellcommands

import (
	
	"fmt"
	
	"os/exec"

	
	
	
)


func SearchPath(command string) bool {


	cmdFound, err := exec.LookPath(command)
	if err != nil {
		return false
	} else {
		fmt.Printf("%s is %s\n",command, cmdFound)
		return true 
	}
	
	
	
}