package main

import (
	"fmt"
	"github.com/gibeautc/goBoat/boat"
)



func main() {
	app := new(boat.App)
	app.Events=make(chan boat.Msg,100)
	app.Idle=true
	app.Conn = boat.ConnectToDB("main.db")
	app.OsmMap = new(boat.MapData)
	app.LocalMap = new(boat.TileSet)
	app.LocalMap.Init()
	app.AllPolly = new(boat.PolySet)
	app.Sensing=boat.NewSensingUnit(app)
	app.Sensing.Run()
	app.HaveRoute=false


	app.CurLocation = new(boat.Point)
	app.CurLocation.Lat = 44.67618
	app.CurLocation.Lon = -123.09918

	app.Destination = new(boat.Point)
	app.Destination.Lat = 44.6378
	app.Destination.Lon = -123.1445

	//app.QueMsg(boat.LoadMapData{})
	app.QueMsg(boat.LoadCurrentTile{})
	//app.QueMsg(boat.FindRoute{})

	app.AddTimer(10000,boat.DoOneTimeTask{},false)
	app.AddTimer(20000,boat.SaveActiveToDisk{},true)

	var event boat.Msg
	for{
		event= app.WaitForEvent()
		switch event.(type){
		case boat.TimeOut:
			fmt.Println("Waiting for Event")
		default:
			if event.IsIdle() && !app.Idle{
				fmt.Println("Not Idle, Event Ignored")
				continue
				//not sure if we want to queue these up somewhere else or just plan to ignore them all together
			}
			event.Handle(app)
		}


	}



}

