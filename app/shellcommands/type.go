package shellcommands

import (
	
	"fmt"
	
	"os/exec"

	// "os"
	
	
)


func SearchPath(command string) bool {
	// path := os.Getenv("PATH")
	// var checked []string 

	cmdFound, err := exec.LookPath(command)
	if err != nil {
		return false
	} else {
		fmt.Printf("%s is %s\n",command, cmdFound)
		return true 
	}
	
	
	// paths := filepath.SplitList(path)
	// for  _, val := range paths{		
		
	// 	if _,err := os.Stat(val); errors.Is(err, fs.ErrNotExist){
			
	// 		continue
	// 	}

	// 	info, _ := os.Stat(val)

	// 	if info.IsDir() {
			
	// 		if !slices.Contains(checked,val){			
	// 			checked = append(checked,val)	
	// 			dirEntries, _ := os.ReadDir(val)
				
	// 			for _,entry := range dirEntries{
	// 				info, _ := entry.Info()
					
	// 				if info.Name() == command {
	// 					mode := info.Mode().Perm() & 0111
					
	// 					if mode != 0 {	
	// 						fmt.Printf("%s is %s/%s\n",command, val,command)
	// 						return true
	// 					} else {
	// 						continue
	// 					}
	// 				}
	// 			}
				
	// 		}
	// 	}
		

	// }
	// return false
}