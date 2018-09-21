package main

import "github.com/gibeautc/goBoat/boat"

func main(){
	//myMap:=new(boat.MapData)
	//st:=time.Now()
	//myMap.Load("mapLarge")
	//fmt.Println("Load Time: ",time.Since(st))
	//fmt.Println("Number of Nodes: ",len(myMap.Data.Nodes))
	//fmt.Println("Number of Ways: ",len(myMap.Data.Ways))
	//fmt.Println("Number of Relations: ",len(myMap.Data.Relations))
	//myMap.ParseForWater()

	//allTests()

	//local:=new(boat.TileSet)
	//id,err:=local.GetNewID()
	//if err!=nil{
	//	fmt.Println("Error getting new ID: ",err.Error())
	//}
	//fmt.Println("New Tile ID: ",id)
	//tile:=boat.NewTile()
	//tile.Id=id
	//tile.Pickle()


	//t:=boat.NewTile()
	//t.Id=9529
	//t.UnPickle()
	//fmt.Println("my ID is :",t.Id)

	//boat.CompressionTest()

	local:=new(boat.TileSet)
	local.CheckMemoryAndCompress()
}



