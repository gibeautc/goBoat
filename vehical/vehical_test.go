package vehical

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func TestGetCords2(t *testing.T) {
	//finding cords for genisis block
	lat, lon := GetCords(44.616429, -123.072626, 2048, 180)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//that will be the southCenter cords
	lat, lon = GetCords(44.6159617, -123.072626, 2048, 270)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//this will be the SW Corner
	//Lat:  44.61596169811975
	//Lon:  -123.07328247537501
	lat, lon = GetCords(44.61596169811975, -123.07328247537501, 4096, 90)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//this will be SE Corner
	//Lat:  44.61596169059874
	//Lon:  -123.07196952462509
	lat, lon = GetCords(44.61596169059874, -123.07196952462509, 4096, 0)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//this will be NE Corner
	//Lat:  44.61689628886895
	//Lon:  -123.07196952462509
	lat, lon = GetCords(44.61689628886895, -123.07196952462509, 4096, 270)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//this will be NW Corner
	//Lat:  44.616896281347714
	//Lon:  -123.07328249650676
	lat, lon = GetCords(44.616896281347714, -123.07328249650676, 4096, 180)
	fmt.Println("Lat: ", lat)
	fmt.Println("Lon: ", lon)
	//this should be SW corner again
	//this is what we will check
	assert.Equal(t, true, math.Abs(lat-44.61596169811975) < 44.61596169811975*.0001)
	assert.Equal(t, true, math.Abs(lon+123.07328247537501) < 123.07328247537501*.0001)
}

func TestGetCords(t *testing.T) {
	lat, lon := GetCords(44.0, -123.0, 538976, -35)
	fmt.Println("Lat")
	fmt.Println("Expected: ", 44.1)
	fmt.Println("Got: ", lat)
	fmt.Println("")
	fmt.Println("Lon")
	fmt.Println("Expected: ", -123.1)
	fmt.Println("Got: ", lon)

	assert.Equal(t, true, math.Abs(lat-44.1) < 44.1*.0001)
	assert.Equal(t, true, math.Abs(lon+123.1) < 123.1*.0001)
}

func TestDistanceBetween(t *testing.T) {
	dist, direction := DistanceBetween(44.0, -123.0, 44.1, -123.1)
	//13.69 Km --> 538976.38
	//324 Deg
	delta := math.Abs(float64(dist - 538976))
	assert.Equal(t, true, delta < float64(dist)*.01)
	assert.Equal(t, 324, direction)
}

func TestTile_Expand(t *testing.T) {

	tile := NewTile()
	fmt.Println("New Tile Created")
	//todo make a small image

	st := time.Now()
	resp := true
	tile.Id = 0
	for resp {
		fmt.Println("Expanding from: ", tile.Size)
		resp = tile.Expand()
		fmt.Println("Done Expanding")
	}
	tile.SaveImage()
	fmt.Println("Total Time To Expand: ", time.Since(st))
}

func TestTile_Compress(t *testing.T) {

	tile := NewTile()
	fmt.Println("New Tile Created")
	tile.Id = 0
	st := time.Now()
	resp := true
	for resp {
		fmt.Println("Compressing from: ", tile.Size)
		resp = tile.Compress()
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
	ts := new(TileSet)
	ts.Init()
	err := ts.DumpDbAndCreateGenisisBlock(true)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestTileSet_GetOldestToCompress(t *testing.T) {
	/*
		not actually creating files for these, so size in DB is bogus
	*/
	ts := new(TileSet)
	ts.Init()
	err := ts.DumpDbAndCreateGenisisBlock(false)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	for x := 0; x < 10; x++ {
		tl := NewTile()
		tl.Id, err = ts.GetNewTileID()
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
		ts.updateTileToDB(*tl, 0)
	}
	_, err = ts.conn.Exec("UPDATE tiles set onDisk=1") //fake them on disk
	_, err = ts.conn.Exec("UPDATE tiles set comp=1 where id=1")
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	id, err := ts.GetOldestToCompress()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	//since we changed the size of id=1 to be 1 (fully compressed) the next oldest should be id=2
	assert.Equal(t, 2, id, "Oldest Should be id 2")
}

func TestTileSet_CheckMemoryAndCompress(t *testing.T) {
	ts := new(TileSet)
	err := ts.Init()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}

	for x := 0; x < 10; x++ {
		tl := NewTile()
		tl.Id, err = ts.GetNewTileID()
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
		}
		ts.Pickle(*tl)
		fmt.Println("Adding Tile: ", x)
		ts.updateTileToDB(*tl, 0)
		ts.CheckMemoryAndCompress()
	}
}

func TestGetDiskSpaceOfPathMB(t *testing.T) {
	/*
		Should return the size of all files in path or a specific file if that is what is defined in path
		return is a float in MB
	*/
	mb := GetDiskSpaceOfPathMB("/home/chadg/logMon")
	fmt.Println("DiskSize: ", mb)
}

func TestTile_UnPickle(t *testing.T) {
	ts:=new(TileSet)
	ts.Init()
	tile := NewTile()
	tile.Id = 0
	tile.Compress()
	ts.Pickle(*tile)


	tile,err := ts.UnPickle(0)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(tile.Size)
}

func TestTileSet_GetIdByPoint(t *testing.T) {
	ts := new(TileSet)
	err := ts.Init()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	var p Point
	p.Lat = 44.6169
	p.Lon = -123.072815
	id, err := ts.GetIdByPoint(p)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	assert.Equal(t, 1, id, "Should be inside genisis block")
}

func TestTile_GetPixelByCords(t *testing.T) {
	tile := NewTile()
	tile.Bounds.NW.Lat = 44.0
	tile.Bounds.NW.Lon = -123.0
	tile.Bounds.NE.Lat = 44.0
	tile.Bounds.NE.Lon = -124.0
	tile.Bounds.SE.Lat = 43.0
	tile.Bounds.SE.Lon = -124.0
	tile.Bounds.SW.Lat = 43.0
	tile.Bounds.SW.Lon = -123.0

	var p Point
	p.Lat = 43.5
	p.Lon = -123.5
	x, y := tile.GetPixelByCords(p)
	fmt.Println("x: ", x)
	fmt.Println("y: ", y)
	assert.Equal(t, 2048, x, "Should be middle")
	assert.Equal(t, 2048, y, "Should be middle")
}

func TestTile_SaveImage(t *testing.T) {
	tile := NewTile()
	tile.Bounds.NW.Lat = 44.0
	tile.Bounds.NW.Lon = -123.0
	tile.Bounds.NE.Lat = 44.0
	tile.Bounds.NE.Lon = -124.0
	tile.Bounds.SE.Lat = 43.0
	tile.Bounds.SE.Lon = -124.0
	tile.Bounds.SW.Lat = 43.0
	tile.Bounds.SW.Lon = -123.0
	tile.Id = 0

	var cur, obj Point
	cur.Lat = 43.1
	cur.Lon = -123.1
	obj.Lat = 43.7
	obj.Lon = -123.7
	st := time.Now()
	err := tile.AddDistanceData(cur, obj)

	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println("AddDistanceData Time: ", time.Since(st))
	st = time.Now()
	err = tile.SaveImage()
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println("SaveImage Time: ", time.Since(st))
}
