package vehical

import (
	"io/ioutil"
	"strconv"
	"strings"
	"path/filepath"
	"fmt"
	"bytes"
	"github.com/hydrogen18/stalecucumber"
	"image"
	"image/color"
	"os"
	"database/sql"
	"os/exec"
	"time"
)



/*
Tiles are square and at 1280 in, they have ~106 ft sides
byte in each little square indicates the liklyhood of it being blocked or impassible.
255-> def cant go there
0-> we have been there so we know its clear

When a Tile is shifted off as being the current one, we can compress the image as needed, further away-> more compression

tiles are always oriented North-> up
if we have Tile Data adjsant to that Tile, the Tile Id will be listed in the struct. 0 means we have no Data
*/

const (


	tileSize=4096  //has to be even multiples of 2
 	activeTileLimit=5  //probably can be higher, but for testing we will keep it low
	maxCompression=1
	maxDiskSpace= 20 //in MB
)



type TileSet struct{
	activeTiles []Tile
	conn *sql.DB

}

func (self *TileSet) Init(){
	self.conn=ConnectToDB(folder+"database/tileSet.db")
	err:=self.dbInit()
	if err!=nil{
		fmt.Println(err.Error())
	}
	self.updateTilesFromDiskToDB()
}


func GetDiskSpaceOfPath(path string) float32{
	out, err := exec.Command("du","-hs", path).Output()
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	fmt.Printf("Memory Used is %s\n", out)
	elems:=strings.Split(string(out)," ")
	fmt.Println(elems)
	var pref float32
	var valueString string
	wholeString:=elems[0]
	wholeString=strings.Replace(wholeString," ","",-1)
	if strings.HasSuffix(wholeString,"B"){
		pref=1/(1000*1000)
		valueString=strings.Replace(elems[0],"B","",1)
	}else if strings.HasSuffix(wholeString,"K"){
		pref=1/1000
		valueString=strings.Replace(elems[0],"K","",1)
	}else if strings.HasSuffix(wholeString,"M"){
		pref=1.0
		valueString=strings.Replace(elems[0],"M","",1)
	}else if strings.HasSuffix(wholeString,"G"){
		pref=1000
		valueString=strings.Replace(elems[0],"G","",1)
	}else{
		fmt.Println("no suffix found....")
		return 0.0
	}
	fmt.Println(valueString)
	value,err:=strconv.ParseFloat(valueString,32)
	if err!=nil{
		fmt.Println(err.Error())
		return 0.0
	}
	return float32(value)*pref
}

func (self *TileSet) SaveAllActiveToDisk(){
	fmt.Println("Saving All Active Files to Disk")
	//todo make a copy of all tiles, and then save those to disk. this way it can be done in a go routine
	for x:=0;x<len(self.activeTiles);x++{
		self.activeTiles[x].Pickle()
	}
}

func (self * TileSet) CheckMemoryAndCompress() {
	/*
	Check the amount of spaced used by /tiles    and maybe /tileImages even thought that probably wont exist out in the wild
	if space used is more then maxDiskSpace then start compressing tiles untill it is below that threshold
	 */

	 //du -hs tiles/
	used:=GetDiskSpaceOfPath("tiles/")
	if used>maxDiskSpace{
		fmt.Println("Need to compress tiles!!!!")
	}

}



func (self * TileSet) LoadTileById(id uint32) (int,error){
	/*
	given an id of a tiles, return the index in activeTiles where it resides

	like loading tile for a point, it could re-oder/replace what is in activeTiles so any previous index should not be trusted
	 */

	//check if its in activeTiles first


	//then check on disk



	//since we are looking for one by id, if we have not found it, return error
	return 0,nil
}




func (self *TileSet) LoadTileForPoint(p Point) (int,error){
	/*
	Given a point, return the index in activeTiles where the tile is that contains the point
	Will search activeTiles first, then tiles on disk. If one is still not found, then create a new tile, put it
	in activeTiles and return index

	Since this action could re-order activeTiles, any previously requested indexes are no longer valid



	 */

	 //todo does not seem to be working!
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
	t.Id,err=self.GetNewTileID()

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
	fmt.Println("Starting Pickle Process")
	st:=time.Now()
	buf := new(bytes.Buffer)
	_,err := stalecucumber.NewPickler(buf).Pickle(&self)
	if err!=nil{
		return nil
	}

	err = ioutil.WriteFile(folder+"tiles/"+strconv.Itoa(int(self.Id)), buf.Bytes(), 0644)
	fmt.Println("Pickling took: ",time.Since(st))
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
	if p.Lon <self.NE.Lon && p.Lon >self.NW.Lon {
		if p.Lat >self.SW.Lat && p.Lat <self.NW.Lat {
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
