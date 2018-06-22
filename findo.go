package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mattn/go-isatty"
)

var flagfileOnly = flag.Bool("f", false, "Select fileonly(Remove directories")
var quotation = flag.Bool("q", false, "Quotation filenames")
var nameOnly = flag.Bool("1", false, "Show nameonly(No Size,timestamp)")
var flagList = flag.Bool("l", false, "Show Size and timestamp")
var startDir = flag.String("d", ".", "Set start Directory")

func main1(args []string) error {
	var pattern string
	if len(args) >= 1 {
		pattern = strings.ToUpper(args[0])
	}

	rich := isatty.IsTerminal(os.Stdout.Fd())
	if *flagList {
		rich = true
	}
	if *nameOnly {
		rich = false
	}

	filepath.Walk(*startDir, func(path_ string, info_ os.FileInfo, err_ error) error {
		name := filepath.Base(path_)
		if name == "." || name == ".." {
			return nil
		}
		if *flagfileOnly && info_.IsDir() {
			return nil
		}
		var matched bool
		if pattern == "" {
			matched = true
		} else {
			var err error
			matched, err = filepath.Match(pattern, strings.ToUpper(name))
			if err != nil {
				matched = false
			}
		}
		if matched {
			if *quotation {
				fmt.Printf("\"%s\"\n", path_)
			} else {
				fmt.Println(path_)
			}
			if rich {
				fmt.Printf("%12s %s\n", humanize.Comma(info_.Size()), info_.ModTime().String())
			}
		}
		return nil
	})
	return nil
}

func main() {
	flag.Parse()
	if err := main1(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
