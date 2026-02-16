package helpers

import (
	"fmt"
	
	"slices"

	"strings"
)


var TabPressed = 0

type CustomCompleter struct {
	commands [][]rune
}

func (cc *CustomCompleter) Do(line []rune, offset int) ([][]rune, int) {
	
	TabPressed += 1

	
	prefix := string(line[:offset])
	var matches [][]rune

	for _, cmd := range cc.commands {
		
		if strings.HasPrefix(string(cmd), prefix) {			
			matches = append(matches, []rune(cmd))
		}
	}
	
	switch len(matches) {
	case 0:
		fmt.Print("\a")
		return [][]rune{line[offset:]}, offset
	case 1:		
		match := matches[0]
		test := fmt.Sprintf("%s ",string(match[offset:]))
		// return [][]rune{match[offset:]}, offset
		return [][]rune{[]rune(test)}, offset
	default:
		if TabPressed == 1 {
			fmt.Print("\a")			
		}
		if TabPressed == 2 {			
			strMatches := []string{}
			for _,match := range(matches){
				strMatches = append(strMatches, string(match))
			}
			slices.Sort(strMatches)
			fmt.Printf("\r\n%s\r\n",strings.Join(strMatches,"  "))
			TabPressed = 0
		}

		return [][]rune{line[offset:]}, offset
	}
}

func NewCustomPrefixCompleter(defaultCommands [][]rune) *CustomCompleter {

	return &CustomCompleter{commands: defaultCommands,}
}