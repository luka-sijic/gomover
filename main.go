package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/fatih/color"
)

// File struct - store files in a slice
type filestore struct {
	directory string
	destination string
	filepaths []string
	flagList []flags
	include []string
	exclude []string
	unsafe bool
	hide bool
}

type flags struct {
	flagName string 
	flagValues []string
}

// List of valid arguments/flags
func argList() []string {
	return []string{"-e", "-i", "-u", "-x"}
}

// Check if an argument contains a valid flag
func contains(arg string, list []string) bool {
	for i := range list {
		if strings.Contains(arg, list[i]) == true {
			return true
		} 
	}
	return false
}

// Check for dupe
func checkIncludes(excludes []string, value string) bool {
	for i := range excludes {
		if value == excludes[i] {
			return true
		}
	}
	return false
}

func main() {
	f := filestore{}
	if len(os.Args) == 1 {
		fmt.Println("gomover: try 'gomover -h' for more information")
		return
	} else if os.Args[1] == "-h" {
		fmt.Println("Usage: gomover [directory source] [destination directory] [option] [string]...")
		fmt.Println(" -e, --exclude <string>...<string> Excludes all files that contain the provided strings")
		fmt.Println(" -i, --include <string>...<string> Includes all files that contain the provided strings")
		fmt.Println(" -x, --hide                        Hides all excluded files for use in large directories")
		fmt.Println(" -u, --unsafe                      Does not ask before moving files, be cautious when using this flag")
		return
	} 
	options(&f)
	//fmt.Printf("Length of os.arguments %d\n", len(os.Args))
	color.Set(color.BgCyan)
	fmt.Printf("Directory: %s\n", f.directory)
	color.Unset()
	color.Set(color.BgBlue)
	fmt.Printf("Directory: %s\n", f.destination)
	color.Unset()

	crawl(&f)

	fmt.Println(f.flagList)

	for i := range f.flagList {
		switch f.flagList[i].flagName {
		case "-e":
			f.exclude = move(&f, f.flagList[i].flagValues)
		case "-i":
			f.include = move(&f, f.flagList[i].flagValues)
		case "-x":
			f.hide = true
		case "-u":
			f.unsafe = true
		}
	}

	if len(f.include) == 0 {
		fmt.Println("No files found with the given flags")
		return
	}

	fmt.Println("Files to Include: ")
	color.Set(color.FgGreen)
	for i := range f.include {
		fmt.Println(f.include[i])
	}
	color.Unset()
	if f.hide != true {
		fmt.Println("Files to Exclude: ")
		color.Set(color.FgRed)
		for i := range f.exclude {
			fmt.Println(f.exclude[i])
		}
		color.Unset()
	}
	
	var final []string 
	for i := range f.include {
		if checkIncludes(f.exclude, f.include[i]) == false {
			final = append(final, f.include[i])
		}
	}

	fmt.Println("Final List of Files: ")
	color.Set(color.FgCyan)
	for i := range final {
		fmt.Println(final[i])
	}
	color.Unset()

	if f.unsafe == false {
		var option string
		fmt.Print("Would you like to move these files(y/n)?\n> ")
		fmt.Scanln(&option)
		if option == "y" {
			execute(&f, final)
		} else {
			fmt.Println("User cancelled move")
			return
		}
	} else {
		execute(&f, final)
	}

	

}

func execute(f *filestore, final []string) {
	for i := range final {
		cmd := exec.Command("mv", final[i], f.destination)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func move(f *filestore, v []string) []string {
	var temp []string
	for i := range f.filepaths {
		if contains(f.filepaths[i], v) {
			temp = append(temp, f.filepaths[i])
		} 
	}
	return temp
}

func crawl(f *filestore) {
	err := filepath.Walk(f.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() == false {
			f.filepaths = append(f.filepaths, path)
		}
		//fmt.Printf("dir: %v: name %s\n", info.IsDir(), path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func options(f *filestore) {
	f.directory = os.Args[1]
	f.destination = os.Args[2]
	for i := 3; i < len(os.Args);i++ {
		//fmt.Printf("main for loop index: %d\n", i)
		if contains(os.Args[i], argList()) {
			flag := flags{}
			flag.flagName = os.Args[i]
			temp := make([]string, 0)
			j := i
			for y := i; y < len(os.Args);y++ {
				//fmt.Printf("J condition: %d\n", j)
				if j < len(os.Args)-1 {
					if contains(os.Args[j+1], argList()) == false {
						temp = append(temp, os.Args[y])
						//fmt.Println("TO ADD: " + os.Args[y])
						j++
					} else {
						temp = append(temp, os.Args[y])
						//fmt.Println(os.Args[y])
						break
					}
				} else {
					temp = append(temp, os.Args[y])
				}
			}
			flag.flagValues = temp[1:]
			f.flagList = append(f.flagList, flag)
		} 
	} 
}


/*
TODO:
	binary <dir to crawl> <dir to move files to> <flag> <string> ...
	options:
		-e string1 string2 string3 <SHOWS ALL FILES THAT DO NOT HAVE THIS STRING IN THEIR NAME>
		-i extension1 extension2 extension3 <SHOWS ALL FILES THAT DO CONTAIN THIS STRING IN THEIR NAME>
		-u unsafe, does not ask before the program moves the files

	possibilities:
		binary ./ -e !q -i mp4 mkv

	additions:
		faster searching methods
*/