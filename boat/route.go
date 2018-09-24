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

func (self *Poly) Verify() bool {
	/*
		Need to verify that the polygon is valid.
		This means its closed and no line segents cross each other.
	*/

	//closed is the easiest to check, do that first
	if self.x[0] != self.x[len(self.x)-1] || self.y[0] != self.y[len(self.y)-1] {
		fmt.Println("First and Last dont match")
		return false
	}

	for i := 0; i < self.corners-2; i++ {
		for j := 0; j < self.corners-2; j++ {
			if i == j {
				continue
			}
			if linesCross(self.x[i], self.y[i], self.x[i+1], self.y[i+1], self.x[j], self.y[j], self.x[j+1], self.y[j+1]) {
				return false
			}
		}
	}
	//routing doesnt want duplicate points, so remove the first node
	self.x = self.x[1:]
	self.y = self.y[1:]
	return true
}

// Given three colinear points p, q, r, the function checks if
// point q lies on line segment 'pr'
func onSegment(p Point, q Point, r Point) bool {
	if q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) && q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y) {
		return true
	}
	return false
}

// To find orientation of ordered triplet (p, q, r).
// The function returns following values
// 0 --> p, q and r are colinear
// 1 --> Clockwise
// 2 --> Counterclockwise
func orientation(p Point, q Point, r Point) int {
	// See https://www.geeksforgeeks.org/orientation-3-ordered-points/
	// for details of below formula.
	val := (q.Y-p.Y)*(r.X-q.X) - (q.X-p.X)*(r.Y-q.Y)

	if val == 0 {
		return 0
	} // colinear

	if val > 0 {
		return 1
	}
	return 2 // clock or counterclock wise
}

func linesCross(l1Sx float64, l1Sy float64, l1Ex float64, l1Ey float64, l2Sx float64, l2Sy float64, l2Ex float64, l2Ey float64) bool {

	var p1, q1, p2, q2 Point
	p1.X = l1Sx
	p1.Y = l1Sy

	q1.X = l1Ex
	q1.Y = l1Ex

	p2.X = l2Sx
	p2.Y = l2Sy

	q2.X = l2Ex
	q2.Y = l2Ey

	//bool doIntersect(Point p1, Point q1, Point p2, Point q2)

	// Find the four orientations needed for general and
	// special cases
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	// General case
	//if o1 != o2 && o3 != o4 {
	//	fmt.Println("General Case")
	//	return true
	//}

	// Special Cases
	// p1, q1 and p2 are colinear and p2 lies on segment p1q1
	if o1 == 0 && onSegment(p1, p2, q1) {
		fmt.Println("General Case")
		return true
	}

	// p1, q1 and q2 are colinear and q2 lies on segment p1q1
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}

	// p2, q2 and p1 are colinear and p1 lies on segment p2q2
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}

	// p2, q2 and q1 are colinear and q1 lies on segment p2q2
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false // Doesn't fall in any of the above cases

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
