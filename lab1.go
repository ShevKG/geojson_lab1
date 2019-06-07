package main

import "fmt"
import "os"
import "io"
import geojson "github.com/paulmach/go.geojson" 
import "github.com/fogleman/gg"

const NameFile = "map.geojson"
var ColorFon, ColorLine string
var Width float64
var XYN [255][2]float64
var Long int

func Read()string{
	file, err := os.Open(NameFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	FileByte := make([]byte, 64)
	var File string

	for {
	rsc, err := file.Read(FileByte)
		if err == io.EOF {
			break
		}
		File += string(FileByte[:rsc])
	}

	return File
}

func Parthing(){
	structl, _ := geojson.UnmarshalFeatureCollection([]byte(Read()))
	//Парсинг параметров
	ColorLine = structl.Features[0].Properties["stroke"].(string)
	Width = structl.Features[0].Properties["stroke-width"].(float64)
	ColorFon = structl.Features[0].Properties["fill"].(string)

	Coordinates := structl.Features[0].Geometry.Polygon
	Long = len(Coordinates[0])
	for i:=0; i<Long; i++{
		XYN[i][0] = Coordinates[0][i][0]
		XYN[i][1] = Coordinates[0][i][1]
	}

}

func Paint (){
	const width = 1366
	const height = 1024
	var coefW float64 = 0.3 * width
	var coefH float64 = 0.3 * height
	var cW float64 = -32400
	var cH float64 = -18000
	
	//Подготовка пространства вывода
	Image := gg.NewContext(width, height)
	Image.InvertY()
	Image.SetRGB(1, 1, 1)
	Image.Clear()
	//Отрисовка фигуры
	for i := 0; i < Long-1; i++ {
		if i == 0 {
			Image.MoveTo((XYN[i][0] * coefW + cW), XYN[i][1] * coefH + cH)
		} else {
			Image.LineTo((XYN[i][0] * coefW + cW), XYN[i][1] * coefH + cH)
		}
		fmt.Println((XYN[i][0] * coefW + cW), XYN[i][1] * coefH + cH)
	}
	Image.ClosePath()
	Image.SetHexColor(ColorFon)
	Image.Fill()

	//Отрисовка контура фигуры
	Image.SetHexColor(ColorLine)
	Image.SetLineWidth(Width)
	for i := 0; i < Long-1; i++ {
		if i == 0 {
			Image.MoveTo((XYN[i][0] * coefW + cW), XYN[i][1] * coefH + cH)
		} else {
			Image.LineTo((XYN[i][0] * coefW + cW), XYN[i][1] * coefH + cH)
		}
	}
	Image.ClosePath()
	Image.Stroke()

	//Сохранение в файл
	Image.SavePNG("output.png")
}


func main() {

	Parthing()
	Paint ()
	
}
