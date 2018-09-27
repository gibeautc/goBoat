package vehical

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
	"time"
	"fmt"
)

func TestGetCords2(t *testing.T) {
	//finding cords for genisis block
	lat,lon:=GetCords(44.616429,-123.072626,2048,180)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//that will be the southCenter cords
	lat,lon=GetCords(44.6159617,-123.072626,2048,270)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//this will be the SW Corner
	//Lat:  44.61596169811975
	//Lon:  -123.07328247537501
	lat,lon=GetCords(44.61596169811975,-123.07328247537501,4096,90)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//this will be SE Corner
	//Lat:  44.61596169059874
	//Lon:  -123.07196952462509
	lat,lon=GetCords(44.61596169059874,-123.07196952462509,4096,0)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//this will be NE Corner
	//Lat:  44.61689628886895
	//Lon:  -123.07196952462509
	lat,lon=GetCords(44.61689628886895,-123.07196952462509,4096,270)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//this will be NW Corner
	//Lat:  44.616896281347714
	//Lon:  -123.07328249650676
	lat,lon=GetCords(44.616896281347714,-123.07328249650676,4096,180)
	fmt.Println("Lat: ",lat)
	fmt.Println("Lon: ",lon)
	//this should be SW corner again
	//this is what we will check
	assert.Equal(t,true,math.Abs(lat-44.61596169811975)<44.61596169811975*.0001)
	assert.Equal(t,true,math.Abs(lon+123.07328247537501)<123.07328247537501*.0001)
}

func TestGetCords(t *testing.T) {
	lat,lon:=GetCords(44.0,-123.0,538976,-35)
	fmt.Println("Lat")
	fmt.Println("Expected: ",44.1)
	fmt.Println("Got: ",lat)
	fmt.Println("")
	fmt.Println("Lon")
	fmt.Println("Expected: ",-123.1)
	fmt.Println("Got: ",lon)


	assert.Equal(t,true,math.Abs(lat-44.1)<44.1*.0001)
	assert.Equal(t,true,math.Abs(lon+123.1)<123.1*.0001)
}

func TestDistanceBetween(t *testing.T) {
	dist,direction:=DistanceBetween(44.0,-123.0,44.1,-123.1)
	//13.69 Km --> 538976.38
	//324 Deg
	delta:=math.Abs(float64(dist-538976))
	assert.Equal(t,true,delta<float64(dist)*.01)
	assert.Equal(t,324,direction)
}


func TestTile_Expand(t *testing.T) {

	tile := NewTile()
	fmt.Println("New Tile Created")
	d:=make([]byte,0)
	d=append(d,128)
	data:=make([][]byte,0)
	data=append(data,d)
	tile.Data=data
	st := time.Now()
	resp := true
	tile.Id=0
	for resp {
		fmt.Println("Expanding from: ", len(tile.Data))
		resp = tile.Expand()
		fmt.Println("Done Expanding")
		tile.Pickle()
	}
	tile.SaveImage()
	fmt.Println("Total Time To Expand: ", time.Since(st))
}

func TestTile_Compress(t *testing.T) {

	tile := NewTile()
	fmt.Println("New Tile Created")
	tile.Id=0
	st := time.Now()
	resp := true
	for resp {
		fmt.Println("Compressing from: ", len(tile.Data))
		resp = tile.Compress()
		tile.Pickle()
	}
	tile.SaveImage()
	fmt.Println("Total Time To Expand: ", time.Since(st))
}


func TestFindRoute_Handle(t *testing.T) {
	//sqare polygon
	//this creates large distance as the polygon corners are read as lat/lon pairs, but thats ok.
	var poly, poly2 Poly
	var polySet PolySet

	poly.AddCorner(0, 0)
	poly.AddCorner(0, 100)
	poly.AddCorner(50, 50)
	poly.AddCorner(100, 100)
	poly.AddCorner(100, 0)
	poly.AddCorner(0, 0)
	polySet.AddPoly(poly)

	poly2.AddCorner(20, 20)
	poly2.AddCorner(20, 40)
	poly2.AddCorner(40, 40)
	poly2.AddCorner(40, 20)
	poly2.AddCorner(20, 20)
	polySet.AddPoly(poly2)

	polySet.Print()

	var start, end Point
	x := 60.0
	start.Lon = 10
	start.Lat = 30
	end.Lon = 70
	end.Lat = x
	if !polySet.Verify() {
		fmt.Println("Polygon Set Failed Verification, should not proced with Shortest Path as it could give bad results")
	}
	r, err := ShortestPath(start, end, polySet)
	r.Print()
	Draw(&polySet, &r, start, end)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}


func TestTileSet_DumpDbAndCreateGenisisBlock(t *testing.T) {
	ts:=new(TileSet)
	ts.Init()
	err:=ts.DumpDbAndCreateGenisisBlock()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
}