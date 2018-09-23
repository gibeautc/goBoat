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

func (self *Route) Print() {
	fmt.Println("--Route--")
	fmt.Println("Start:")
	fmt.Println("X: ", self.start.X)
	fmt.Println("Y: ", self.start.Y)
	fmt.Println("Distance: ", self.Distance)
	fmt.Println("")

	for x := 0; x < self.Count; x++ {
		fmt.Println("Leg-", x)
		cP := self.Points[x]

		fmt.Println("X: ", cP.X)
		fmt.Println("Y: ", cP.Y)
		fmt.Println("D: ", cP.totalDistance)
		fmt.Println("")
	}

	fmt.Println("End:")
	fmt.Println("X: ", self.end.X)
	fmt.Println("Y: ", self.end.Y)

}

type Poly struct {
	y       []float64
	x       []float64
	corners int
}

func (self *Poly) AddCorner(x float64, y float64) {
	self.x = append(self.x, x)
	self.y = append(self.y, y)
	self.corners++
}

func (self *Poly) Print() {
	for x := 0; x < self.corners; x++ {
		fmt.Printf("Corner-> x: %f y: %f\n", self.x[x], self.y[x])
	}
	fmt.Println("")
}

type PolySet struct {
	count int
	poly  []Poly
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
		for y := 0; y < len(self.poly[x].x); y++ {
			if self.poly[x].x[y] < minX {
				minX = self.poly[x].x[y]
			}
			if self.poly[x].x[y] > maxX {
				maxX = self.poly[x].x[y]
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
		for y := 0; y < len(self.poly[x].y); y++ {
			if self.poly[x].y[y] < minY {
				minY = self.poly[x].y[y]
			}
			if self.poly[x].x[y] > maxY {
				maxY = self.poly[x].y[y]
			}
		}
	}
	return minY, maxY
}

func Draw(pS *PolySet, rT *Route, start Point, end Point) {
	minX, maxX := pS.MinMaxX()
	minY, maxY := pS.MinMaxY()
	imgSize := 2000
	r := image.Rect(0, 0, imgSize+int(float64(imgSize)*.25), imgSize+int(float64(imgSize)*.25))
	var sP, eP Point
	img := image.NewAlpha(r)

	for x := 0; x < pS.count; x++ {
		for y := 0; y < pS.poly[x].corners-1; y++ {
			sP.X = pS.poly[x].x[y]
			sP.Y = pS.poly[x].y[y]
			eP.X = pS.poly[x].x[y+1]
			eP.Y = pS.poly[x].y[y+1]
			drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
		}

		//last element back to start
		sP.X = pS.poly[x].x[0]
		sP.Y = pS.poly[x].y[0]
		eP.X = pS.poly[x].x[len(pS.poly[x].y)-1]
		eP.Y = pS.poly[x].y[len(pS.poly[x].y)-1]
		drawRtLine(img, sP, eP, minX, maxX, minY, maxY, imgSize, false)
	}

	SaveImage(img, "world.jpg")
	if rT != nil {
		if rT.Count == 0 {
			//direct route, only need to draw start to finish
			drawRtLine(img, start, end, minX, maxX, minY, maxY, imgSize, true)
			SaveImage(img, "route.jpg")
			return
		}
		//we have more nodes
		drawRtLine(img, start, rT.Points[0], minX, maxX, minY, maxY, imgSize, true)
		imageCount := 0
		SaveImage(img, "route"+strconv.Itoa(imageCount)+".jpg")
		imageCount++
		for x := 0; x < rT.Count-1; x++ {
			drawRtLine(img, rT.Points[x], rT.Points[x+1], minX, maxX, minY, maxY, imgSize, true)
			SaveImage(img, "route"+strconv.Itoa(imageCount)+".jpg")
			imageCount++
		}
		drawRtLine(img, rT.Points[rT.Count-1], end, minX, maxX, minY, maxY, imgSize, true)
		SaveImage(img, "route"+strconv.Itoa(imageCount)+".jpg")
	}

}
func drawRtLine(img draw.Image, start Point, end Point, minX float64, maxX float64, minY float64, maxY float64, imgSize int, isRoute bool) {
	var sP, eP image.Point
	sP.X = mapFloatToInt(start.X, minX, maxX, 0, imgSize)
	sP.Y = mapFloatToInt(start.Y, minY, maxY, 0, imgSize)
	eP.X = mapFloatToInt(end.X, minX, maxX, 0, imgSize)
	eP.Y = mapFloatToInt(end.Y, minY, maxY, 0, imgSize)
	fmt.Printf("Dawing line from x: %d y: %d TO x:%d y:%d\n", sP.X, sP.Y, eP.X, eP.Y)
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
	X, Y          float64
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
			if allPolys.poly[polyI].y[i] < testY && allPolys.poly[polyI].y[j] >= testY || allPolys.poly[polyI].y[j] < testY && allPolys.poly[polyI].y[i] >= testY {
				if allPolys.poly[polyI].x[i]+(testY-allPolys.poly[polyI].y[i])/(allPolys.poly[polyI].y[j]-allPolys.poly[polyI].y[i])*(allPolys.poly[polyI].x[j]-allPolys.poly[polyI].x[i]) < testX {
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

			sX := allPolys.poly[polyI].x[i] - testSX
			sY := allPolys.poly[polyI].y[i] - testSY
			eX := allPolys.poly[polyI].x[j] - testSX
			eY := allPolys.poly[polyI].y[j] - testSY
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
returns the list of internal nodes of the route. So excluding start and end points
Route also contains the number of nodes. If the count is zero, and no errors, then we have a straight
line route from start to end
*/

func pointDis(allPolys PolySet, from Point, to Point) float64 {
	if lineInPolygonSet(from.X, from.Y, to.X, to.Y, allPolys) {
		return calcDist(from.X, from.Y, to.X, to.Y)
	}
	return math.MaxFloat64
}

func ShortestPath(start Point, end Point, allPolys PolySet) (Route, error) {
	t := time.Now()
	var route Route
	pointList := make([]Point, 0)
	route.start = start
	route.end = end
	sX := start.X
	sY := start.Y
	eX := end.X
	eY := end.Y

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
		return route, nil
	}

	pointList = append(pointList, start)
	for polyI := 0; polyI < allPolys.count; polyI++ {
		for i := 0; i < allPolys.poly[polyI].corners; i++ {
			var tempP Point
			tempP.X = allPolys.poly[polyI].x[i]
			tempP.Y = allPolys.poly[polyI].y[i]
			tempP.totalDistance = math.MaxFloat64
			pointList = append(pointList, tempP)
		}
	}
	end.totalDistance = math.MaxFloat64
	pointList = append(pointList, end)
	pointList[0].totalDistance = 0.0

	for i := 0; i < len(pointList); i++ {
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

	//to get the route, have to work from end to
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
