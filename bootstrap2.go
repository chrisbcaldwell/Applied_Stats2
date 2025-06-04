package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"bootstrap"

	"github.com/go-gota/gota/dataframe"
)

func main() {
	start := time.Now()
	df, err := readToDataFrame(filepath)
	if err != nil {
		fmt.Println("Terminating process")
		log.Fatal(err)
	}

	cols := []string{"AVG", "BB%", "R"}
	iterations := 5000
	bootstart := time.Now()                        // start the bootstrapping clock
	median := bootstrap.NewQuantileAggregator(0.5) // define the median for bootstrap
	resampler := bootstrap.NewBasicResampler(median, iterations)
	for _, col := range cols {
		// possible insertion of concurrency here
		// would need to create new resmapler here instead of before loop
		s := df.Col(col)
		vals := s.Float()
		resampler.Resample(vals)
		sd := resampler.StdErr()
		fmt.Println("Variable:", col)
		fmt.Println("Median", s.Median())
		fmt.Println("Standard Error:", sd)
		fmt.Println()
		resampler.Reset()
	}

	fmt.Println("Bootstrapping run time:", time.Since(bootstart))
	fmt.Println("Total run time:", time.Since(start))
}

const (
	filepath = "baseball.csv" // can alter this to any path where CSV file is located
)

func readToDataFrame(p string) (dataframe.DataFrame, error) {
	f, err := os.Open(p)
	if err != nil {
		fmt.Println("Unable to read input file " + filepath)
		return dataframe.New(), err
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)
	return df, nil
}
