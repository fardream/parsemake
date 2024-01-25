package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/fardream/parsemake/parser"
)

func main() {
	var target string
	flag.StringVar(&target, "target", "", "target to see")
	flag.Parse()

	files := flag.Args()

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			panic(err)
		}
		p, err := parser.Parse(bytes.NewReader(b))
		if err != nil {
			panic(err)
		}

		for _, t := range p.Targets {
			if t.Name == target {
				fmt.Printf("%v\n", t)
			}
		}
	}
}
