package lq

import (
	// "flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

type counts struct{ dirs, files int }

func walk(dir, prefix string, depth, maxDepth int, all bool, c *counts) {
	if depth >= maxDepth {
		return
	}
	entries, err := os.ReadDir(dir) // already sorted by filename
	if err != nil {
		fmt.Fprintf(os.Stderr, "lq: %v\n", err)
		return
	}

	visible := entries[:0:0]
	for _, e := range entries {
		if !all && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		visible = append(visible, e)
	}

	for i, e := range visible {
		conn, next := "├── ", "│   "
		if i == len(visible)-1 {
			conn, next = "└── ", "    "
		}

		path := filepath.Join(dir, e.Name())
		name := e.Name()
		if e.Type()&os.ModeSymlink != 0 {
			if target, err := os.Readlink(path); err == nil {
				name += " -> " + target
			}
		}
		fmt.Println(prefix + conn + name)

		if e.IsDir() { // symlinks report false, so they are never followed
			c.dirs++
			walk(path, prefix+next, depth+1, maxDepth, all, c)
		} else {
			c.files++
		}
	}
}

func Main(args []string) {
	var _depth int
	var depth *int
	var _all bool
	var all *bool 
	var err error 
	
	// depth := flag.Int("L", 5, "maximum depth to descend")
	_depth , err = strconv.Atoi(args[0]) 
	if err != nil {
		if len(args) == 0 {
			_depth = 0 
		} else {
			
		
		panic(err) 
		}
	}
	depth = &_depth
	if len(args) < 2 and args[1] == "-a" {
		_all = false 
	} else {
		_all = true 
	}
	all = &_all
	// all := flag.Bool("a", false, "include hidden (dot) files and directories")
	// flag.Usage = func() {
		// fmt.Fprintln(os.Stderr, "usage: lq [-L depth] [-a] [dir]")
		// flag.PrintDefaults()
	// }
	// flag.Parse()

	root := "."
	// if flag.NArg() > 0 {
		// root = flag.Arg(0)
	// }
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "lq: %s: not a readable directory\n", root)
		os.Exit(1)
	}

	fmt.Println(root)
	var c counts
	walk(root, "", 0, *depth, *all, &c)
	fmt.Printf("\n%d directories, %d files\n", c.dirs, c.files)
}












































