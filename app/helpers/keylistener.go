package helpers

import (
	"fmt"

	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

type BellListener struct {
    Completer *readline.PrefixCompleter
}

type ChangeListener struct {
	Completer *CustomCompleter
}

func (cl *ChangeListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	    if key != 9 { 
			TabPressed = 0
		}
		return nil, 0, false
}

func SetCmdPaths(completer *CustomCompleter) {
	childPaths := [][]rune{}
	childPaths = append(childPaths, completer.commands...)
	paths := strings.Split(os.Getenv("PATH"),":")
	defaultCommands := []string{}
	for _,val := range(completer.commands){
		defaultCommands = append(defaultCommands, string(val))
	}
	for _,val := range(paths) {

		f, err := os.Open(val)
		if err != nil {			

			continue
		}
		files, err := f.ReadDir(0)
		if err != nil {
			
			continue
		}
		for _, v := range files {
			if !v.IsDir(){
				if _, err := os.Stat(val); os.IsNotExist(err) {      		
	 		 		continue 
   				}
			
				_, err := os.Stat(filepath.Join(val,v.Name()))
				if err != nil {
					
					continue
				}
				if !slices.Contains(defaultCommands,v.Name()){
					// log.Printf("adding %s to list",v.Name())
					childPaths = append(childPaths, []rune(v.Name()))								
				}
				
				
			}
		}
	}
	// completer.commands = append(completer.commands, childPaths...)
	completer.commands = childPaths
}



func (l *BellListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
    // Check if the key pressed is Tab (usually 9)
	cmds := []string{}
	tabPressed := 0
    if key == 9 && tabPressed == 0 { 
		tabPressed += 1
		
        lineStr := string(line[:pos])
		
        // Check if the current input has any valid completions
		
		results := l.Completer.GetChildren()
		for _,val := range(results) {
			cmds = append(cmds, string(val.GetName()))
			
		}

		_, offset := l.Completer.Do([]rune(line), pos)


        if !slices.Contains(cmds,lineStr){
			fmt.Printf("\a")
		}
        if offset == 0 {
            // No matches found, trigger the terminal bell
			// fmt.Println("tab pressed with no completion")
            fmt.Print("\a")
        }
    }
	if key == 9 && tabPressed > 1{
		tabPressed += 1
		fmt.Println("Tab pressed for second time")
	}
    return nil, 0, false
}


func SetPaths(completer *readline.PrefixCompleter) {
	childPaths := []readline.PrefixCompleterInterface{}
	childPaths = append(childPaths, completer.GetChildren()...)
	paths := strings.Split(os.Getenv("PATH"),":")
	for _,val := range(paths) {

		f, err := os.Open(val)
		if err != nil {
			

			continue
		}
		files, err := f.ReadDir(0)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		for _, v := range files {
			if !v.IsDir(){
		if _, err := os.Stat(val); os.IsNotExist(err) {
      		// fmt.Println(val, "does not exist")
	 		 continue 
   		}
			
				_, err := os.Stat(filepath.Join(val,v.Name()))
				if err != nil {
					// fmt.Println(err)
					continue
				}
				childPaths = append(childPaths, readline.PcItem(v.Name()))
				completer.SetChildren(childPaths)
				
			}
		}
	}
}