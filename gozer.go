package main

import (
	"github.com/go-martini/martini"
	"github.com/jessevdk/go-flags"
	"github.com/kr/fs"
	"index/suffixarray"
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

func index_search(path string, stf string) {
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
				index_search(walker.Path(), stf)
			} else {
				fmt.Printf("Indexing file: %s\n", walker.Path())
				data, err := ioutil.ReadFile(walker.Path())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				index := suffixarray.New(data)
				s := index.Lookup([]byte(stf), 1)
				if len(s) > 0 {
					fmt.Printf("Found match at: %s", s[0])
				}
			}
		}
	} else {
		fmt.Printf("Path does not exist: %s", path)
	}
}

func search(path string, stf string) {
	if path != "" {
		index_search(path, stf)
	} else {
		index_search("./", stf)
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

	// le martini stub
	if opts.Daemonize {
		m := martini.Classic()
		m.Get("/", func() string {
			search(opts.Path, opts.Search)
			return "Search finished."
		})
		m.Run()
	} else {
		search(opts.Path, opts.Search)
	}

}