package shellcommands

import (
	"errors"
	"fmt"
	"io/fs"
	// "io/fs"

	"os"
	"path/filepath"
	"slices"
)


func SearchPath(command string) bool {
	path := os.Getenv("PATH")
	var checked []string 
	// paths := strings.Split(path,":")
	paths := filepath.SplitList(path)
	for  _, val := range paths{		
		
		if _,err := os.Stat(val); errors.Is(err, fs.ErrNotExist){
			// fmt.Printf("path: %s\n",val)
			continue
		}
		// errors.Is(err, fs.ErrNotExist).
		// if !fs.ValidPath(val){
		// 	continue
		// }
		info, _ := os.Stat(val)
		// if err != nil {
		// 	// log.Println(err)
		// 	continue
		// }
		if info.IsDir() {
			// fmt.Printf("check path: %s\n",val)
			if !slices.Contains(checked,val){			
				checked = append(checked,val)	
				dirEntries, _ := os.ReadDir(val)
				
				for _,entry := range dirEntries{
					info, _ := entry.Info()
					// fmt.Println(info.Name())
					if info.Name() == command {
						mode := info.Mode().Perm() & 0111
						// mode := 1
						if mode != 0 {	
							fmt.Printf("%s is %s/%s\n",command, val,command)
							return true
						} else {
							continue
						}
					}
				}
				
			}
		}
		

	}
	return false
}