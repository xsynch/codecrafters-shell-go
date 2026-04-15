package shellcommands

import (
	"bufio"
	"log"

	"strconv"

	"fmt"
	"os"

	"github.com/chzyer/readline"
)

var _ = log.Printf 

func ReturnHistory(shell *readline.Instance,limit ...int) (string,error){
	count := 1
	histData, err := os.Open(shell.Config.HistoryFile)	
	if err != nil {
		return "", err 
	}
	defer histData.Close()
	histScanner := bufio.NewScanner(histData)
	for histScanner.Scan(){
		fmt.Printf("    %d %s\n",count, histScanner.Text())
		count++
	}
	if err := histScanner.Err(); err != nil {
		return "",err
	}
	return "",nil

}

func PrintHistory(cmd Command, shell *readline.Instance) (string,error){
	var histSize int64 
	var err error 
	if len(cmd.Args) > 0 {
		histSize, err = strconv.ParseInt(cmd.Args[0],10,64)
		if err != nil {
			return "",err 
		}
		// log.Printf("history size: %d\n",histSize)
	}
	count := 1
	histData, err := os.Open(shell.Config.HistoryFile)	
	if err != nil {
		return "", err 
	}
	defer histData.Close()
	histScanner := bufio.NewScanner(histData)
	// log.Printf("history size: %d\n",histSize)
	for histScanner.Scan(){
		// log.Printf("count: %d\n",count)
		
		if histSize !=0 && count <= int(histSize) {
			fmt.Printf("    %d %s\n",count, histScanner.Text())
			count++
			
		} else if histSize == 0{
			fmt.Printf("    %d %s\n",count, histScanner.Text())
			count++
		} else {
			break
		} //else {
		// 	fmt.Printf("    %d %s\n",count, histScanner.Text())
			
		// }
		// count++
	}
	if err := histScanner.Err(); err != nil {
		return "",err
	}
	return "",nil

}