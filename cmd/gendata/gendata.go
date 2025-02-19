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

const EarthRadius float64 = 6372.8
const RadiansPerDegree = math.Pi / 180.0

type Cluster struct {
	Mul    float64
	Offset float64
}

func pickRandomCluster(r *rand.Rand, clusters []Cluster) Cluster {
	return clusters[r.UintN(uint(len(clusters)))]
}

func pickRandomCluster2(r *rand.Rand, clusters []Cluster) Cluster {
	c := pickRandomCluster(r, clusters)
	c.Mul *= 2
	c.Offset *= 2
	return c
}

func square(n float64) float64 {
	return n * n
}

func referenceHaversine(x0, y0, x1, y1, radius float64) float64 {
	lat1 := y0
	lat2 := y1
	lon1 := x0
	lon2 := x1

	dLat := (lat2 - lat1) * RadiansPerDegree
	dLon := (lon2 - lon1) * RadiansPerDegree

	lat1 *= RadiansPerDegree
	lat2 *= RadiansPerDegree

	a := square(math.Sin(dLat/2.0)) + math.Cos(lat1)*math.Cos(lat2)*square(math.Sin(dLon/2.0))
	c := 2.0 * math.Asin(math.Sqrt(a))

	return c * radius
}

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

	var clusters []Cluster

	for range 8 {
		mul := rng.UintN(180)
		minOffset := -90
		maxOffset := 90 - int(mul)
		offset := minOffset + rng.IntN(maxOffset-minOffset)

		clusters = append(clusters, Cluster{
			Mul:    float64(mul),
			Offset: float64(offset),
		})
	}

	clusterX0 := pickRandomCluster2(rng, clusters)
	clusterY0 := pickRandomCluster(rng, clusters)
	clusterX1 := pickRandomCluster2(rng, clusters)
	clusterY1 := pickRandomCluster(rng, clusters)
	clusterSwapChance := 1.0 / float64(rows)

	var sum float64

	for i := 0; i < rows; i++ {
		lastRow := i == rows-1

		if n := rng.Float64(); n <= clusterSwapChance {
			clusterX0 = pickRandomCluster2(rng, clusters)
		}
		if n := rng.Float64(); n <= clusterSwapChance {
			clusterY0 = pickRandomCluster(rng, clusters)
		}
		if n := rng.Float64(); n <= clusterSwapChance {
			clusterX1 = pickRandomCluster2(rng, clusters)
		}
		if n := rng.Float64(); n <= clusterSwapChance {
			clusterY1 = pickRandomCluster(rng, clusters)
		}

		x0 := rng.Float64()*clusterX0.Mul + clusterX0.Offset
		y0 := rng.Float64()*clusterY0.Mul + clusterY0.Offset
		x1 := rng.Float64()*clusterX1.Mul + clusterX1.Offset
		y1 := rng.Float64()*clusterY1.Mul + clusterY1.Offset

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

		sum += referenceHaversine(x0, y0, x1, y1, EarthRadius)
	}

	bw.WriteString("]}\n")

	if err = bw.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %d rows.\n", rows)
	fmt.Printf("Took %s.\n", time.Since(start))

	fmt.Printf("Seed: %d\n", seed)

	average := sum / float64(rows)
	fmt.Printf("Haversine average: %f\n", average)

	return nil
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
