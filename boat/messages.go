package boat

import (
	"fmt"
	"time"
)

type Msg interface {
	Handle(app *App)
	IsIdle() bool
}


type LoadCurrentTile struct{}
func (self LoadCurrentTile) IsIdle() bool{return true}
func (self LoadCurrentTile) Handle(app *App){
	var err error
	app.LocalIndex, err = app.LocalMap.LoadTileForPoint(*app.CurLocation)
	fmt.Println("We are currently in tile with index:", app.LocalIndex)

	if err != nil {
		fmt.Println("Failed to get current local tile....")
		fmt.Println(err.Error())
		return
	}
}


type FindRoute struct{}
func (self FindRoute) IsIdle() bool{return true}
func (self FindRoute) Handle(app *App){
	fmt.Println("Spawning Routing Thread")
	go app.DoRoute()
}


type DoOneTimeTask struct{}
func (self DoOneTimeTask) IsIdle() bool {return true}
func (self DoOneTimeTask) Handle(app *App) {
	fmt.Println("Doing a one time task!")
}


type TimeOut struct{}
func (self TimeOut) IsIdle() bool {return false}
func (self TimeOut) Handle(app *App){
	fmt.Println("Default Timeout Handler")
}




type LoadMapData struct{}
func (self LoadMapData) IsIdle() bool {return true}
func (self LoadMapData) Handle(app *App){
	st := time.Now()
	err:= app.OsmMap.Load("largeMap.osm")
	if err != nil {
		fmt.Println("Failed to load map")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Load Time: ", time.Since(st))
		fmt.Println("Number of Nodes: ", len(app.OsmMap.Data.Nodes))
		fmt.Println("Number of Ways: ", len(app.OsmMap.Data.Ways))
		fmt.Println("Number of Relations: ", len(app.OsmMap.Data.Relations))
		app.OsmMap.ParseForWater()

		for i := 0; i < len(app.OsmMap.Data.Ways); i++ {
			p := new(Poly)
			for x := 0; x < len(app.OsmMap.Data.Ways[i].Nds); x++ {
				nid := app.OsmMap.Data.Ways[i].Nds[x].ID
				for y := 0; y < len(app.OsmMap.Data.Nodes); y++ {
					if nid == app.OsmMap.Data.Nodes[y].ID {
						n := app.OsmMap.Data.Nodes[y]
						p.AddCorner(n.Lng, n.Lat)
						break
					}
				}

			}
			if p.Verify() {
				app.AllPolly.AddPoly(*p)
			}

		}
	}
}



type SaveActiveToDisk struct{}
func (self SaveActiveToDisk) IsIdle() bool {return true}
func (self SaveActiveToDisk) Handle(app *App){
	app.LocalMap.SaveAllActiveToDisk()
}


type SensorData struct{
	angles []int
	distances []int

}
func (self SensorData) IsIdle() bool {return false}
func (self SensorData) Handle(app *App){
	fmt.Println("Got Sensor Data....what to do with it?")
	fmt.Println("Number of points: ",len(self.angles))
}