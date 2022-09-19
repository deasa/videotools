package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// targetVideoFilesPath = "/Users/brendanashton/Desktop/small but mighty"
	targetVideoFilesPath = "/Users/brendan.ashton/tmp/tmp2"
	//sourceVideoFilesPath = "/Volumes/Untitled/DCIM/100MEDIA"
	sourceVideoFilesPath = "/Users/brendan.ashton/tmp"
)

type copyFunc func(string, string) (int64, error)

func main() {
	d, err := os.ReadDir(sourceVideoFilesPath)
	if err != nil {
		log.Fatalf("could not read directory: %v", err)
	}
	if len(d) == 0 {
		fmt.Println("no files found in directory: " + sourceVideoFilesPath)
		return
	}

	CopyFiles(d, targetVideoFilesPath, copy)
	RenameFiles(d)
}

func copy(src, dst string) (int64, error) {
	if src == "" || dst == "" {
		return 0, fmt.Errorf("src and dst are both required")
	}

	// if it starts with ".", ignore it.
	if string(filepath.Base(src)[0]) == "." {
		return 0, nil
	}

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyFiles(d []fs.DirEntry, targetPath string, copyFunc copyFunc) {
	if len(d) == 0 || targetPath == "" {
		fmt.Println("directory must not be empty and target path must be specified")
		return
	}

	filepath.WalkDir(filepath.WalkFunc)
	for _, f := range d {
		_, err := copyFunc(f.Name(), path.Join(targetPath, f.Name()))
		if err != nil {
			log.Fatalf("error copying file: %v", err)
		}
	}
}

func RenameFiles(d []fs.DirEntry) {
	for _, v := range d {
		info, _ := v.Info()
		if v.Type().IsRegular() && strings.Contains(v.Name(), "DJI_0") && !strings.Contains(info.Name(), strconv.Itoa(time.Now().Year())) {
			fmt.Println("original name: " + v.Name())
			fmt.Println("name to rename to: " + attachTimeToFileName(targetVideoFilesPath, info))
			err := os.Rename(sourceVideoFilesPath+"/"+v.Name(), attachTimeToFileName(sourceVideoFilesPath, info))
			if err != nil {
				log.Fatalf("error renaming file: %v", err)
			}
		}
	}
}

func attachTimeToFileName(targetDir string, info fs.FileInfo) string {
	return fmt.Sprintf("%s/DJI_%s%s", targetDir, info.ModTime().Local().Format("Mon Jan _2 2006 03-04"), filepath.Ext(info.Name()))
}
