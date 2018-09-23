package boat

import (
	"fmt"
	"math/rand"
	"time"
)

func CompressionTest() {
	st := time.Now()
	t := NewTile()
	newData := make([][]byte, 0)
	rand.Seed(time.Now().Unix())
	for x := 0; x < 10; x++ {
		row := make([]byte, 0)
		for y := 0; y < 10; y++ {
			n := rand.Int() % 256
			row = append(row, byte(n))
		}
		newData = append(newData, row)
	}
	t.Data = newData
	fmt.Println("New Tile Created")
	t.PrintData()

	resp := true
	for resp {
		fmt.Println("Expanding from: ", len(t.Data))
		resp = t.Expand()
		t.Pickle()
	}
	t.SaveImage()
	fmt.Println("Total Time: ", time.Since(st))
}

func SquareTest() bool {
	//sqare polygon
	var poly, poly2 Poly
	var polySet PolySet

	poly.AddCorner(0, 0)
	poly.AddCorner(0, 100)
	poly.AddCorner(50, 50)
	poly.AddCorner(100, 100)
	poly.AddCorner(100, 0)
	polySet.AddPoly(poly)

	poly2.AddCorner(20, 20)
	poly2.AddCorner(20, 40)
	poly2.AddCorner(40, 40)
	poly2.AddCorner(40, 20)
	polySet.AddPoly(poly2)

	polySet.Print()

	var start, end Point
	x := 60.0
	start.X = 10
	start.Y = 30
	end.X = 70
	end.Y = x

	r, err := ShortestPath(start, end, polySet)
	r.Print()
	Draw(&polySet, &r, start, end)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if r.Count != 0 {
		fmt.Println("Route Count should be Zero")
		return false
	}
	return true
}
