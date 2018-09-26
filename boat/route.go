package boat

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"time"
)

type Route struct {
	Points   []Point
	Count    int
	Distance float64
	start    Point
	end      Point
}

func (self *Route) MinMaxX() (float64, float64) {
	var minX float64
	var maxX float64
	minX = 400.0
	maxX = -400.0
	for x := 0; x < self.Count; x++ {
		if self.Points[x].Lon > maxX {
			maxX = self.Points[x].Lon
		}
		if self.Points[x].Lon < minX {
			minX = self.Points[x].Lon
		}
	}
	return minX, maxX
}

func (self *Route) MinMaxY() (float64, float64) {
	var minY float64
	var maxY float64
	minY = 400.0
	maxY = -400.0
	for x := 0; x < self.Count; x++ {
		if self.Points[x].Lat > maxY {
			maxY = self.Points[x].Lat
		}
		if self.Points[x].Lat < minY {
			minY = self.Points[x].Lat
		}
	}
	return minY, maxY
}

func (self *Route) Print() {
	fmt.Println("--Route--")
	fmt.Println("Start:")
	fmt.Println("Lon: ", self.start.Lon)
	fmt.Println("Lat: ", self.start.Lat)
	fmt.Println("Distance: ", self.Distance)
	fmt.Println("")

	for x := 0; x < self.Count; x++ {
		fmt.Println("Leg-", x)
		cP := self.Points[x]

		fmt.Println("Lon: ", cP.Lon)
		fmt.Println("Lat: ", cP.Lat)
		fmt.Println("D: ", cP.totalDistance)
		fmt.Println("")
	}

	fmt.Println("End:")
	fmt.Println("Lon: ", self.end.Lon)
	fmt.Println("Lat: ", self.end.Lat)

}

type Poly struct {
	latLst  []float64
	lonLst  []float64
	corners int
}

func (self *Poly) AddCorner(lon float64, lat float64) {
	self.lonLst = append(self.lonLst, lon)
	self.latLst = append(self.latLst, lat)
	self.corners++

}

func (self *Poly) Print() {
	for x := 0; x < self.corners; x++ {
		fmt.Printf("Corner-> lonLst: %f latLst: %f\n", self.lonLst[x], self.latLst[x])
	}
	fmt.Println("")
}

func (self *Poly) Verify() bool {
	/*
		Need to verify that the polygon is valid.
		This means its closed and no line segents cross each other.
	*/

	//closed is the easiest to check, do that first
	if self.lonLst[0] != self.lonLst[len(self.lonLst)-1] || self.latLst[0] != self.latLst[len(self.latLst)-1] {
		fmt.Println("First and Last dont match")
		return false
	}

	//routing doesnt want duplicate points, so remove the first node
	self.lonLst = self.lonLst[1:]
	self.latLst = self.latLst[1:]
	if len(self.lonLst) != len(self.latLst) {
		//not sure how this would happen, but still bad
		return false
	}
	self.corners = len(self.lonLst)
	return true
}

type PolySet struct {
	count int
	poly  []Poly
}

func (self *PolySet) Verify() bool {
	for x := 0; x < self.count-1; x++ {
		if !self.poly[x].Verify() {
			fmt.Println("Polygon failed at index: ", x)
			return false
		}
	}
	fmt.Println("All Polygons Verified")
	return true
}

func (self *PolySet) Print() {
	for x := 0; x < self.count; x++ {
		fmt.Println("Polygon ", x)
		self.poly[x].Print()
	}
}

func (self *PolySet) MinMaxX() (float64, float64) {
	var minX float64
	var maxX float64
	minX = 400.0
	maxX = -400.0
	for x := 0; x < self.count; x++ {
		for y := 0; y < len(self.poly[x].lonLst); y++ {
			if self.poly[x].lonLst[y] < minX {
				minX = self.poly[x].lonLst[y]
			}
			if self.poly[x].lonLst[y] > maxX {
				maxX = self.poly[x].lonLst[y]
			}
		}
	}
	return minX, maxX
}

func (self *PolySet) MinMaxY() (float64, float64) {
	var minY float64
	var maxY float64
	minY = 400.0
	maxY = -400.0
	for x := 0; x < self.count; x++ {
		for y := 0; y < len(self.poly[x].latLst); y++ {
			if self.poly[x].latLst[y] < minY {
				minY = self.poly[x].latLst[y]
			}
			if self.poly[x].latLst[y] > maxY {
				maxY = self.poly[x].latLst[y]
			}
		}
	}
	return minY, maxY
}

func Draw(pS *PolySet, rT *Route, start Point, end Point) {
	//for this view, min and max should be for the Route, not the polySet
	//minX, maxX := pS.MinMaxX()
	//minY, maxY := pS.MinMaxY()
	minX, maxX := rT.MinMaxX()
	minY, maxY := rT.MinMaxY()
	imgSize := 2000
	r := image.Rect(0, 0, imgSize+int(float64(imgSize)*.25), imgSize+int(float64(imgSize)*.25))
	var sP, eP Point
	img := image.NewAlpha(r)

	for x := 0; x < pS.count; x++ {
		for y := 0; y < pS.poly[x].corners-1; y++ {
			sP.Lon = pS.poly[x].lonLst[y]
			sP.Lat = pS.poly[x].latLst[y]
			eP.Lon = pS.poly[x].lonLst[y+1]
			eP.Lat = pS.poly[x].latLst[y+1]
			drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
		}

		//last element back to start
		sP.Lon = pS.poly[x].lonLst[0]
		sP.Lat = pS.poly[x].latLst[0]
		eP.Lon = pS.poly[x].lonLst[len(pS.poly[x].latLst)-1]
		eP.Lat = pS.poly[x].latLst[len(pS.poly[x].latLst)-1]
		drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
	}

	SaveImage(img, "world.jpg")
	if rT != nil {
		if rT.Count == 0 {
			//direct Route, only need to draw start to finish
			drawRtLine(img, start, end, minX, maxX, minY, maxY, imgSize, true)
			SaveImage(img, "Route.jpg")
			return
		}
		//we have more nodes
		drawRtLine(img, start, rT.Points[0], minX, maxX, minY, maxY, imgSize, true)
		imageCount := 0
		SaveImage(img, "Route"+strconv.Itoa(imageCount)+".jpg")
		imageCount++
		for x := 0; x < rT.Count-1; x++ {
			drawRtLine(img, rT.Points[x], rT.Points[x+1], minX, maxX, minY, maxY, imgSize, true)
			SaveImage(img, "Route"+strconv.Itoa(imageCount)+".jpg")
			imageCount++
		}
		drawRtLine(img, rT.Points[rT.Count-1], end, minX, maxX, minY, maxY, imgSize, true)
		SaveImage(img, "Route"+strconv.Itoa(imageCount)+".jpg")
	}

}
func DrawWorld(pS *PolySet) {
	minX, maxX := pS.MinMaxX()
	minY, maxY := pS.MinMaxY()
	imgSize := 2000
	r := image.Rect(0, 0, imgSize+int(float64(imgSize)*.25), imgSize+int(float64(imgSize)*.25))
	var sP, eP Point
	img := image.NewAlpha(r)

	for x := 0; x < pS.count; x++ {
		for y := 0; y < pS.poly[x].corners-1; y++ {
			sP.Lon = pS.poly[x].lonLst[y]
			sP.Lat = pS.poly[x].latLst[y]
			eP.Lon = pS.poly[x].lonLst[y+1]
			eP.Lat = pS.poly[x].latLst[y+1]
			drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
		}

		//last element back to start
		sP.Lon = pS.poly[x].lonLst[0]
		sP.Lat = pS.poly[x].latLst[0]
		eP.Lon = pS.poly[x].lonLst[len(pS.poly[x].latLst)-1]
		eP.Lat = pS.poly[x].latLst[len(pS.poly[x].latLst)-1]
		drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
	}

	SaveImage(img, "world.jpg")
}
func drawRtLine(img draw.Image, start Point, end Point, minX float64, maxX float64, minY float64, maxY float64, imgSize int, isRoute bool) {
	var sP, eP image.Point
	sP.X = mapFloatToInt(start.Lon, minX, maxX, 0, imgSize)
	sP.Y = mapFloatToInt(start.Lat, minY, maxY, 0, imgSize)
	eP.X = mapFloatToInt(end.Lon, minX, maxX, 0, imgSize)
	eP.Y = mapFloatToInt(end.Lat, minY, maxY, 0, imgSize)
	//fmt.Printf("Drawing line from lonLst: %d latLst: %d TO lonLst:%d latLst:%d\n", sP.X, sP.Y, eP.X, eP.Y)
	if isRoute {
		c := color.Alpha16{A: 0xFF0F}
		drawLine(img, sP, eP, c)
	} else {
		drawLine(img, sP, eP, color.White)
	}

}

func (self *PolySet) AddPoly(p Poly) {
	self.poly = append(self.poly, p)
	self.count++
}

type Point struct {
	Lon, Lat      float64
	totalDistance float64
	prev          int
}

func calcDist(sX float64, sY float64, eX float64, eY float64) float64 {
	//todo change to use a more accurate distance calculation
	eX -= sX
	eY -= sY
	return math.Sqrt(eX*eX + eY*eY)
}

func pointInPolygonSet(testX float64, testY float64, allPolys PolySet) bool {
	oddNodes := false
	var j int
	for polyI := 0; polyI < allPolys.count; polyI++ {
		for i := 0; i < allPolys.poly[polyI].corners; i++ {
			j = i + 1
			if j == allPolys.poly[polyI].corners {
				j = 0
			}
			if allPolys.poly[polyI].latLst[i] < testY && allPolys.poly[polyI].latLst[j] >= testY || allPolys.poly[polyI].latLst[j] < testY && allPolys.poly[polyI].latLst[i] >= testY {
				if allPolys.poly[polyI].lonLst[i]+(testY-allPolys.poly[polyI].latLst[i])/(allPolys.poly[polyI].latLst[j]-allPolys.poly[polyI].latLst[i])*(allPolys.poly[polyI].lonLst[j]-allPolys.poly[polyI].lonLst[i]) < testX {
					oddNodes = !oddNodes
				}
			}
		}
	}
	return oddNodes
}

func lineInPolygonSet(testSX float64, testSY float64, testEX float64, testEY float64, allPolys PolySet) bool {
	testEX -= testSX
	testEY -= testSY
	dist := math.Sqrt(testEX*testEX + testEY*testEY)
	theCos := testEX / dist
	theSin := testEY / dist
	var j int
	for polyI := 0; polyI < allPolys.count; polyI++ {
		for i := 0; i < allPolys.poly[polyI].corners; i++ {
			j = i + 1
			if j == allPolys.poly[polyI].corners {
				j = 0
			}

			sX := allPolys.poly[polyI].lonLst[i] - testSX
			sY := allPolys.poly[polyI].latLst[i] - testSY
			eX := allPolys.poly[polyI].lonLst[j] - testSX
			eY := allPolys.poly[polyI].latLst[j] - testSY
			if sX == 0. && sY == 0. && eX == testEX && eY == testEY || eX == 0. && eY == 0. && sX == testEX && sY == testEY {
				return true
			}

			rotSX := sX*theCos + sY*theSin
			rotSY := sY*theCos - sX*theSin
			rotEX := eX*theCos + eY*theSin
			rotEY := eY*theCos - eX*theSin
			if rotSY < 0. && rotEY > 0. || rotEY < 0. && rotSY > 0. {
				crossX := rotSX + (rotEX-rotSX)*(0.-rotSY)/(rotEY-rotSY)
				if crossX >= 0. && crossX <= dist {
					return false
				}
			}

			if rotSY == 0. && rotEY == 0. && (rotSX >= 0. || rotEX >= 0.) && (rotSX <= dist || rotEX <= dist) && (rotSX < 0. || rotEX < 0. || rotSX > dist || rotEX > dist) {
				return false
			}
		}
	}
	return pointInPolygonSet(testSX+testEX/2., testSY+testEY/2., allPolys)
}

/*
returns the list of internal nodes of the Route. So excluding start and end points
Route also contains the number of nodes. If the count is zero, and no errors, then we have a straight
line Route from start to end
*/

func pointDis(allPolys PolySet, from Point, to Point) float64 {
	if lineInPolygonSet(from.Lon, from.Lat, to.Lon, to.Lat, allPolys) {
		//return calcDist(from.Lon, from.Lat, to.Lon, to.Lat)
		dist, _ := DistanceBetween(from.Lon, from.Lat, to.Lon, to.Lat)
		return float64(dist)
	}
	return math.MaxFloat64
}

func ShortestPath(start Point, end Point, allPolys PolySet) (Route, error) {
	t := time.Now()
	var route Route
	pointList := make([]Point, 0)
	route.start = start
	route.end = end
	sX := start.Lon
	sY := start.Lat
	eX := end.Lon
	eY := end.Lat

	//check to make sure start and end points are inside polys
	//not sure if this is really needed as our starting location will be our current location, and we know we get to where we are.
	if !pointInPolygonSet(sX, sY, allPolys) {
		return route, errors.New("start point not inside polys")
	}

	//todo may need this to be more of a warning as scan data may "appear" to put this point in an unreachable place
	if !pointInPolygonSet(eX, eY, allPolys) {
		return route, errors.New("end point not inside polys")
	}

	//check if straight line solution works
	if lineInPolygonSet(sX, sY, eX, eY, allPolys) {
		route.Count = 0
		fmt.Print("StaightLine Route Works")
		return route, nil
	}

	pointList = append(pointList, start)
	for polyI := 0; polyI < allPolys.count; polyI++ {
		for i := 0; i < allPolys.poly[polyI].corners; i++ {
			var tempP Point
			tempP.Lon = allPolys.poly[polyI].lonLst[i]
			tempP.Lat = allPolys.poly[polyI].latLst[i]
			tempP.totalDistance = math.MaxFloat64
			pointList = append(pointList, tempP)
		}
	}
	end.totalDistance = math.MaxFloat64
	pointList = append(pointList, end)
	pointList[0].totalDistance = 0.0

	for i := 0; i < len(pointList); i++ {
		fmt.Println(float64(i) / float64(len(pointList)))
		if pointList[i].totalDistance == math.MaxFloat64 {
			continue
		}
		for j := 0; j < len(pointList); j++ {
			if i == j {
				continue
			}
			dist := pointDis(allPolys, pointList[i], pointList[j]) + pointList[i].totalDistance

			if dist < pointList[j].totalDistance {
				pointList[j].totalDistance = dist
				pointList[j].prev = i

			}
		}
	}

	//to get the Route, have to work from end to
	backwardsRoute := make([]Point, 0)
	index := len(pointList) - 1
	for true {
		index = pointList[index].prev
		if index == 0 {
			break
		}
		backwardsRoute = append(backwardsRoute, pointList[index])

	}

	for index := len(backwardsRoute) - 1; index >= 0; index-- {
		route.Points = append(route.Points, backwardsRoute[index])
	}
	route.Count = len(route.Points)
	route.Distance = pointList[len(pointList)-1].totalDistance
	fmt.Println("Time to Calculate Route: ", time.Since(t))
	return route, nil
}

//going to implement diksters.... and actualy swap points in pointList
func ShortestPath2(start Point, end Point, allPolys PolySet) (Route, error) {
	t := time.Now()
	var route Route
	pointList := make([]Point, 0)
	route.start = start
	route.end = end
	sX := start.Lon
	sY := start.Lat
	eX := end.Lon
	eY := end.Lat

	//check to make sure start and end points are inside polys
	//not sure if this is really needed as our starting location will be our current location, and we know we get to where we are.
	if !pointInPolygonSet(sX, sY, allPolys) {
		return route, errors.New("start point not inside polys")
	}

	//todo may need this to be more of a warning as scan data may "appear" to put this point in an unreachable place
	if !pointInPolygonSet(eX, eY, allPolys) {
		return route, errors.New("end point not inside polys")
	}

	//check if straight line solution works
	if lineInPolygonSet(sX, sY, eX, eY, allPolys) {
		route.Count = 0
		fmt.Print("StaightLine Route Works")
		return route, nil
	}

	pointList = append(pointList, start)
	for polyI := 0; polyI < allPolys.count; polyI++ {
		for i := 0; i < allPolys.poly[polyI].corners; i++ {
			var tempP Point
			tempP.Lon = allPolys.poly[polyI].lonLst[i]
			tempP.Lat = allPolys.poly[polyI].latLst[i]
			tempP.totalDistance = math.MaxFloat64
			pointList = append(pointList, tempP)
		}
	}
	end.totalDistance = math.MaxFloat64
	pointList = append(pointList, end)
	pointList[0].totalDistance = 0.0

	tc:=0
	bestJ:=0
	bestDist:=math.MaxFloat64
	for bestJ!=len(pointList)-1{
		for j := tc+1; j < len(pointList); j++ {
			dist := pointDis(allPolys, pointList[tc], pointList[j]) + pointList[tc].totalDistance

			if dist < pointList[j].totalDistance {
				pointList[j].totalDistance = dist
				pointList[j].prev = tc
				bestJ=j
				bestDist=dist
			}
		}

	}
	_=bestDist
	//to get the Route, have to work from end to
	backwardsRoute := make([]Point, 0)
	index := len(pointList) - 1
	for true {
		index = pointList[index].prev
		if index == 0 {
			break
		}
		backwardsRoute = append(backwardsRoute, pointList[index])

	}

	for index := len(backwardsRoute) - 1; index >= 0; index-- {
		route.Points = append(route.Points, backwardsRoute[index])
	}
	route.Count = len(route.Points)
	route.Distance = pointList[len(pointList)-1].totalDistance
	fmt.Println("Time to Calculate Route: ", time.Since(t))
	return route, nil
}

func mapFloatToInt(input float64, inMin float64, inMax float64, outMin int, outMax int) int {
	slope := 1.0 * (float64(outMax) - float64(outMin)) / (inMax - inMin)
	output := float64(outMin) + slope*(input-inMin)
	return int(output)
}
