package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	search  = flag.String("s", "", "Search regex pattern")
	replace = flag.String("r", "", "Replacement string, $ for submatch,$1 represents the text of the first submatch")
	preview = flag.Bool("p", false, "Preview mod, do not rename")
	dir     = flag.String("d", ".", "Working directory")
)

func main() {
	flag.Parse()
	if search == nil || len(*search) == 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	searchReg, err := regexp.Compile(*search)
	if err != nil {
		log.Fatal(err)
	}

	replaceStr := ""
	if replace != nil {
		replaceStr = *replace
	}

	previewMod := false
	if preview != nil {
		previewMod = *preview
	}

	var workingDir string
	if dir != nil {
		workingDir, err = filepath.Abs(*dir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		workingDir, err = filepath.Abs(".")
		if err != nil {
			log.Fatal(err)
		}
	}

	fis, err := ioutil.ReadDir(workingDir)
	if err != nil {
		log.Fatal(err)
	}

	lenFis := len(fis)
	for i, f := range fis {
		fName := f.Name()
		newName := searchReg.ReplaceAllString(fName, replaceStr)
		if fName == newName {
			fmt.Printf("(%d/%d) No change [%s]\n", i+1, lenFis, fName)
			continue
		}
		if previewMod {
			fmt.Printf("(%d/%d) Will rename [%s] to [%s]\n", i+1, lenFis, fName, newName)
			continue
		}
		fmt.Printf("(%d/%d) Rename [%s] to [%s]\n", i+1, lenFis, fName, newName)
		err := os.Rename(filepath.Join(workingDir, fName), filepath.Join(workingDir, newName))
		if err != nil {
			log.Fatal(err)
		}
	}
}
