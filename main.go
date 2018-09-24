package main

import (
	"database/sql"
	"fmt"
	"github.com/gibeautc/goBoat/boat"
	"time"
)

type App struct {
	conn        *sql.DB
	osmMap      *boat.MapData
	localMap    *boat.TileSet
	curLocation *boat.Point
	localIndex  int
	destination *boat.Point
	route       boat.Route
	allPoly     *boat.PolySet
}

func main() {
	var err error
	app := new(App)
	app.conn = boat.ConnectToDB()
	app.osmMap = new(boat.MapData)
	app.localMap = new(boat.TileSet)
	app.localMap.Init()
	app.allPoly = new(boat.PolySet)
	st := time.Now()
	err = app.osmMap.Load("largeMap.osm")
	if err != nil {
		fmt.Println("Failed to load map")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Load Time: ", time.Since(st))
		fmt.Println("Number of Nodes: ", len(app.osmMap.Data.Nodes))
		fmt.Println("Number of Ways: ", len(app.osmMap.Data.Ways))
		fmt.Println("Number of Relations: ", len(app.osmMap.Data.Relations))
		app.osmMap.ParseForWater()

		for i := 0; i < len(app.osmMap.Data.Ways); i++ {
			p := new(boat.Poly)
			for x := 0; x < len(app.osmMap.Data.Ways[i].Nds); x++ {
				nid := app.osmMap.Data.Ways[i].Nds[x].ID
				for y := 0; y < len(app.osmMap.Data.Nodes); y++ {
					if nid == app.osmMap.Data.Nodes[y].ID {
						n := app.osmMap.Data.Nodes[y]
						p.AddCorner(n.Lng, n.Lat)
						break
					}
				}

			}
			if p.Verify() {
				app.allPoly.AddPoly(*p)
			}

		}
	}
	app.curLocation = new(boat.Point)
	app.curLocation.Lat = 44.67618
	app.curLocation.Lon = -123.09918

	app.destination = new(boat.Point)
	app.destination.Lat = 44.6378
	app.destination.Lon = -123.1445
	st = time.Now()
	app.localIndex, err = app.localMap.LoadTileForPoint(*app.curLocation)
	fmt.Println("Load Time: ", time.Since(st))
	fmt.Println("We are currently in tile with index:", app.localIndex)
	if err != nil {
		fmt.Println("Failed to get current local tile....exiting")
		fmt.Println(err.Error())
		return
	}


	//keyPoints:= app.localMap.GetPolygons(app.localIndex)
	//fmt.Println(keyPoints)
	//return
	app.route, err = boat.ShortestPath(*app.curLocation, *app.destination, *app.allPoly)
	if err != nil {
		fmt.Println("Routing Error!")
		fmt.Println(err.Error())
		boat.DrawWorld(app.allPoly)
	}
	app.route.Print()

	if err == nil {
		boat.Draw(app.allPoly, &app.route, *app.curLocation, *app.destination)
	}

	fmt.Println("nothing else to do for now")
}
