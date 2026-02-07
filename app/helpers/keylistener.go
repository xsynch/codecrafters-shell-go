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

func (l *BellListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
    // Check if the key pressed is Tab (usually 9)
	cmds := []string{}
    if key == 9 { 
		
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
    return nil, 0, false
}


func SetPaths(completer *readline.PrefixCompleter) {
	childPaths := []readline.PrefixCompleterInterface{}
	childPaths = append(childPaths, completer.GetChildren()...)
	paths := strings.Split(os.Getenv("PATH"),":")
	for _,val := range(paths) {
		// fmt.Printf("%s\n",val)
		// if _, err := os.Stat(val); os.IsNotExist(err) {
      	// 	fmt.Println(val, "does not exist")
	 	// 	 continue 
   		// }
		f, err := os.Open(val)
		if err != nil {
			// fmt.Println(err)

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
				// fmt.Printf("%s with perms%s\n",v.Name(), info.Mode())
			}
		}
	}
}