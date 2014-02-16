package geo

import (
	"github.com/sanjarj/go-georeversed/geojson"
	"errors"
	"math"
)

func FromGeoJson(json *geojson.FeatureCollection) ([]*Country, error) {
	result := make([]*Country, len(json.Features))
	for i, feature := range json.Features {
		code := feature.Properties["ISO2"].(string)
		name := feature.Properties["NAME"].(string)

		switch feature.Geometry.Type{
		case "MultiPolygon":
			result[i] = &Country{code, name, *newMultiPolygon(feature.Geometry.Coordinates)}
		case "Polygon":
			result[i] = &Country{code, name, multiPolygon{*newPolygon(feature.Geometry.Coordinates)}}
		default:
			return nil, errors.New("unexpected type "+feature.Geometry.Type)
		}
	}
	return result, nil
}

func newMultiPolygon(mp []interface{}) *multiPolygon {
	result := make(multiPolygon, len(mp))
	for i := range mp {
		result[i] = *newPolygon(mp[i].([]interface{}))
	}
	return &result
}

func newPolygon(p []interface{}) *polygon {
	result := make(polygon, len(p))
	for i := range p {
		result[i] = *newLinearLing(p[i].([]interface{}))
	}
	return &result
}

func newLinearLing(lr []interface{}) *linearRing {
	result := make(linearRing, len(lr))
	for i := range lr {
		p := lr[i].([]interface{})
		result[i] = point{p[0].(float64), p[1].(float64)}
	}
	return &result
}

type Country struct {
	Code     string
	Name     string
	geometry multiPolygon
}

func (country *Country) Contains(lat, lon float64) bool {
	return country.geometry.contains(&point{lon, lat})
}

type multiPolygon []polygon

func (this *multiPolygon) contains(p *point) bool {
	for _, polygon := range *this {
		if polygon.contains(p) {
			return true
		}
	}
	return false
}

type polygon []linearRing

func (this *polygon) contains(p *point) bool {
	rings := *this
	if rings[0].contains(p) {
		// check holes
		for i := 1; i < len(rings); i++ {
			if rings[i].contains(p) {
				return false
			}
		}
		return true
	}
	return false
}

type linearRing []point

func (this *linearRing) contains(p *point) bool {
	points := *this
	if len(points) < 4 {
		return false
	}
	extreme := &point{1000, p.y}
	count := 0
	for i := 1; i < len(points); i++ {
		prev, cur := &points[i - 1], &points[i]
		if doIntersect(prev, cur, p, extreme) {
			if orientation(prev, p, cur) == 0 {
				return p.onSegment(prev, cur)
			}
			count++
		}
	}
	return count%2 == 1
}

type point struct {
	x, y float64
}

func (p *point) onSegment(q, r *point) bool {
	return p.x <= math.Max(q.x, r.x) && p.x >= math.Min(q.x, r.x) && p.y <= math.Max(q.y, r.y) && p.y >= math.Min(q.y, r.y)
}

func orientation(p, q, r *point) byte {
	v := (q.y - p.y) * (r.x - q.x) - (q.x - p.x) * (r.y - q.y)
	if v > 0 {
		return 1
	} else if v < 0 {
		return 2
	}
	return 0
}

func doIntersect(p1, q1, p2, q2 *point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}

	if o1 == 0 && p2.onSegment(p1, q1) {
		return true
	}
	if o2 == 0 && q2.onSegment(p1, q1) {
		return true
	}
	if o3 == 0 && p1.onSegment(p2, q2) {
		return true
	}
	if o4 == 0 && q1.onSegment(p2, q2) {
		return true
	}
	return false
}
