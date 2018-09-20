package boat

import (
	"io/ioutil"
	"strconv"
	"strings"
	"path/filepath"
	"fmt"
	"math/rand"
	"time"
	"bytes"
	"github.com/hydrogen18/stalecucumber"
	"os"
	"image"
	"image/color"
)

/*
Tiles are square and at 1000 in, they have 83.333 ft sides
byte in each little square indicates the liklyhood of it being blocked or impassible.
255-> def cant go there
0-> we have been there so we know its clear

When a Tile is shifted off as being the current one, we can compress the image as needed, further away-> more compression

tiles are always oriented North-> up
if we have Tile Data adjsant to that Tile, the Tile Id will be listed in the struct. 0 means we have no Data
 */


 /*
 Sizes and time values are on ubuntu nuk
 tile starts as 10x10 and expands to the max
1280	7.9M		5.091896117s
640		2.0M		1.384764334s
320		502K		417.147612ms
160		126K		113.71121ms
80		32K			24.222966ms
40		8.2K		4.543498ms
20		2.3K		1.897164ms
10		772			826.44µs

 */


const tileSize=1280  //has to be even multiples of 2
const activeTileLimit=5  //probably can be higher, but for testing we will keep it low
const maxCompression=10



type TileSet struct{
	activeTiles []Tile

}



/*
Given a point, return the index in activeTiles where the tile is that contains the point
Will search activeTiles first, then tiles on disk. If one is still not found, then create a new tile, put it
in activeTiles and return index


Since this action could re-order activeTiles, any previously requested indexes are no longer valid

 */
func (self *TileSet) LoadTileForPoint(p Point) (int,error){
	for x:=0;x<len(self.activeTiles);x++{
		if self.activeTiles[x].isPointInTile(p){
			return x,nil
		}
	}

	//now check tiles on disk
	files, err := filepath.Glob("tiles/*")
	if err != nil {
		return 0,err
	}
	t:=NewTile()
	for x:=0;x<len(files);x++{
		str:=files[x]
		str=strings.Replace(str,"tiles/","",1)
		id,err:=strconv.Atoi(str)
		if err!=nil{
			fmt.Println("Bad File Name...ignoring: ",str)
			continue
		}

		t.Id=uint32(id)
		t.UnPickle()
		if t.isPointInTile(p){
			//found one that works, now add it to self
			return self.AddTile(*t),nil
		}
	}
	//if we get to here, we dont have it on disk, or in memory, so use new tile, but get new ID for it
	t.Id,err=self.GetNewID()
	if err!=nil{
		return 0,err
	}
	return self.AddTile(*t),nil

}


//Adds tile to TileSet, and returns the index in activeTiles it lives in
func (self *TileSet) AddTile(t Tile) int{
	if len(self.activeTiles)<activeTileLimit{
		self.activeTiles=append(self.activeTiles,t)
		return len(self.activeTiles)-1
	}
	//we need to move somthing to disk....for now it will just be the last one
	//todo come up with a better way to remove one.  Either by time since we have used it, or by distance
	self.activeTiles[activeTileLimit-1].Pickle()
	self.activeTiles[activeTileLimit-1]=t
	return activeTileLimit-1
}

func (self *TileSet) GetNewID() (uint32,error){
	used:=make([]uint32,0)
	for x:=0;x<len(self.activeTiles);x++{
		used=append(used,self.activeTiles[x].Id)
	}

	files, err := filepath.Glob("tiles/*")
	if err != nil {
		return 0,err
	}
	for x:=0;x<len(files);x++{
		str:=files[x]
		str=strings.Replace(str,"tiles/","",1)
		id,err:=strconv.Atoi(str)
		if err!=nil{
			fmt.Println("Bad File Name...ignoring: ",str)
			continue
		}
		used=append(used,uint32(id))
	}
	var id uint32
	for id<=0{
		rand.Seed(time.Now().Unix())
		tryNum:=rand.Uint32()/1000000
		if tryNum>10000{
			fmt.Println("too large")
			continue
		}
		if !sliceContainsUint32(used,tryNum){
			id=tryNum
			break
		}
	}

	return id,nil
}

type Tile struct{
	Data [][]byte
	Id   uint32
	IdN  uint32
	IdS  uint32
	IdW  uint32
	IdE  uint32
	NW   Point //points for the four corners of the tile, can also be used to make a square to see if a point is in a given tile
	SW   Point
	NE   Point
	SE   Point
}



func NewTile() *Tile {
	t:=new(Tile)
	for x:=0;x<tileSize;x++{
		d:=make([]byte,tileSize)
		for y:=0;y<tileSize;y++{
			d[y]=128  //new Tile will be filled with middle numbers (unknown)
		}
		t.Data =append(t.Data,d)
	}



	return t
}

func (self *Tile) Pickle() error{
	buf := new(bytes.Buffer)
	_,err := stalecucumber.NewPickler(buf).Pickle(&self)
	if err!=nil{
		return nil
	}

	err = ioutil.WriteFile("tiles/"+strconv.Itoa(int(self.Id)), buf.Bytes(), 0644)
	return err
}

/*
File saved as:		2416712524
loading back as: 	2416712523
 */


func (self *Tile) UnPickle() error{
	var err error
	//data, err := ioutil.ReadFile("tiles/"+strconv.Itoa(int(self.Id)))

	f, err := os.Open("tiles/"+strconv.Itoa(int(self.Id)))
	if err!=nil{
		return err
	}

	err = stalecucumber.UnpackInto(&self).From(stalecucumber.Unpickle(f))
	if err!=nil{
		return err
	}
	return nil
}


func(self *Tile) isPointInTile(p Point) bool{
	if p.X<self.NE.X && p.X>self.NW.X{
		if p.Y>self.SW.Y && p.Y<self.NW.Y{
			return true
		}
	}

	return false
}

func (self *Tile) Compress() bool{
	if len(self.Data)<=maxCompression{
		return false
	}
	newData:=make([][]byte,0)
	newSize:=len(self.Data)/2
	for x:=0;x<newSize;x++{
		newRow:=make([]byte,0)
		for y:=0;y<newSize;y++{
			sum:=int(self.Data[x*2][y*2])+int(self.Data[x*2+1][y*2])+int(self.Data[x*2][y*2+1])+int(self.Data[x*2+1][y*2+1])
			newRow=append(newRow,byte(sum/4))
		}
		newData=append(newData,newRow)
	}
	self.Data=newData
	return true
}


func (self *Tile) Expand() bool{
	if len(self.Data)>=tileSize{
		return false
	}
	newData:=make([][]byte,0)
	newSize:=len(self.Data)*2
	for x:=0;x<newSize;x+=2{
		newRow:=make([]byte,0)
		for y:=0;y<newSize;y+=2{
			value:=self.Data[x/2][y/2]
			newRow=append(newRow,value)
			newRow=append(newRow,value)

		}
		newData=append(newData,newRow)
		newData=append(newData,newRow)
	}
	self.Data=newData
	return true
}

func (self *Tile) PrintData(){
	for x:=0;x<len(self.Data);x++{
		fmt.Println(self.Data[x])
	}

}


func (self *Tile) SaveImage() error{
	imgSize:=len(self.Data)
	r:=image.Rect(0,0,imgSize,imgSize)
	img:=image.NewAlpha(r)
	for x:=0;x<len(self.Data);x++{
		for y:=0;y<len(self.Data);y++{
			c:=color.Opaque
			c.A=uint16(self.Data[x][y])*256
			img.Set(x,y,c)
		}

	}
	return SaveImage(img,"tileImage/"+strconv.Itoa(int(self.Id))+".jpg")
}