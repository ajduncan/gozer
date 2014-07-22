package gozer

import (
	"bytes"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/kr/fs"
	"github.com/pmylund/go-cache"
	"index/suffixarray"
	"io/ioutil"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

const CACHE_DB string = "./gatekeeper.db"

var Opts struct {
	Daemonize bool   `short:"d"               description:"Run continuous indexing as a daemon."`
	Help      bool   `short:"h" long:"help"   description:"Show help"`
	Path      string `short:"p" long:"path"   description:"Path to begin search"`
	Search    string `short:"s" long:"search" description:"String to search for."`
}

var gatekeeper = cache.New(30*time.Minute, 180*time.Second)

type Keymaster struct {
	Index int
	Path string
	Context string
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

func save_result(stf string, index int, path string, context string) {
	if km_docs, found := gatekeeper.Get(stf); found {
		di := Keymaster{index, path, context}
		km_docs = append(km_docs.([]Keymaster), di)
		gatekeeper.Set(stf, km_docs, 0)
	} else {
		var document_indexes []Keymaster
		di := Keymaster{index, path, context}
		document_indexes = append(document_indexes, di)
		gatekeeper.Set(stf, document_indexes, 0)
	}

	fmt.Printf("%s(%s): %s\n", path, strconv.Itoa(index), context)
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
					padding := s[0] + 25
					if padding > cap(data) { padding = cap(data) }
					save_result(stf, s[0], walker.Path(), string(data[s[0]:padding]))
				}
			}
		}
	} else {
		fmt.Printf("Path does not exist: %s", path)
	}
}

func KM_string(km_doc []Keymaster) (km_string string) {
    var result bytes.Buffer

	for _, value := range km_doc {
		result.WriteString(value.Path + "(" + strconv.Itoa(value.Index) + ") ... " + value.Context + " ... ")
	}

	return result.String()
}

func Search(path string, stf string) (result []Keymaster) {
	if key, found := gatekeeper.Get(stf); found {
		return key.([]Keymaster)
	} else {
		if path != "" {
			index_search(path, stf)
		} else {
			index_search("./", stf)
		}
		if key, found := gatekeeper.Get(stf); found {
			return key.([]Keymaster)
		} else {
			return nil
		}
	}
}

func Daemonize() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})
	m.Get("/search/", func (params martini.Params, request *http.Request, r render.Render) {
		search_query := request.URL.Query().Get("search")
		fmt.Printf("Searching for: %s", search_query)
		results := Search(Opts.Path, search_query)
		fmt.Printf("Got results: %s", results)
		content := map[string]interface{}{"results": results}
		r.HTML(200, "results", content)
	})
	m.Run()
}

func LoadCache() {
	// load cache from disk if available.
	if _, err := os.Stat(CACHE_DB); err == nil {
		fmt.Printf("Loading gatekeeper db from disk.")
		if err := gatekeeper.LoadFile(CACHE_DB); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func SaveCache() {
	// save cache periodically.
	gk_ticker      := time.NewTicker(60 * time.Second)
	gk_ticker_quit := make(chan struct{})
	go func() {
		for {
			select {
				case <- gk_ticker.C:
					fmt.Printf("Saving DB to disk...")
					if err := gatekeeper.SaveFile(CACHE_DB); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				case <- gk_ticker_quit:
					gk_ticker.Stop()
					return
			}
		}
	}()
	// close(gk_ticker_quit) to end the timer.
}
