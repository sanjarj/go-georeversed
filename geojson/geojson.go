package geojson

import (
	"io"
	"encoding/json"
)

type base struct {
	Type string `json:"type"`
}

type FeatureCollection struct {
	base
	Features []Feature `json:"features"`
}

type Feature struct {
	base
	Properties map[string]interface{} `json:"properties"`
	Geometry   Geometry               `json:"geometry"`
}

type Geometry struct {
	base
	Coordinates []interface{} `json:"coordinates"`
}

func Decode(r io.Reader) (*FeatureCollection, error) {
	result := new(FeatureCollection)
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
