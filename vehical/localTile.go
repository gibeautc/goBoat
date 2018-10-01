package vehical

import (
	"database/sql"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	tileSize        = 4096 //has to be even multiples of 2
	activeTileLimit = 5    //probably can be higher, but for testing we will keep it low
	maxCompression  = 1
	maxDiskSpace    = 200 //in MB
)

type Bounds struct {
	SW  Point
	SE  Point
	NW  Point
	NE  Point
	IdN int
	IdS int
	IdW int
	IdE int
}

type TileSet struct {
	activeTiles []Tile
	conn        *sql.DB
}

func (self *TileSet) Init() error {
	self.conn = ConnectToDB(folder + "database/tileSet.db")
	self.ClearTileCache()
	self.DumpDbAndCreateGenisisBlock(true)
	err := self.dbInit()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func GetDiskSpaceOfPathMB(path string) float32 {
	/*
		Returns the size of folder or file @path in MB
	*/
	out, err := exec.Command("du", "-hs", path).Output()
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	//fmt.Printf("Memory Used is %s\n", out)
	elems := strings.Split(string(out), "\t")
	//fmt.Println(elems)
	//fmt.Println("Num of elements: ",len(elems))
	var pref float32
	var valueString string
	wholeString := elems[0]
	//fmt.Println("ValueString: ",wholeString)
	//fmt.Println("Length: ",len(wholeString))
	wholeString = strings.Replace(wholeString, " ", "", -1)
	if strings.HasSuffix(wholeString, "B") {
		//fmt.Println("Size in Bytes")
		pref = 1. / (1000. * 1000.)
		valueString = strings.Replace(elems[0], "B", "", 1)
	} else if strings.HasSuffix(wholeString, "K") {
		//fmt.Println("Size in KB")
		pref = 1. / 1000.
		valueString = strings.Replace(elems[0], "K", "", 1)
	} else if strings.HasSuffix(wholeString, "M") {
		//fmt.Println("Size in MB")
		pref = 1.0
		valueString = strings.Replace(elems[0], "M", "", 1)
	} else if strings.HasSuffix(wholeString, "G") {
		//fmt.Println("Size in GB")
		pref = 1000.
		valueString = strings.Replace(elems[0], "G", "", 1)
	} else {
		fmt.Println("no suffix found....")
		return -100 //shouldnt happen once I actually write this function correctly.......
	}
	//fmt.Println(valueString)
	value, err := strconv.ParseFloat(valueString, 32)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	//fmt.Printf("Value: %f \tPrefix: %f\n",value,pref)
	return float32(value) * pref
}

func (self *TileSet) SaveAllActiveToDisk() {
	fmt.Println("Saving All Active Files to Disk")
	//todo make a copy of all tiles, and then save those to disk. this way it can be done in a go routine
	for x := 0; x < len(self.activeTiles); x++ {
		self.Pickle(self.activeTiles[x])
	}
}

func (self *TileSet) CheckMemoryAndCompress() error {
	/*
		Check the amount of spaced used by /tiles    and maybe /tileImages even thought that probably wont exist out in the wild
		if space used is more then maxDiskSpace then start compressing tiles untill it is below that threshold
	*/

	//du -hs tiles/
	used := GetDiskSpaceOfPathMB(folder + "tileImage/")
	for used > maxDiskSpace {
		fmt.Println("Need to compress tiles!!!!")
		id, err := self.GetOldestToCompress()
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if id == 0 {
			fmt.Println("Returned ID of zero.....")
			break
		}
		fmt.Println("Going to compress tile with ID: ", id)
		t, err := self.UnPickle(uint32(id))
		if err != nil {
			return err
		}
		if !t.Compress() {
			fmt.Println("could not compres the tile we were given")
			return errors.New("could not compres the tile we were given")
		}
		err = self.Pickle(*t)
		if err != nil {
			fmt.Println("Pickle Failed!!!!")
			return err
		}
		self.conn.Exec("UPDATE tiles set comp=$1 where id=$2", t.Size, t.Id)
		used = GetDiskSpaceOfPathMB(folder + "tileImage/")

	}
	fmt.Println("No Need to Compress at this time")
	return nil
}

func (self *TileSet) LoadTileById(id uint32) (int, error) {
	/*
		given an id of a tiles, return the index in activeTiles where it resides

		like loading tile for a point, it could re-oder/replace what is in activeTiles so any previous index should not be trusted
	*/

	//check if its in activeTiles first
	for x := 0; x < len(self.activeTiles); x++ {
		if self.activeTiles[x].Id == id {
			return x, nil
		}
	}

	//then check on disk
	files, err := filepath.Glob(folder + "tileImage/*")
	if err != nil {
		return -1, err
	}
	for x := 0; x < len(files); x++ {
		str := files[x]
		str = strings.Replace(str, folder+"tilesImages/", "", 1)
		fid, err := strconv.Atoi(str)
		if err != nil {
			fmt.Println("Bad File Name...ignoring: ", str)
			continue
		}
		if id == uint32(fid) {

			t, err := self.UnPickle(id)
			if err != nil {
				return 0, err
			}
			t.FullyExpand()
			//place tile in activeTiles
			return self.AddTile(*t), nil

		}

	}

	//since we are looking for one by id, if we have not found it, return error
	return 0, nil
}

func (self *TileSet) ClearTileCache() error {

	err := os.RemoveAll(folder + "tileImage/")
	if err != nil {
		return err
	}
	err = os.Mkdir(folder+"tileImage/", 0777)
	return err

}

func (self *TileSet) LoadTileForPoint(p Point) (int, error) {
	/*
		Given a point, return the index in activeTiles where the tile is that contains the point
		Will search activeTiles first, then tiles on disk. If one is still not found, then create a new tile, put it
		in activeTiles and return index

		Since this action could re-order activeTiles, any previously requested indexes are no longer valid
	*/

	//todo does not seem to be working!  might have been because of path issue, need to test again now
	for x := 0; x < len(self.activeTiles); x++ {
		if self.isPointInTile(p, self.activeTiles[x].Id) {
			return x, nil
		}
	}

	//now check tiles on disk
	files, err := filepath.Glob(folder + "tileImage/*")
	if err != nil {
		return 0, err
	}
	for x := 0; x < len(files); x++ {
		id, err := self.GetIdByPoint(p)
		if err != nil {
			return 0, err
		}
		t, err := self.UnPickle(uint32(id))
		if err != nil {
			return 0, err
		}
		t.FullyExpand()
		return self.AddTile(*t), nil

	}
	//if we get to here, we dont have it on disk, or in memory, so use new tile, but get new ID for it
	t := NewTile()
	t.Id, err = self.GetNewTileID()
	//todo need to create tile bounds, .....not trival
	fmt.Println("had to create new tile for point, but it has not bounds.......")
	if err != nil {
		return 0, err
	}
	return self.AddTile(*t), nil

}

//Adds tile to TileSet, and returns the index in activeTiles it lives in
func (self *TileSet) AddTile(t Tile) int {
	if len(self.activeTiles) < activeTileLimit {
		self.activeTiles = append(self.activeTiles, t)
		return len(self.activeTiles) - 1
	}
	//we need to move somthing to disk....for now it will just be the last one
	//todo come up with a better way to remove one.  Either by time since we have used it, or by distance
	self.Pickle(self.activeTiles[activeTileLimit-1])
	self.activeTiles[activeTileLimit-1] = t
	return activeTileLimit - 1
}

func (self *TileSet) AddDistanceDataSet(curLocation Point, sensorLocation Point) error {
	/*
		Given our current location, and a location of a "hit" update tile(s) data accordingly
		increase pixel at location by some value IncreasePercentage
		decrease all pixels between curLocation and sensorLocation by some value DecreasePercentage

		Should always be doing this on fully expanded tiles, but should try and make it work either way
	*/

	//are both points on same tile?
	sid, err := self.GetIdByPoint(curLocation)
	if err != nil {
		return err
	}
	eid, err := self.GetIdByPoint(sensorLocation)
	if err != nil {
		return err
	}
	if sid == eid {
		//both on same tile!
		index, err := self.LoadTileById(uint32(sid))
		if err != nil {
			return err
		}

		return self.activeTiles[index].AddDistanceData(curLocation, sensorLocation)
	} else {
		//on seperate tiles, will need to find boundary point that is included on both tiles

	}
	return nil
}

func (self *TileSet) Pickle(t Tile) error {
	fmt.Println("Starting Pickle Process")
	st := time.Now()
	err := self.updateTileToDB(t, -1)
	if err != nil {
		return err
	}
	err = t.SaveImage()
	if err != nil {
		return err
	}

	fmt.Println("Pickling took: ", time.Since(st))
	return nil
}

func (self *TileSet) UnPickle(id uint32) (*Tile, error) {
	var err error
	t := NewTile()
	t.Id = id
	t.Bounds,err = self.GetBounds(id)
	if err!=nil{
		return t,err
	}
	t.Img, err = LoadImage(int(id))
	if err != nil {
		return t, err
	}
	return t, nil
}

type Tile struct {
	//Data [][]byte
	Size   int
	Id     uint32
	Bounds Bounds
	Img    *image.Gray
}

func NewTile() *Tile {
	t := new(Tile)
	t.Size = tileSize
	r := image.Rect(0, 0, t.Size, t.Size)
	t.Img = image.NewGray(r)
	for x := 0; x < t.Size; x++ {
		for y := 0; y < t.Size; y++ {
			c := color.Gray{}
			c.Y = 128
			t.Img.Set(x, y, c)
		}

	}

	return t
}

func (self *TileSet) isPointInTile(p Point, id uint32) bool {
	//todo write this just look up the data in the database

	return false
}

func (self *Tile) Compress() bool {
	fmt.Println("Compressing tile: ", self.Id)
	fmt.Println("Starting Size: ", self.Size)

	if self.Size <= maxCompression {
		return false
	}

	self.Size = self.Size / 2
	r := image.Rect(0, 0, self.Size, self.Size)
	newImg := image.NewGray(r)

	var c color.Gray
	for x := 0; x < self.Size; x++ {
		for y := 0; y < self.Size; y++ {
			sum := int(self.Img.GrayAt(x*2, y*2).Y) + int(self.Img.GrayAt(x*2+1, y*2).Y) + int(self.Img.GrayAt(x*2, y*2+1).Y) + int(self.Img.GrayAt(x*2+1, y*2+1).Y)
			c.Y = uint8(sum / 4)
			newImg.Set(x, y, c)
		}
	}
	self.Img = newImg
	fmt.Println("New Size: ", self.Size)
	return true
}

func (self *Tile) FullyExpand() {
	for self.Expand() {
		fmt.Println("Expanded to: ", self.Size)
	}
}

func (self *Tile) Expand() bool {
	if self.Size >= tileSize {
		return false
	}

	self.Size = self.Size * 2
	r := image.Rect(0, 0, self.Size, self.Size)
	newImg := image.NewGray(r)
	var c color.Gray
	for x := 0; x < self.Size; x += 2 {
		for y := 0; y < self.Size; y += 2 {
			c.Y = self.Img.GrayAt(x/2, y/2).Y
			newImg.Set(x/2, y/2, c)
			newImg.Set(x/2+1, y/2, c)
			newImg.Set(x/2, y/2+1, c)
			newImg.Set(x/2+1, y/2+1, c)

		}

	}
	self.Img = newImg
	return true
}

func (self *Tile) SaveImage() error {

	return SaveImage(self.Img, folder+"tileImage/"+strconv.Itoa(int(self.Id))+".png")
}

func (self *Tile) GetPixelByCords(p Point) (int, int) {

	percentX := (p.Lon - self.Bounds.SW.Lon) / (self.Bounds.SE.Lon - self.Bounds.SW.Lon)
	percentY := (p.Lat - self.Bounds.SW.Lat) / (self.Bounds.NW.Lat - self.Bounds.SW.Lat)
	return int(float64(self.Size) * percentX), int(float64(self.Size) * percentY)

}

func (self *Tile) AddDistanceData(curLocation Point, objectLocation Point) error {
	//we know both points are on the tile already, just need to update pixel info

	var ls, le image.Point
	ls.X, ls.Y = self.GetPixelByCords(curLocation)
	le.X, le.Y = self.GetPixelByCords(objectLocation)
	drawLine(self.Img, ls, le, color.White)
	return nil
}
