package main

import (
	"encoding/csv"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
)

func main() {
	// we open the csv file from the disk
	f, err := os.Open("/Users/Sara/GolandProjects/cmpscfa23team2/cuda/data2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// we create a new csv reader specifying
	// the number of columns it has
	salesData := csv.NewReader(f)
	salesData.FieldsPerRecord = 21

	// we read all the records
	records, err := salesData.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	header := records[0]
	// by slicing the records we skip the header
	records = records[1:]

	// we iterate over all the records
	// and keep track of all the gathered values
	// for each column
	columnsValues := map[int]plotter.Values{}
	for i, record := range records {
		// we want one histogram per column,
		// so we will iterate over all the columns we have
		// and gather the date for each in a separate value set
		// in columnsValues
		// we are skipping the ID column and the Date,
		// so we start on index 2
		for c := 2; c < salesData.FieldsPerRecord; c++ {
			if _, found := columnsValues[c]; !found {
				columnsValues[c] = make(plotter.Values, len(records))
			}
			// we parse each close value and add it to our set
			floatVal, err := strconv.ParseFloat(record[c], 64)
			if err != nil {
				log.Fatal(err)
			}
			columnsValues[c][i] = floatVal
		}
	}

	// once we have all the data, we draw each graph
	for c, values := range columnsValues {
		// create a new plot
		p := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histogram of %s", header[c])

		// create a new normalized histogram
		// and add it to the plot
		h, err := plotter.NewHist(values, 16)
		if err != nil {
			log.Fatal(err)
		}
		h.Normalize(1)
		p.Add(h)

		// create a PNG file:
		// Create a new blank image with width 200 and height 100
		width := 400
		height := 400
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// Fill the image with a blue color
		blue := color.RGBA{0, 0, 255, 255}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, blue)
			}
		}

		// Create a new PNG file
		file, err := os.Create("test_hist.png")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// save the plot to a PNG file.
		if err := p.Save(
			10*vg.Centimeter,
			10*vg.Centimeter,
			"test_hist.png",
		); err != nil {
			log.Fatal(err)
		}
	}
}
