package web

import (
	"github.com/sanjarj/go-georeversed/search"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func ListenAndServe(address string, index search.Index) error {
	// country by location
	http.HandleFunc("/country/", func(w http.ResponseWriter, req *http.Request) {
			req.ParseForm()
			l, err := getLocation(req)
			if (err != nil) {
				writeClientError(w, err)
				return
			}
			if code, ok := index.Find(l.lat, l.lon); ok {
				writeJson(w, response{"status":"success", "country":code})
				return
			}
			writeJson(w, response{"status":"success"})
		})
	// TODO countries by area - func(area) -> []countries
	// TODO countries in area - func(area, []countries) -> []countriesInArea
	return http.ListenAndServe(address, nil)
}

func writeClientError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	writeJson(w, response{"status":"error", "message":err.Error()})
}

func writeJson(w http.ResponseWriter, resp response) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type response map[string]interface{}

type location struct {
	lat, lon float64
}

type area struct {
	lat, lon, radius float64
}

func getLocation(req *http.Request) (*location, error) {
	lat, err := getFloatParam(req, "lat")
	if err != nil {
		return nil, err
	}
	lon, err := getFloatParam(req, "lon")
	if err != nil {
		return nil, err
	}
	return &location{lat, lon}, nil
}

func getFloatParam(req *http.Request, name string) (float64, error) {
	if value, prs := req.Form[name]; prs {
		if len(value) == 1 {
			return strconv.ParseFloat(value[0], 64)
		}
		return 0, errors.New("exactly one "+name+" expected")
	}
	return 0, errors.New(name+" expected")
}
