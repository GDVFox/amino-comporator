package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gonum.org/v1/gonum/mat"
)

var (
	dataFile string
)

func init() {
	flag.StringVar(&dataFile, "data", "", "Name of input file")
}

func main() {
	flag.Parse()

	f, err := os.Open(dataFile)
	if err != nil {
		log.Fatalf("can not open file: %s", err)
	}
	defer f.Close()

	seqs := make([]*AminoSequence, 0, 70)
	p := NewFastaParser(f)
	for {
		seq, err := p.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("processing error: %s", err)
		}
		seqs = append(seqs, seq)
	}

	comparisonMatrix := BuildComparisonMatrixConcurrent(seqs, conf.Threads)
	fa := mat.Formatted(comparisonMatrix, mat.Squeeze())
	fmt.Println(fa)
}
