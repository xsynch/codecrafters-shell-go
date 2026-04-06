package shellcommands

import (
	"fmt"
	"slices"

	"os/exec"
)

func SearchPath(command string) bool {

	cmdFound, err := exec.LookPath(command)
	if err != nil {
		return false
	} else {
		fmt.Printf("%s is %s\n", command, cmdFound)
		return true
	}

}

func ReturnSearchPath(command string) (string, bool) {
	if slices.Contains(Commands, command){
		return fmt.Sprintf("%s is a shell builtin",command),true 
	}
	cmdFound, err := exec.LookPath(command)
	if err != nil {
		return "", false
	} else {
		results := fmt.Sprintf("%s is %s\n", command, cmdFound)
		return results,true 
	}

}
