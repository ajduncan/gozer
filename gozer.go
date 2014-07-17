package main

import (
	"github.com/jessevdk/go-flags"
	"io/ioutil"
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

func is_directory(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false, err
	}

	switch mode := fi.Mode(); {
		case mode.IsDir():
			return true, nil
	}

	return false, nil
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
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			r, err := is_directory(f.Name())
			if err != nil {
				fmt.Print("Error 2\n")
			}
			if r {
				fmt.Printf("Directory: %s\n", f.Name())
				kwic_search(f.Name())
			} else {
				fmt.Printf("Indexing file: %s\n", f.Name())
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