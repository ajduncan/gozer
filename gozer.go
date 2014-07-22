package main

import (
	"bytes"
	"github.com/go-martini/martini"
	"github.com/jessevdk/go-flags"
	"github.com/kr/fs"
	"github.com/pmylund/go-cache"
	"index/suffixarray"
	"io/ioutil"
	"fmt"
	"os"
	"strconv"
	"time"
)

var opts struct {
	Daemonize bool   `short:"d"               description:"Run continuous indexing as a daemon."`
	Help      bool   `short:"h" long:"help"   description:"Show help"`
	Path      string `short:"p" long:"path"   description:"Path to begin search"`
	Search    string `short:"s" long:"search" description:"String to search for."`
}

type Keymaster struct {
	index int
	path string
	context string
}

var gatekeeper = cache.New(30*time.Minute, 180*time.Second)


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
				data, err := ioutil.ReadFile(walker.Path())
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					continue
				}
				index := suffixarray.New(data)
				s := index.Lookup([]byte(stf), 1)
				if len(s) > 0 {
					fmt.Printf("Indexing file: %s\n", walker.Path())
					padding := s[0] + 25
					if padding > cap(data) { padding = cap(data) }

					if km_docs, found := gatekeeper.Get(stf); found {
						di := Keymaster{s[0], walker.Path(), string(data[s[0]:padding])}
						km_docs = append(km_docs.([]Keymaster), di)
						gatekeeper.Set(stf, km_docs, 0)
					} else {
						var document_indexes []Keymaster
						di := Keymaster{s[0], walker.Path(), string(data[s[0]:padding])}
						document_indexes = append(document_indexes, di)
						gatekeeper.Set(stf, document_indexes, 0)
					}

					fmt.Printf("%s(%s): %s\n", walker.Path(), strconv.Itoa(s[0]), data[s[0]:padding])
				}
			}
		}
	} else {
		fmt.Printf("Path does not exist: %s", path)
	}
}

func km_string(km_doc []Keymaster) (km_string string) {
    var result bytes.Buffer

	for _, value := range km_doc {
		result.WriteString(value.path + "(" + strconv.Itoa(value.index) + ") ... " + value.context + " ... \n")
	}

	return result.String()
}

func search(path string, stf string) (result string) {
	if key, found := gatekeeper.Get(stf); found {
		return km_string(key.([]Keymaster))
	} else {
		if path != "" {
			index_search(path, stf)
		} else {
			index_search("./", stf)
		}
		if key, found := gatekeeper.Get(stf); found {
			return km_string(key.([]Keymaster))
		} else {
			return "None found."
		}
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
			return "Main index... use /search/<key>."
		})
		m.Get("/search/:key", func (params martini.Params) string {
			return search(opts.Path, params["key"])
		})
		m.Run()
	} else {
		fmt.Printf("%s\n", search(opts.Path, opts.Search))
	}

}