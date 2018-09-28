package main

import (
	"fmt"
	"github.com/gibeautc/goBoat/vehical"
)

func main() {
	app := new(vehical.App)
	app.Init()

	//app.QueMsg(vehical.LoadMapData{})
	app.QueMsg(vehical.LoadCurrentTile{})
	//app.QueMsg(vehical.FindRoute{})

	app.AddTimer(10000, vehical.DoOneTimeTask{}, false)
	app.AddTimer(20000, vehical.SaveActiveToDisk{}, true)
	app.AddTimer(60000, vehical.CheckMemoryCompress{}, true)
	var event vehical.Msg
	for {
		event = app.WaitForEvent()
		switch event.(type) {
		case vehical.TimeOut:
			fmt.Println("Waiting for Event")
		default:
			if event.IsIdle() && !app.Idle {
				fmt.Println("Not Idle, Event Ignored")
				continue
				//not sure if we want to queue these up somewhere else or just plan to ignore them all together
			}
			event.Handle(app)
		}

	}

}
