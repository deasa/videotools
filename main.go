package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	targetVideoFilesPath = "/Users/brendanashton/Desktop/small but mighty"
	sourceVideoFilesPath = "/Volumes/Untitled/DCIM/100MEDIA"
)

func main() {
	d, err := os.ReadDir(sourceVideoFilesPath)
	if err != nil {
		log.Fatalf("could not read directory: %v", err)
	}

	for _, v := range d {
		info, _ := v.Info()
		if v.Type().IsRegular() && strings.Contains(v.Name(), "DJI_0") && !strings.Contains(info.Name(), strconv.Itoa(time.Now().Year())) {
			fmt.Println("original name: " + v.Name())
			fmt.Println("name to rename to: " + attachTimeToFileName(targetVideoFilesPath, info))
			err := os.Rename(sourceVideoFilesPath+ "/" + v.Name(), attachTimeToFileName(sourceVideoFilesPath, info))
			if err != nil {
				log.Fatalf("error renaming file: %v", err)
			}	
		}
		
		// fmt.Printf("name: '%s' type: '%v' info: '%v'\n", v.Name(), v.Type(), info.ModTime())
	}
}

func attachTimeToFileName(targetDir string, info fs.FileInfo) string {
	return fmt.Sprintf("%s/DJI_%s%s", targetDir, info.ModTime().Local().Format("Mon Jan _2 2006 03-04"), filepath.Ext(info.Name()))
}