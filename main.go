package main

import (
	"context"
	"fmt"
	"github.com/akamensky/argparse"
	"hopper/app"
	"os"
)

type input struct {
}

func main() {

	parser := argparse.NewParser("print", "Prints provided string to stdout")

	input := parser.String("i", "input", &argparse.Options{Required: true, Help: "The input json file"})
	threshold := parser.Int("t", "threshold", &argparse.Options{Required: false, Help: "The threshold on which to alert", Default: 1000})
	traces := parser.StringList("s", "trace", &argparse.Options{Required: false, Help: "A list of traces. If not specified, all traces will be considered."})
	findDocument := parser.Flag("f", "find-document", &argparse.Options{Required: false, Help: "Flag to make the parser search for a document in the line. Helpful for when there is a prefix to the document in line.", Default: false})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	fileParser := app.Parser{
		Traces:       *traces,
		Input:        *input,
		FindDocument: *findDocument,
	}

	res, err := fileParser.Parse(context.Background())
	if err != nil {
		fmt.Printf("Error parsing: %w", err)
	}

	res.PrintHops(int64(*threshold))
}
