package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
)

// PROJECTID : BigQuery にアクアスするアカウントのプロジェクトID
const PROJECTID = "sample-1385"

// TODO: 全体的にリファクタリングが必要
func main() {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Number of births (USA)"
	p.X.Label.Text = "year"
	p.Y.Label.Text = "number"

	var resultArray [100][2]int
	resultSize, err := query(&resultArray)
	if err != nil {
		panic(err)
	}
	fmt.Print("resultSize=")
	fmt.Print(resultSize)
	fmt.Print("\n")

	if err := plotutil.AddLinePoints(p, "sample", plotData(resultArray, 20)); err != nil {
		panic(err)
	}

	if err := p.Save(5*vg.Inch, 5*vg.Inch, "graph.png"); err != nil {
		panic(err)
	}

}

func query(resultArray *[100][2]int) (int, error) {
	resultSize := 0
	jsonData, err := ioutil.ReadFile("client.json")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	conf, err := google.JWTConfigFromJSON(jsonData, bigquery.BigqueryScope)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	client := conf.Client(oauth2.NoContext)

	// query の例
	conn, err := bigquery.New(client)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	query := "SELECT year, COUNT(1) as count FROM publicdata:samples.natality GROUP BY year ORDER BY year"
	result, err := conn.Jobs.Query(PROJECTID, &bigquery.QueryRequest{
		Query: query,
	}).Do()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	for x, row := range result.Rows {
		for y, cell := range row.F {
			value, err := strconv.Atoi(cell.V.(string))
			if err != nil {
				log.Fatal(err)
				return resultSize, err
			}
			(*resultArray)[x][y] = value
			resultSize++
		}
	}
	return resultSize, nil
}

func plotData(data [100][2]int, dataSize int) plotter.XYs {

	pts := make(plotter.XYs, dataSize)
	for i := 0; i < dataSize; i++ {
		pts[i].X = float64(data[i][0])
		pts[i].Y = float64(data[i][1])
		//fmt.Print(i)
		//fmt.Print(": ")
		//fmt.Print(pts[i].X)
		//fmt.Print(", ")
		//fmt.Print(pts[i].Y)
		//fmt.Print("\n")
	}

	return pts
}
