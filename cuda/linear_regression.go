package main

import (
	"github.com/gocarina/gocsv"
	"os"
)

type KCHouseData struct {
	ID           string  `csv:"id"`
	Date         string  `csv:"date"`
	Price        int     `csv:"price"`
	Bedrooms     int     `csv:"bedrooms"`
	Bathrooms    float32 `csv:"bathrooms"`
	SqftLiving   int     `csv:"sqft_living"`
	SqftLot      int     `csv:"sqft_lot"`
	Floors       string  `csv:"floors"`
	Waterfront   int     `csv:"waterfront"`
	View         int     `csv:"view"`
	Condition    int     `csv:"condition"`
	Grade        int     `csv:"grade"`
	SqftAbove    int     `csv:"sqft_above"`
	SqftBasement int     `csv:"sqft_basement"`
	YrBuilt      int     `csv:"yr_built"`
	YrRenovated  int     `csv:"yr_renovated"`
	Zipcode      string  `csv:"zipcode"`
	Lat          float32 `csv:"lat"`
	Long         float32 `csv:"long"`
	SqftLiving15 int     `csv:"sqft_living15"`
	SqftLot15    int     `csv:"sqft_lot15"`
}

func main() {
	file, err := os.Create("data2.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	kc_house_data := []*KCHouseData{
		{"7129300520", "20141013T000000", 221900, 3, 1, 1180, 5650, "1", 0, 0, 3, 7, 1180, 0, 1955, 0, "98178", 47.5112, -122.257, 1340, 5650},
		{"6414100192", "20141209T000000", 538000, 3, 2.25, 2570, 7242, "2", 0, 0, 3, 7, 2170, 400, 1951, 1991, "98125", 47.721, -122.319, 1690, 7639},
		{"5631500400", "20150225T000000", 180000, 2, 1, 770, 10000, "1", 0, 0, 3, 6, 770, 0, 1933, 0, "98028", 47.7379, -122.233, 2720, 8062},
	}

	if err := gocsv.MarshalFile(&kc_house_data, file); err != nil {
		panic(err)
	}

}
