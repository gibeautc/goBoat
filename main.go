package main

import (
	"fmt"
	"github.com/gibeautc/goBoat/vehical"
)



func main() {
	app := new(vehical.App)
	app.Events=make(chan vehical.Msg,100)
	app.Idle=true
	app.Conn = vehical.ConnectToDB("database/main.db")
	app.OsmMap = new(vehical.MapData)
	app.LocalMap = new(vehical.TileSet)
	app.LocalMap.Init()
	app.AllPolly = new(vehical.PolySet)
	app.Sensing= vehical.NewSensingUnit(app)
	app.Sensing.Run()
	app.HaveRoute=false

	app.CurLocation = new(vehical.Point)
	app.CurLocation.Lat = 44.67618
	app.CurLocation.Lon = -123.09918

	app.Destination = new(vehical.Point)
	app.Destination.Lat = 44.6378
	app.Destination.Lon = -123.1445

	//app.QueMsg(vehical.LoadMapData{})
	app.QueMsg(vehical.LoadCurrentTile{})
	//app.QueMsg(vehical.FindRoute{})

	app.AddTimer(10000, vehical.DoOneTimeTask{},false)
	app.AddTimer(20000, vehical.SaveActiveToDisk{},true)
	app.AddTimer(60000,vehical.CheckMemoryCompress{},true)
	var event vehical.Msg
	for{
		event= app.WaitForEvent()
		switch event.(type){
		case vehical.TimeOut:
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

