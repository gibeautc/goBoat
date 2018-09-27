package vehical

import (
	"time"
	"fmt"
)

type SensingUnit struct{
	//have to be careful what we do here, we will be in another thread
	app *App
	shouldRun bool
}


func NewSensingUnit(app *App) *SensingUnit{
	n:=new(SensingUnit)
	n.app=app
	n.shouldRun=true
	return n
}

func (self *SensingUnit) Run(){
	go self.runThread()
}


func (self *SensingUnit) runThread(){
	fmt.Println("Starting Sensor Thread")
	for self.shouldRun{
		time.Sleep(10*time.Second)
		var s SensorData
		s.angles=make([]int,0)
		s.distances=make([]int,0)
		s.angles=append(s.angles,40)
		s.angles=append(s.angles,50)
		s.distances=append(s.distances,10)
		s.distances=append(s.distances,40)
		self.app.QueMsg(s)

	}
	fmt.Println("Exiting Sensor Thread")
}