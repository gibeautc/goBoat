package boat

import (
	"math"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strconv"
)

type Route struct{
	Points []Point
	Count int
	start Point
	end Point
}

func (self *Route) Print(){
	fmt.Println("--Route--")
	fmt.Println("Start:")
	fmt.Println("X: ",self.start.X)
	fmt.Println("Y: ",self.start.Y)
	fmt.Println("")

	for x:=0;x<self.Count;x++{
		fmt.Println("Leg-",x)
		cP:=self.Points[x]

		fmt.Println("X: ",cP.X)
		fmt.Println("Y: ",cP.Y)
		fmt.Println("D: ",cP.totalDistance)
		fmt.Println("")
	}

	fmt.Println("End:")
	fmt.Println("X: ",self.end.X)
	fmt.Println("Y: ",self.end.Y)

}

type Poly struct{
	y []float64
	x []float64
	corners int
}

func (self *Poly) AddCorner(x float64,y float64){
	self.x=append(self.x,x)
	self.y=append(self.y,y)
	self.corners++
}

func (self *Poly) Print(){
	for x:=0;x<self.corners;x++{
		fmt.Printf("Corner-> x: %f y: %f\n",self.x[x],self.y[x])
	}
	fmt.Println("")
}


type PolySet struct{
	count int
	poly []Poly
}

func (self *PolySet) Print(){
	for x:=0;x<self.count;x++{
		fmt.Println("Polygon ",x)
		self.poly[x].Print()
	}
}

func (self *PolySet) MinMaxX() (float64,float64){
	var minX float64
	var maxX float64
	minX=400.0
	maxX=-400.0
	for x:=0;x<self.count;x++{
		for y:=0;y<len(self.poly[x].x);y++{
			if self.poly[x].x[y]<minX{
				minX=self.poly[x].x[y]
			}
			if self.poly[x].x[y]>maxX{
				maxX=self.poly[x].x[y]
			}
		}
	}
	return minX,maxX
}

func (self *PolySet) MinMaxY() (float64,float64){
	var minY float64
	var maxY float64
	minY=400.0
	maxY=-400.0
	for x:=0;x<self.count;x++{
		for y:=0;y<len(self.poly[x].y);y++{
			if self.poly[x].y[y]<minY{
				minY=self.poly[x].y[y]
			}
			if self.poly[x].x[y]>maxY{
				maxY=self.poly[x].y[y]
			}
		}
	}
	return minY,maxY
}


func Draw(pS *PolySet,rT *Route,start Point,end Point){
	minX,maxX:=pS.MinMaxX()
	minY,maxY:=pS.MinMaxY()
	imgSize:=2000
	r:=image.Rect(0,0,imgSize+int(float64(imgSize)*.25),imgSize+int(float64(imgSize)*.25))
	var sP,eP Point
	img:=image.NewAlpha(r)

	for x:=0;x<pS.count;x++{
		for y:=0;y<pS.poly[x].corners-1;y++{
			sP.X=pS.poly[x].x[y]
			sP.Y=pS.poly[x].y[y]
			eP.X=pS.poly[x].x[y+1]
			eP.Y=pS.poly[x].y[y+1]
			drawRtLine(img,sP,eP,minX,maxX,minY,maxY,imgSize,false)
		}

		//last element back to start
		sP.X=pS.poly[x].x[0]
		sP.Y=pS.poly[x].y[0]
		eP.X=pS.poly[x].x[len(pS.poly[x].y)-1]
		eP.Y=pS.poly[x].y[len(pS.poly[x].y)-1]
		drawRtLine(img,sP,eP,minX,maxX,minY,maxY,imgSize,false)
	}


	
	SaveImage(img,"world.jpg")
	if rT!=nil{
		if rT.Count==0{
			//direct route, only need to draw start to finish
			drawRtLine(img,start,end,minX,maxX,minY,maxY,imgSize,true)
			SaveImage(img,"route.jpg")
			return
		}
		//we have more nodes
		drawRtLine(img,start,rT.Points[0],minX,maxX,minY,maxY,imgSize,true)
		imageCount:=0
		SaveImage(img,"route"+strconv.Itoa(imageCount)+".jpg")
		imageCount++
		for x:=0;x<rT.Count-1;x++{
			drawRtLine(img,rT.Points[x],rT.Points[x+1],minX,maxX,minY,maxY,imgSize,true)
			SaveImage(img,"route"+strconv.Itoa(imageCount)+".jpg")
			imageCount++
		}
		drawRtLine(img,rT.Points[rT.Count-1],end,minX,maxX,minY,maxY,imgSize,true)
		SaveImage(img,"route"+strconv.Itoa(imageCount)+".jpg")
	}


}
func drawRtLine(img draw.Image,start Point,end Point,minX float64,maxX float64,minY float64,maxY float64,imgSize int,isRoute bool){
	var sP,eP image.Point
	sP.X=mapFloatToInt(start.X,minX,maxX,0,imgSize)
	sP.Y=mapFloatToInt(start.Y,minY,maxY,0,imgSize)
	eP.X=mapFloatToInt(end.X,minX,maxX,0,imgSize)
	eP.Y=mapFloatToInt(end.Y,minY,maxY,0,imgSize)
	fmt.Printf("Dawing line from x: %d y: %d TO x:%d y:%d\n",sP.X,sP.Y,eP.X,eP.Y)
	if isRoute{
		c:=color.Alpha16{A:0xFF0F}
		drawLine(img,sP,eP,c)
	}else{
		drawLine(img,sP,eP,color.White)
	}

}

func(self *PolySet) AddPoly(p Poly){
	self.poly=append(self.poly,p)
	self.count++
}

type Point struct{
	X,Y float64
	totalDistance float64
	prev int
}

func calcDist(sX float64,sY float64,eX float64,eY float64) float64{
	eX-=sX
	eY-=sY
	return math.Sqrt(eX*eX+eY*eY)
}

func swapPoints(a *Point,b *Point){
	swap:=*a
	*a=*b
	*b=swap
}



func pointInPolygonSet(testX float64,testY float64,allPolys PolySet) bool{
	oddNodes:=false
	var j int
	for polyI:=0;polyI<allPolys.count;polyI++{
		for i:=0;i<allPolys.poly[polyI].corners;i++{
			j=i+1
			if j==allPolys.poly[polyI].corners{
				j=0
			}
			if allPolys.poly[polyI].y[i]<testY && allPolys.poly[polyI].y[j]>=testY || allPolys.poly[polyI].y[j]< testY && allPolys.poly[polyI].y[i]>=testY{
				if 	allPolys.poly[polyI].x[i]+(testY-allPolys.poly[polyI].y[i])/ (allPolys.poly[polyI].y[j]-allPolys.poly[polyI].y[i])*(allPolys.poly[polyI].x[j]-allPolys.poly[polyI].x[i])<testX{
					oddNodes=!oddNodes				
				}
			} 
		}
	}
	return oddNodes
}

func lineInPolygonSet(testSX float64,testSY float64,testEX float64,testEY float64,allPolys PolySet) bool {
	testEX -= testSX
	testEY -= testSY
	dist := math.Sqrt(testEX*testEX + testEY*testEY)
	theCos := testEX / dist;
	theSin := testEY / dist;
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
					return false;
				}
			}

			if rotSY == 0. && rotEY == 0. && (rotSX >= 0. || rotEX >= 0.) && (rotSX <= dist || rotEX <= dist) && (rotSX < 0. || rotEX < 0. || rotSX > dist || rotEX > dist) {
				return false
			}
		}
	}
	return pointInPolygonSet(testSX+testEX/2.,testSY+testEY/2.,allPolys)
}

/*
returns the list of internal nodes of the route. So excluding start and end points
Route also contains the number of nodes. If the count is zero, and no errors, then we have a straight
line route from start to end
 */
func ShortestPath(start Point,end Point,allPolys PolySet) (Route,error){
	var route Route
	pointList:=make([]Point,0)	
	var treeCount int
	route.start=start
	route.end=end
	sX:=start.X
	sY:=start.Y
	eX:=end.X
	eY:=end.Y
	
	//check to make sure start and end points are inside polys
	if !pointInPolygonSet(sX,sY,allPolys){
		return route,errors.New("start point not inside polys")
	}

	if !pointInPolygonSet(eX,eY,allPolys){
		return route,errors.New("end point not inside polys")
	}

	//check if straight line solution works
	if lineInPolygonSet(sX,sY,eX,eY,allPolys){
		route.Count=0
		return route,nil
	}


	pointList=append(pointList,start)
	for polyI:=0;polyI<allPolys.count;polyI++{
		for i:=0;i<allPolys.poly[polyI].corners;i++{
			var tempP Point
			tempP.X=allPolys.poly[polyI].x[i]
			tempP.Y=allPolys.poly[polyI].y[i]
			pointList=append(pointList,tempP)
		}
	}
	//not sure if if matters yet that the endpoint is actually put at the end, but that is the way the original code was
	pointList=append(pointList,end)
	treeCount=1
	pointList[0].totalDistance=0.0


	bestJ:=0
	bestI:=0

	for bestJ!=len(pointList)-1{
		bestDist:=math.MaxFloat64
		for i:=0;i<treeCount;i++{
			for j:=treeCount;j<len(pointList);j++{

				if lineInPolygonSet(pointList[i].X,pointList[i].Y,pointList[j].X,pointList[j].Y,allPolys){
					newDist:=pointList[treeCount].totalDistance+calcDist(pointList[i].X,pointList[i].Y,pointList[j].X,pointList[j].X)
					if newDist<bestDist{
						bestDist=newDist
						bestI=treeCount
						bestJ=j
					}
				}
			}

			if bestDist==math.MaxFloat64{
				return route,errors.New("noRoutePossible")
			}

		}
		pointList[bestJ].prev=bestI
		pointList[bestJ].totalDistance=bestDist
		fmt.Println("bestJ: ",bestJ)
		fmt.Println("treeCount: ",treeCount)
		fmt.Println("pointListLen: ",len(pointList))
		swapPoints(&pointList[bestJ],&pointList[treeCount])
		treeCount++
	}


	for cnt:=1;cnt<len(pointList);cnt++{
		curPoint:=pointList[cnt]
		if curPoint.X==end.X && curPoint.Y==end.Y{
			break
		}
		route.Points=append(route.Points,curPoint)
	}
	route.Count=len(route.Points)
	return route,nil
}


func mapFloatToInt(input float64,inMin float64,inMax float64,outMin int,outMax int) int{
	slope:= 1.0 * (float64(outMax) - float64(outMin)) / (inMax - inMin)
	output:= float64(outMin) + slope * (input - inMin)
	return int(output)
}


/*
//  Public-domain code by Darel Rex Finley, 2006.



//  (This function automatically knows that enclosed polygons are "no-go"
//  areas.)

boolean pointInPolygonSet(double testX, double testY, polySet allPolys) {

  bool  oddNodes=NO ;
  int   polyI, i, j ;

  for (polyI=0; polyI<allPolys.count; polyI++) {
    for (i=0;    i< allPolys.poly[polyI].corners; i++) {
      j=i+1; if (j==allPolys.poly[polyI].corners) j=0;
      if   ( allPolys.poly[polyI].y[i]< testY
      &&     allPolys.poly[polyI].y[j]>=testY
      ||     allPolys.poly[polyI].y[j]< testY
      &&     allPolys.poly[polyI].y[i]>=testY) {
        if ( allPolys.poly[polyI].x[i]+(testY-allPolys.poly[polyI].y[i])
        /   (allPolys.poly[polyI].y[j]       -allPolys.poly[polyI].y[i])
        *   (allPolys.poly[polyI].x[j]       -allPolys.poly[polyI].x[i])<testX) {
          oddNodes=!oddNodes; }}}}

  return oddNodes; }


//  This function should be called with the full set of *all* relevant polygons.
//  (The algorithm automatically knows that enclosed polygons are “no-go”
//  areas.)
//
//  Note:  As much as possible, this algorithm tries to return YES when the
//         test line-segment is exactly on the border of the polygon, particularly
//         if the test line-segment *is* a side of a polygon.

bool lineInPolygonSet(
double testSX, double testSY, double testEX, double testEY, polySet allPolys) {

  double  theCos, theSin, dist, sX, sY, eX, eY, rotSX, rotSY, rotEX, rotEY, crossX ;
  int     i, j, polyI ;

  testEX-=testSX;
  testEY-=testSY; dist=sqrt(testEX*testEX+testEY*testEY);
  theCos =testEX/ dist;
  theSin =testEY/ dist;

  for (polyI=0; polyI<allPolys.count; polyI++) {
    for (i=0;    i< allPolys.poly[polyI].corners; i++) {
      j=i+1; if (j==allPolys.poly[polyI].corners) j=0;

      sX=allPolys.poly[polyI].x[i]-testSX;
      sY=allPolys.poly[polyI].y[i]-testSY;
      eX=allPolys.poly[polyI].x[j]-testSX;
      eY=allPolys.poly[polyI].y[j]-testSY;
      if (sX==0. && sY==0. && eX==testEX && eY==testEY
      ||  eX==0. && eY==0. && sX==testEX && sY==testEY) {
        return YES; }

      rotSX=sX*theCos+sY*theSin;
      rotSY=sY*theCos-sX*theSin;
      rotEX=eX*theCos+eY*theSin;
      rotEY=eY*theCos-eX*theSin;
      if (rotSY<0. && rotEY>0.
      ||  rotEY<0. && rotSY>0.) {
        crossX=rotSX+(rotEX-rotSX)*(0.-rotSY)/(rotEY-rotSY);
        if (crossX>=0. && crossX<=dist) return NO; }

      if ( rotSY==0.   && rotEY==0.
      &&  (rotSX>=0.   || rotEX>=0.  )
      &&  (rotSX<=dist || rotEX<=dist)
      &&  (rotSX< 0.   || rotEX< 0.
      ||   rotSX> dist || rotEX> dist)) {
        return NO; }}}

  return pointInPolygonSet(testSX+testEX/2.,testSY+testEY/2.,allPolys); }



double calcDist(double sX, double sY, double eX, double eY) {
  eX-=sX; eY-=sY; return sqrt(eX*eX+eY*eY); }



void swapPoints(point *a, point *b) {
  point swap=*a; *a=*b; *b=swap; }



bool shortestPath(double sX, double sY, double eX, double eY, polySet allPolys,
double *solutionX, double *solutionY, int *solutionNodes) {

  #define  INF  9999999.     //  (larger than total solution dist could ever be)

  point  pointList[1000] ;   //  (enough for all polycorners plus two)
  int    pointCount      ;

  int     treeCount, polyI, i, j, bestI, bestJ ;
  double  bestDist, newDist ;

  //  Fail if either the startpoint or endpoint is outside the polygon set.
  if (!pointInPolygonSet(sX,sY,allPolys)
  ||  !pointInPolygonSet(eX,eY,allPolys)) {
    return NO; }

  //  If there is a straight-line solution, return with it immediately.
  if (lineInPolygonSet(sX,sY,eX,eY,allPolys)) {
    (*solutionNodes)=0; return YES; }

  //  Build a point list that refers to the corners of the
  //  polygons, as well as to the startpoint and endpoint.
  pointList[0].x=sX;
  pointList[0].y=sY; pointCount=1;
  for (polyI=0; polyI<allPolys.count; polyI++) {
    for (i=0; i<allPolys.poly[polyI].corners; i++) {
      pointList[pointCount].x=allPolys.poly[polyI].x[i];
      pointList[pointCount].y=allPolys.poly[polyI].y[i]; pointCount++; }}
  pointList[pointCount].x=eX;
  pointList[pointCount].y=eY; pointCount++;

  //  Initialize the shortest-path tree to include just the startpoint.
  treeCount=1; pointList[0].totalDist=0.;

  //  Iteratively grow the shortest-path tree until it reaches the endpoint
  //  -- or until it becomes unable to grow, in which case exit with failure.
  bestJ=0;


  while (bestJ<pointcount-1) {
    bestDist=INF;
    for (i=0; i<treeCount; i++) {
      for (j=treeCount; j<pointCount; j++) {
        if (lineInPolygonSet(
        pointList[i].x,pointList[i].y,
        pointList[j].x,pointList[j].y,allPolys)) {
          newDist=pointList[i].totalDist+calcDist(
          pointList[i].x,pointList[i].y,
          pointList[j].x,pointList[j].y);
          if (newDist<bestDist) {
            bestDist=newDist; bestI=i; bestJ=j; }}}}
    if (bestDist==INF) return NO;   //  (no solution)
    pointList[bestJ].prev     =bestI   ;
    pointList[bestJ].totalDist=bestDist;
    swapPoints(&pointList[bestJ],&pointList[treeCount]); treeCount++; }

  //  Load the solution arrays.
  (*solutionNodes)= -1; i=treeCount-1;
  while (i> 0) {
    i=pointList[i].prev; (*solutionNodes)++; }
  j=(*solutionNodes)-1; i=treeCount-1;
  while (j>=0) {
    i=pointList[i].prev;
    solutionX[j]=pointList[i].x;
    solutionY[j]=pointList[i].y; j--; }

  //  Success.
  return YES; }




 */
