package main

import (
	"github.com/sanjarj/go-georeversed/geojson"
	"github.com/sanjarj/go-georeversed/geo"
	"github.com/sanjarj/go-georeversed/search"
	"github.com/sanjarj/go-georeversed/web"
	"flag"
	"fmt"
	"os"
)

func main() {
	port := flag.Int("p", 8080, "port to listen to")
	flag.Parse()

	file := flag.Arg(0)
	if file == "" {
		printUsage()
		return
	}

	countries, err := readBorders(file)
	if (err != nil) {
		panic(err)
	}

	index := search.NewIndex(countries)
	if err := web.ListenAndServe(fmt.Sprintf(":%d", *port), index); err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go-georeversed [-p=port] <world-borders-geojson-file>")
}

func readBorders(filename string) ([]*geo.Country, error) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		if features, err := geojson.Decode(file); err == nil {
			return geo.FromGeoJson(features)
		}else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
