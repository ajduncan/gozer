package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/kr/fs"
	"fmt"
	"os"
//	"path/filepath"
)


var opts struct {
	Daemonize bool   `short:"d"               description:"Run continuous indexing as a daemon."`
	Help      bool   `short:"h" long:"help"   description:"Show help"`
	Path      string `short:"p" long:"path"   description:"Path to begin search"`
	Search    string `short:"s" long:"search" description:"String to search for."`
}

func path_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { 
		return true, nil 
	}

	if os.IsNotExist(err) { 
		return false, nil 
	}

	return false, err
}

func kwic_search(path string) {
	r, err := path_exists(path)
	if err != nil {
		fmt.Print("Error 1.\n")
	}

	if r {
		fmt.Printf("Walking path: %s\n", path)
		walker := fs.Walk(path)
		for walker.Step() {
			// if the path we're searching is the path we're on then let's
			// take a step back and ask ourselves what problem are we solving?
			if path == walker.Path() {
				walker.Step()
				continue
			}

			if err := walker.Err(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fi := walker.Stat()
			if fi.IsDir() {
				kwic_search(walker.Path())
			} else {
				fmt.Printf("Indexing file: %s\n", walker.Path())
			}

		}
	} else {
		fmt.Printf("Path does not exist: %s", path)
	}
}

func main() {
	argparser := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)

	_, err := argparser.Parse()
	if err != nil {
		return
	}

	if opts.Help || len(os.Args[1:]) == 0 {
		argparser.WriteHelp(os.Stdout)
		return
	}

	// le debug
	if len(os.Args[1:]) > 0 {
		fmt.Printf("Search: %s\n", opts.Search)
		fmt.Printf("Path: %s\n", opts.Path)
		fmt.Printf("Help: %t\n", opts.Help)
		fmt.Printf("Daemonize: %t\n", opts.Daemonize)
	}

	if opts.Path != "" {
		kwic_search(opts.Path)
	} else {
		kwic_search("./")
	}

}