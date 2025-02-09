package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

const AverageRowSize = len(
	`    {"x0":0.00000000000000000,"y0":0.00000000000000000,"x1":0.00000000000000000,"y1":0.00000000000000000},`,
)

func run() error {
	defaultSeed := uint64(time.Now().UnixMilli())

	var filename string
	var rows int
	var seed uint64
	flag.IntVar(&rows, "rows", 1_000, "number of rows to generate")
	flag.StringVar(&filename, "out", "data.json", "output filename")
	flag.Uint64Var(&seed, "seed", defaultSeed, "random number generator seed")
	flag.Parse()

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	bw := bufio.NewWriterSize(file, min(rows*AverageRowSize, math.MaxInt32))

	bw.WriteString(`{"pairs":[` + "\n")

	pcg := rand.NewPCG(0, seed)
	rng := rand.New(pcg)

	for i := 0; i < rows; i++ {
		lastRow := i == rows-1

		x0 := rng.Float64()
		y0 := rng.Float64()
		x1 := rng.Float64()
		y1 := rng.Float64()

		bw.WriteString(`    {"x0":`)
		bw.WriteString(strconv.FormatFloat(x0, 'f', 17, 64))
		bw.WriteString(`,"y0":`)
		bw.WriteString(strconv.FormatFloat(y0, 'f', 17, 64))
		bw.WriteString(`,"x1":`)
		bw.WriteString(strconv.FormatFloat(x1, 'f', 17, 64))
		bw.WriteString(`,"y1":`)
		bw.WriteString(strconv.FormatFloat(y1, 'f', 17, 64))

		if lastRow {
			bw.WriteString("}\n")
		} else {
			bw.WriteString("},\n")
		}
	}

	bw.WriteString("]}\n")

	if err = bw.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %d rows.\n", rows)
	fmt.Printf("Took %s.\n", time.Since(start))

	return nil
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
