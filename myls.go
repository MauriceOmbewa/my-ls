package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Parsing command-line flags
	showAll := flag.Bool("a", false, "Show hidden files")
	detailed := flag.Bool("l", false, "Show detailed information")
	// detailed := flag.Bool("R", false, "Show detailed information")
	// detailed := flag.Bool("r", false, "Show detailed information")
	// detailed := flag.Bool("t", false, "Show detailed information")
	flag.Parse()

	// Get the directory to list (default is current directory)
	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	// Open and read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Loop through the files
	for _, file := range files {
		name := file.Name()

		// Skip hidden files if not showing all
		if !*showAll && name[0] == '.' {
			continue
		}

		// If detailed flag is set, show file info
		if *detailed {
			info, err := file.Info()
			if err != nil {
				fmt.Println("Error getting file info:", err)
				continue
			}
			modTime := info.ModTime().Format(time.RFC822)
			fmt.Printf("%-10s %5d %s %s\n", info.Mode().String(), info.Size(), modTime, name)
		} else {
			fmt.Println(name)
		}
	}
}
