package vehical

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"time"
)

const folder = "/home/chadg/go/src/github.com/gibeautc/goBoat/"

/*
latLst,lonLst are decimal degrees  (float64)
distance is measured in inches (int64)
angles are measred in degrees (int)
speed in mph
*/


func DistanceBetween(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (int64, int) {
	R := 6371000.0 //m
	latDelta := toRadians(lat2 - lat1)
	lonDelta := toRadians(lon2 - lon1)
	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)
	a := math.Sin(latDelta/2)*math.Sin(latDelta/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(lonDelta/2)*math.Sin(lonDelta/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c      // in m
	d = 39.3701 * d //inches

	y := math.Sin(lonDelta) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(lonDelta)
	brng := math.Atan2(y, x)

	//return distance, bearing
	return int64(d), int(toDegrees(brng))
}

func toRadians(deg float64) float64 {
	return deg / 360 * 2 * math.Pi
}

func toDegrees(rad float64) float64 {
	return rad / (2 * math.Pi) * 360
}

func GetCords(lat float64, lon float64, distance int64, direction int) (float64, float64) {
	R := 6378.1 //km
	angle := toRadians(float64(direction))
	d := float64(distance) * .0000254
	startLat := toRadians(lat)
	startLon := toRadians(lon)
	lat2 := math.Asin(math.Sin(startLat)*math.Cos(d/R) + math.Cos(startLat)*math.Sin(d/R)*math.Cos(angle))
	lon2 := startLon + math.Atan2(math.Sin(angle)*math.Sin(d/R)*math.Cos(startLat), math.Cos(d/R)-math.Sin(startLat)*math.Sin(lat2))
	return toDegrees(lat2), toDegrees(lon2)
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

type App struct {
	Conn        *sql.DB
	OsmMap      *MapData
	LocalMap    *TileSet
	CurLocation *Point
	LocalIndex  int
	Destination *Point
	Route       Route
	AllPolly    *PolySet
	Events      chan Msg
	Idle        bool
	Sensing     *SensingUnit
	HaveRoute   bool //set to true when we are in the process of finding a route
}

func (app *App) QueMsg(msg Msg) {
	app.Events <- msg
}

func (app *App) WaitForEvent() Msg {
	if len(app.Events) == 0 {
		//only start the timer if we dont already have something to process
		app.AddTimer(1000, TimeOut{}, false)
	}
	for {
		ev := <-app.Events
		if !app.Idle && ev.IsIdle() {
			//we want to delay the event
			fmt.Println("Delaying Event: ", ev)
			app.AddTimer(500, ev, false)
			continue
		}
		return ev
	}

}

func (app *App) AddTimer(interval int, msg Msg, repeating bool) {
	if repeating {
		go app.repeatingTimer(msg, interval)
	} else {
		go app.nonRepeatingTimer(msg, interval)
	}
}

func (app *App) repeatingTimer(msg Msg, interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		app.Events <- msg
	}
}

func (app *App) nonRepeatingTimer(msg Msg, interval int) {
	time.Sleep(time.Duration(interval) * time.Millisecond)
	app.Events <- msg
}

func (app *App) DoRoute() {
	//careful! we are in a go routine!

	var err error
	app.Route, err = ShortestPath(*app.CurLocation, *app.Destination, *app.AllPolly)
	if err != nil {
		app.HaveRoute = false
		fmt.Println("Routing Error!")
		fmt.Println(err.Error())
		DrawWorld(app.AllPolly)
	} else {
		app.HaveRoute = true
	}
	app.Route.Print()

	if err == nil {
		Draw(app.AllPolly, &app.Route, *app.CurLocation, *app.Destination)
	}
}

func (app *App) PrintState() {
	fmt.Println("**********Current State**********")
	fmt.Println("Length of activeTiles: ", len(app.LocalMap.activeTiles))
	fmt.Println("")
}

func (app *App) Init() {
	st := time.Now()
	app.Events = make(chan Msg, 100)
	app.Idle = true
	app.Conn = ConnectToDB("database/main.db")
	app.OsmMap = new(MapData)
	app.LocalMap = new(TileSet)
	app.LocalMap.Init()
	app.AllPolly = new(PolySet)
	app.Sensing = NewSensingUnit(app)
	app.Sensing.Run()
	app.HaveRoute = false
	//44.616028, -123.073269
	app.CurLocation = new(Point)
	app.CurLocation.Lat = 44.616028
	app.CurLocation.Lon = -123.073269

	app.Destination = new(Point)
	app.Destination.Lat = 44.6378
	app.Destination.Lon = -123.1445
	fmt.Println("App initialized time: ", time.Since(st))
}
