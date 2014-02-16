go-georeversed
===

**go-georeversed** is a simple reverse geocode service written in Go.

This is a toy project, which i wrote just to play with Go language. Use it at your own risk.


World borders data
---

Sources doesn't contain actual geographical data. You should provide world borders data in GeoJson format when running the service:

1. Download and unzip world borders shapefile TM_WORLD_BORDERS-0.3.zip from [thematicmapping.org](http://thematicmapping.org/downloads/world_borders.php)
2. Convert unzipped shapefile to GeoJson format using `ogr2ogr` tool which is part of [Geospatial Data Abstraction Library](http://www.gdal.org/). As an example, this command will generate `world.borders.geo.json` GeoJson file from shapefile: 
	
		$ ogr2ogr -f GeoJSON world.borders.geo.json TM_WORLD_BORDERS-0.3.shp


Endpoints
---

There is only one endpoint - `/country/?lat=<latitude>&lon=<longitude>` which looks up country by location.

	$ curl -s "http://localhost:8080/country/?lat=52.0&lon=13.0"
	{
	  "country": "DE",
	  "status": "success"
	}
