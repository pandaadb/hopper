# Hopper

A very simple trace-hop-finder* 



\* a bunny term

# Usage 

## Help 

```
❯ ./hopper 
[-i|--input] is required
usage: print [-h|--help] -i|--input "<value>" [-t|--threshold <integer>]
             [-s|--trace "<value>" [-s|--trace "<value>" ...]]
             [-f|--find-document]

             Prints provided string to stdout

Arguments:

  -h  --help           Print help information
  -i  --input          The input json file
  -t  --threshold      The threshold on which to alert. Default: 1000
  -s  --trace          A list of traces. If not specified, all traces will be
                       considered.
  -f  --find-document  Flag to make the parser search for a document in the
                       line. Helpful for when there is a prefix to the document
                       in line.. Default: false


```

* input 

A valid file path to a file containing json, 1 document per line 

* threshold 

The threshold in millisecond applied to find hops 

* trace 

A list of IDs to look for and group together. Expects a field called 'trace_id'


## What it does 

Hopper will read your input and sort logs based on traces and time. It will then find documents that are logged and their time delta is bigger than the given threshold. 

Hopper's expectations: 

* `trace_id` must be a root-level document field and a string 
* `time` must be a root-level document field and a string 
* `time` is formatted to correspond to `RFC3339` 
* `input` points to a file containing 1 json document per line. 

## What I use it for 

Sometimes calls are slow. And looking at each log line and head-diffing the time is annoying. This somewhat helps that task :) 