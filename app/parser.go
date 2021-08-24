package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Parser struct {
	Traces []string
	Input  string
}

type ParseResult struct {
	Traces map[string]timedDocumentSlice
}

type TimedDocument struct {
	Time  time.Time
	Trace string
	doc   map[string]interface{}
}

type timedDocumentSlice []TimedDocument

func (td timedDocumentSlice) Len() int {
	return len(td)
}

func (td timedDocumentSlice) Less(i, j int) bool {
	return td[i].Time.Before(td[j].Time)
}

func (td timedDocumentSlice) Swap(i, j int) {
	td[i], td[j] = td[j], td[i]
}

func (p Parser) Parse(context context.Context) (ParseResult, error) {
	file, err := os.Open(p.Input)
	if err != nil {
		return ParseResult{}, fmt.Errorf("file open error: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	traces := make(map[string]timedDocumentSlice)

	traceSet := map[string]struct{}{}
	for _, traceId := range p.Traces {
		traceSet[traceId] = struct{}{}
	}

	for scanner.Scan() {
		var doc map[string]interface{}

		if err := json.Unmarshal(scanner.Bytes(), &doc); err != nil {
			fmt.Printf("error parsing json: %w\n", err.Error())
			continue
		}

		if _, ok := doc["time"]; !ok {
			fmt.Printf("no 'time' field\n")
			continue
		}

		if _, ok := doc["trace_id"]; !ok {
			continue
		}

		trace := fmt.Sprintf("%s", doc["trace_id"])

		if _, ok := traceSet[trace]; !ok && len(p.Traces) > 0 {
			continue
		}

		parsedTime, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", doc["time"]))
		if err != nil {
			fmt.Printf("'time' parse failure: %s\n", err.Error())
			continue
		}

		traces[trace] = append(traces[trace], TimedDocument{
			Time:  parsedTime,
			Trace: trace,
			doc:   doc,
		})
	}

	if err := scanner.Err(); err != nil {
		return ParseResult{}, fmt.Errorf("scanner error: %w\n", err)
	}

	for _, timedDocuments := range traces {
		sort.Sort(timedDocuments)
	}

	return ParseResult{
		Traces: traces,
	}, nil
}

func (p ParseResult) PrintHops(threshold int64) {
	for traceId, docs := range p.Traces {

		for i, doc := range docs {
			if i == len(docs)-1 {
				break
			}

			first := doc.Time
			second := docs[i+1].Time

			duration := second.Sub(first)

			if duration.Milliseconds() > threshold {
				fmt.Printf("\n\n========================\n\n")
				fmt.Printf("Hop found:\n")
				fmt.Printf("Trace ID : %s\n", traceId)
				fmt.Printf("First document at: %s\n", first)
				fmt.Printf("Second document at: %s\n", second)
				fmt.Printf("Duration:%dms\n", duration.Milliseconds())

				fmt.Printf("\n\nFirst document:\n\n")
				printDoc(doc.doc)
				fmt.Printf("\n\nSecond document:\n\n")
				printDoc(docs[i+1].doc)
			}
		}

	}

}

func printDoc(in interface{}) {
	firstDocString, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling first document: %s\n", err)
		fmt.Printf("%s\n", in)
	}

	fmt.Printf("%s\n", firstDocString)
}
