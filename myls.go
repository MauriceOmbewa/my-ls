package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"
)

// Function to list directory contents
func listDir(path string, showAll, detailed, reverse, sortByTime, recursive bool, isRoot bool) {
	// Read the directory contents
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Sort the files alphabetically by default
	sort.Slice(files, func(i, j int) bool {
		// If sorting by time, use modification time
		if sortByTime {
			info1, _ := files[i].Info()
			info2, _ := files[j].Info()
			return info1.ModTime().After(info2.ModTime())
		}
		return files[i].Name() < files[j].Name()
	})

	// Reverse the order if `-r` flag is set (for files)
	if reverse {
		for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
			files[i], files[j] = files[j], files[i]
		}
	}

	// Print ".:" at the beginning if this is the root and recursive (-R) is true
	if isRoot && recursive {
		fmt.Println(".:")
	}

	// Print the directory name if not the root level (in recursive calls)
	if !isRoot {
		fmt.Printf("\n%s:\n", path)
	}

	// Prepare a list of directories to be processed recursively
	var dirsToRecurse []string

	// Print the files
	for _, file := range files {
		name := file.Name()

		// Skip hidden files if not showing all (i.e., no -a flag)
		if !showAll && name[0] == '.' {
			continue
		}

		// If detailed flag `-l` is set, show file info
		if detailed {
			info, err := file.Info()
			if err != nil {
				fmt.Println("Error getting file info:", err)
				continue
			}
			modTime := info.ModTime().Format(time.RFC822)
			fmt.Printf("%-10s %5d %s %s\n", info.Mode().String(), info.Size(), modTime, name)
		} else {
			// Print file names horizontally with a single space between them
			fmt.Print(name + "  ")
		}

		// Collect directories for recursive listing (do not recurse into hidden directories unless -a is set)
		if recursive && file.IsDir() && (showAll || name[0] != '.') {
			dirsToRecurse = append(dirsToRecurse, name)
		}
	}

	// Print newline after listing files (if not in detailed mode)
	if !detailed {
		fmt.Println()
	}

	// Reverse the directories to recurse into if `-r` is set
	if reverse {
		for i, j := 0, len(dirsToRecurse)-1; i < j; i, j = i+1, j-1 {
			dirsToRecurse[i], dirsToRecurse[j] = dirsToRecurse[j], dirsToRecurse[i]
		}
	}

	// Recursively list subdirectories if -R is set
	for _, dir := range dirsToRecurse {
		listDir(path+"/"+dir, showAll, detailed, reverse, sortByTime, recursive, false)
	}
}

func main() {
	// Parsing command-line flags
	showAll := flag.Bool("a", false, "Show hidden files")
	detailed := flag.Bool("l", false, "Show detailed information")
	reverse := flag.Bool("r", false, "Reverse the sorting order")
	sortByTime := flag.Bool("t", false, "Sort by modification time")
	recursive := flag.Bool("R", false, "List directories recursively")
	flag.Parse()

	// Get the directory to list (default is current directory)
	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Args()[0]
	}

	// Call the function to list the directory (root level)
	listDir(dir, *showAll, *detailed, *reverse, *sortByTime, *recursive, true)
}
