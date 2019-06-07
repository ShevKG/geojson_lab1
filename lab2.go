package main

import (
	"fmt" 
	"io/ioutil"
	"math"
	"github.com/davvo/mercator"
	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

const width, height = 256, 256
const mercatorMaxValue float64 = 20037508.342789244
const mercatorToCanvasScaleFactorX = float64(width) / (mercatorMaxValue)
const mercatorToCanvasScaleFactorY = float64(height) / (mercatorMaxValue)

func Paint(z, x, y float64) {
	var err error
	var img string
	
	var featureCollectionJSON []byte
	var filePath = "rf.geojson"

	if featureCollectionJSON, err = ioutil.ReadFile(filePath); err != nil {
		fmt.Println(err.Error())
	}

	if img, err = CreatePNG(featureCollectionJSON, z, x, y); err != nil {
		fmt.Println(err.Error())
	}

	println(img)
}

func CreatePNG(featureCollectionJSON []byte, z float64, x float64, y float64) (string, error) {
	var coordinates [][][][][]float64
	var err error

	if coordinates, err = getUserCoordinates(featureCollectionJSON); err != nil {
		return err.Error(), err
	}

	dc := gg.NewContext(width, height)
	scale := 1.0

	dc.InvertY()
	//отрисовка полигонов
	Polygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetRGB(1, 0, 0)
		PaintPolygonCoordinates(dc, polygonCoordinates, scale, dc.Fill, z, x, y)
	})
	//отрисовка границ
	dc.SetLineWidth(2)
	Polygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetRGB(0, 0.5, 0.5)
		PaintPolygonCoordinates(dc, polygonCoordinates, scale, dc.Stroke, z, x, y)
	})
	var out = "output.png"
	dc.SavePNG(out)
	return out, nil
}
func getUserCoordinates(featureCollectionJSON []byte) ([][][][][]float64, error) {
	var featureCollection *geojson.FeatureCollection
	var err error

	if featureCollection, err = geojson.UnmarshalFeatureCollection(featureCollectionJSON); err != nil {
		return nil, err
	}
	var features = featureCollection.Features
	var coordinates [][][][][]float64
	for i := 0; i < len(features); i++ {
		coordinates = append(coordinates, features[i].Geometry.MultiPolygon)
	}
	return coordinates, nil
}
func Polygon(dc *gg.Context, coordinates [][][][][]float64, callback func([][]float64)) {
	for i := 0; i < len(coordinates); i++ {
		for j := 0; j < len(coordinates[i]); j++ {
			callback(coordinates[i][j][0])
		}
	}
}
func PaintPolygonCoordinates(dc *gg.Context, coordinates [][]float64, scale float64, method func(), z float64, xTile float64, yTile float64) {

	scale = scale * math.Pow(2, z)

	dx := float64(dc.Width())*(xTile) - 138.5*scale
	dy := float64(dc.Height())*(math.Pow(2, z)-1-yTile) - 128*scale

	for index := 0; index < len(coordinates)-1; index++ {
		x, y := mercator.LatLonToMeters(coordinates[index][1], convertX(coordinates[index][0]))

		x, y = centerPolygon(x, y)

		x *= mercatorToCanvasScaleFactorX * scale * 0.5
		y *= mercatorToCanvasScaleFactorY * scale * 0.5

		x -= dx
		y -= dy

		dc.LineTo(x, y)
	}
	dc.ClosePath()
	method()
}
func centerPolygon(x float64, y float64) (float64, float64) {
	var west = float64(1635093.15883866)

	if x > 0 {
		x -= west
	} else {
		x += 2*mercatorMaxValue - west
	}

	return x, y
}
func convertX(x float64) float64 {
	if x < 0 {
		x = x - 360
	}
	return x
}
func main() {
	var z, x, y float64
	fmt.Scan(&z, &x, &y)
	Paint(z, x, y)
}
