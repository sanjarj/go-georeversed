package search

import "github.com/sanjarj/go-georeversed/geo"

type Index interface {
	Find(lat, lon float64) (string, bool)
}

type simpleIndex struct {
	countries []*geo.Country
}

func (index *simpleIndex) Find(lat, lon float64) (string, bool) {
	for _, country := range index.countries {
		if country.Contains(lat, lon) {
			return country.Code, true
		}
	}
	return "", false
}

func NewIndex(countries []*geo.Country) Index {
	return &simpleIndex{countries}
}
