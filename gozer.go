package main

import (
	"fmt"
	"os"

	"github.com/ajduncan/gozer/lib"
	"github.com/jessevdk/go-flags"
)


func main() {
	argparser := flags.NewParser(&gozer.Opts, flags.PrintErrors|flags.PassDoubleDash)

	_, err := argparser.Parse()
	if err != nil {
		return
	}

	if gozer.Opts.Help || len(os.Args[1:]) == 0 {
		argparser.WriteHelp(os.Stdout)
		return
	}

	gozer.Init()
	gozer.LoadCache()
	gozer.SaveCache()

	if gozer.Opts.Daemonize {
		gozer.Daemonize()
	} else {
		fmt.Printf("%s\n", gozer.KM_string(gozer.Search(gozer.Opts.Path, gozer.Opts.Search)))
	}

}