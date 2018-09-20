package main

import (
	"github.com/gibeautc/goBoat/boat"
	"fmt"
)

func main(){
	//myMap:=new(boat.MapData)
	//st:=time.Now()
	//myMap.Load("mapLarge")
	//fmt.Println("Load Time: ",time.Since(st))
	//fmt.Println("Number of Nodes: ",len(myMap.Data.Nodes))
	//fmt.Println("Number of Ways: ",len(myMap.Data.Ways))
	//fmt.Println("Number of Relations: ",len(myMap.Data.Relations))
	//myMap.ParseForWater()

	allTests()
}




func allTests(){
	if squareTest(){
		fmt.Println("Square: PASS")
	} else{
		fmt.Println("Square: FAIL")
	}
}






func squareTest() bool{
	//sqare polygon
	var poly boat.Poly
	var polySet boat.PolySet
	poly.AddCorner(0,0)
	poly.AddCorner(0,10)
	poly.AddCorner(5,5)
	poly.AddCorner(10,10)
	poly.AddCorner(10,0)
	polySet.AddPoly(poly)
	polySet.Print()

	var start,end boat.Point
	x:=6.0
	start.X=1
	start.Y=x
	end.X=7
	end.Y=x

	r,err:=boat.ShortestPath(start,end,polySet)
	r.Print()
	boat.Draw(&polySet,&r,start,end)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}

	if r.Count!=0{
		fmt.Println("Route Count should be Zero")
		return false
	}
	return true
}



