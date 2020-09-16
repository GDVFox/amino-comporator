package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pelletier/go-toml"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

const (
	// OutputModeMatrix as math matrix
	OutputModeMatrix = "matrix"
	// OutputModePython as python two dim array
	OutputModePython = "python"
	// OutputModeMatlab as matlab matrix
	OutputModeMatlab = "matlab"
	// OutputModeHuman as human readable text
	OutputModeHuman = "human"
)

var (
	configFile string
	dataFile   string
	outputMode string
)

func init() {
	flag.StringVar(&configFile, "config", "", "Name of config file")
	flag.StringVar(&dataFile, "data", "", "Name of input file")
	flag.StringVar(&outputMode, "output", OutputModeMatrix, "Possible values matrix|python|matlab|human")
}

func printHuman(seqs []*AminoSequence, matrix *mat.Dense) {
	for k, seq := range seqs {
		fmt.Println(">" + seq.ID + "|" + seq.Description)

		_, cols := matrix.Dims()
		rowValues := make([]float64, cols)
		for i := 0; i < cols; i++ {
			// because upper angle matrix
			if k > i {
				rowValues[i] = matrix.At(i, k)
			} else if k < i {
				rowValues[i] = matrix.At(k, i)
			} else {
				// completely identical to itself
				rowValues[i] = 1.
			}
		}

		indexes := make([]int, len(rowValues))
		floats.Argsort(rowValues, indexes)

		for i := len(indexes) - 2; i >= 0; i-- {
			fmt.Printf("[%d]%8.6f|%s|%s\n", len(indexes)-i-1, rowValues[i], seqs[indexes[i]].ID, seqs[indexes[i]].Description)
		}
	}
}

func main() {
	flag.Parse()

	tree, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatalf("can not read config: %s", err)
	}
	if err := tree.Unmarshal(conf); err != nil {
		log.Fatalf("can not unmarshal config: %s", err)
	}

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

	switch outputMode {
	case OutputModeMatrix:
		fmt.Println(mat.Formatted(comparisonMatrix, mat.Squeeze()))
	case OutputModePython:
		fmt.Println(mat.Formatted(comparisonMatrix, mat.FormatPython()))
	case OutputModeMatlab:
		fmt.Println(mat.Formatted(comparisonMatrix, mat.FormatMATLAB()))
	case OutputModeHuman:
		printHuman(seqs, comparisonMatrix)
	default:
		fmt.Println("Unknown output mode. Using default matrix mode.")
		fmt.Println(mat.Formatted(comparisonMatrix, mat.Squeeze()))
	}

}
