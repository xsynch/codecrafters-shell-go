package helpers 

import (
	"fmt"
	"slices"
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