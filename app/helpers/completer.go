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
		formattedMatch := fmt.Sprintf("%s ",string(match[offset:]))
		
		return [][]rune{[]rune(formattedMatch)}, offset
	default:
		strMatches := []string{}
		for _,match := range(matches){
			strMatches = append(strMatches, string(match))
		}
			slices.Sort(strMatches)
		if TabPressed == 1 {
			fmt.Print("\a")			
			m := lcp(matches)		
			m = []rune(strings.Trim(string(m),string(line)))

			return [][]rune{m},offset
		}
		if TabPressed == 2 {						

			fmt.Printf("\r\n%s\r\n",strings.Join(strMatches,"  "))
			TabPressed = 0
		}

		return [][]rune{line[offset:]}, offset
	}
}

func NewCustomPrefixCompleter(defaultCommands [][]rune) *CustomCompleter {

	return &CustomCompleter{commands: defaultCommands,}
}


func lcp(matches [][]rune) ([]rune){

	for charIndex :=0; charIndex < len(matches[0]); charIndex+=1{
		for stringIndex := 1; stringIndex < len(matches); stringIndex += 1{
			if(len(matches[stringIndex]) <= charIndex) || matches[stringIndex][charIndex] != matches[0][charIndex]{
				// log.Println(string(matches[0][:charIndex]))
				return matches[0][:charIndex]
			}
		}

	}


	return matches[0] 
}