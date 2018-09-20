package boat

import (
	"github.com/glaslos/go-osm"
	"io/ioutil"
	"fmt"
	"bytes"
	"github.com/hydrogen18/stalecucumber"
)

type MapData struct{
	Data *osm.Map
}



func (self *MapData) Load(filename string) error{
	var err error
	data, err := ioutil.ReadFile(filename)
	if err!=nil{
		return err
	}

	self.Data,err= osm.DecodeString(string(data))
	if err!=nil{
		return err
	}
	return nil
}


func (self *MapData) Save() error{
	buf := new(bytes.Buffer)
	_,err := stalecucumber.NewPickler(buf).Pickle(&self.Data)
	if err!=nil{
		return nil
	}

	err = ioutil.WriteFile("pickleMap", buf.Bytes(), 0644)
	return err
}


func (self *MapData) ParseForWater() error{
	//search way tags for water related tag, add to new list
	newWay:=make([]osm.Way,0)
	for x:=0;x<len(self.Data.Ways);x++{
		for y:=0;y<len(self.Data.Ways[x].RTags);y++{
			curTag:=self.Data.Ways[x].RTags[y]
			if curTag.Key=="waterway"{
				fmt.Println("Found Waterway:",self.Data.Ways[x].ID)
				newWay=append(newWay,self.Data.Ways[x])
				break
			}
		}
	}


	//go through new list and see what nodes we need to keep--> new list,
	newNodes:=make([]osm.Node,0)
	for x:=0;x<len(newWay);x++{
		wayPercent:=float64(x)/float64(len(newWay))
		fmt.Println(wayPercent)
		for y:=0;y<len(newWay[x].Nds);y++{
			id:=newWay[x].Nds[y].ID
			for i:=0;i<len(self.Data.Nodes);i++{
				if self.Data.Nodes[i].ID==id{
					fmt.Print("+")
					newNodes=append(newNodes,self.Data.Nodes[i])
					break
				}
				fmt.Println("")
			}
		}

	}
	fmt.Println("Original Number of Ways: ",len(self.Data.Ways))
	fmt.Println("Original Number of Nodes: ",len(self.Data.Nodes))
	self.Data.Ways=newWay
	self.Data.Nodes=newNodes
	fmt.Println("Current Number of Ways: ",len(self.Data.Ways))
	fmt.Println("Current Number of Nodes: ",len(self.Data.Nodes))
	return nil
}

