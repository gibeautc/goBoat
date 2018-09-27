package vehical

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
	"path/filepath"
	"strings"
	"fmt"
)



func ConnectToDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}
	return db
}



func (self *TileSet) dbInit() error{
	_,err:=self.conn.Exec("CREATE TABLE IF NOT EXISTS tiles(id INTEGER PRIMARY KEY AUTOINCREMENT ,onDisk INTEGER,comp INTEGER,inRam INTEGER,lastUsed INTEGER,N INTEGER,S INTEGER,E INTEGER,W INTEGER,NELat REAL,NWLat REAL,SELat REAL,SWLat REAL,NELon REAL,NWLon REAL,SELon REAL,SWLon REAL)")
	if err!=nil{
		return err
	}
	return nil
}


func (self *TileSet) GetNewTileID() (uint32,error){
	res,err:=self.conn.Exec("INSERT INTO tiles(lastUsed) VALUES($1)",time.Now().Unix())
	if err!=nil{
		return 0,err
	}
	id,err:=res.LastInsertId()
	if err!=nil{
		return 0,err
	}
	return uint32(id),nil

}

func(self *TileSet) updateTileToDB(tile Tile,index int) error{
	onDisk:=Exists(folder+"tiles/"+strconv.Itoa(int(tile.Id)))
	//var onDisk int
	//if onDiskBool{
	//	onDisk=1
	//}else{
	//	onDisk=0
	//}
	_,err:=self.conn.Exec("REPLACE INTO tiles(id,onDisk,comp,inRam,lastUsed,NWLat,NWLon,NELat,NELon,SELat,SELon,SWLat,SWLon,N,S,E,W) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)",
							tile.Id,onDisk,len(tile.Data),index,time.Now().Unix(),tile.NW.Lat,tile.NW.Lon,tile.NE.Lat,tile.NE.Lon,tile.SE.Lat,tile.SE.Lon,tile.SW.Lat,tile.SW.Lon,tile.IdN,tile.IdS,tile.IdE,tile.IdW)
	return err
}


func(self *TileSet) UpdateAllTilesInDB()error{
	for x:=0;x<len(self.activeTiles);x++{
		err:=self.updateTileToDB(self.activeTiles[x],x)
		if err!=nil{
			return err
		}
	}
	return nil
}


func(self *TileSet) updateTilesFromDiskToDB() error{
	//now check tiles on disk
	files, err := filepath.Glob("tiles/*")
	if err != nil {
		return err
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
		err=self.updateTileToDB(*t,-1)
		if err!=nil{
			return err
		}
	}
	return nil
}



func(self *TileSet) DumpDbAndCreateGenisisBlock(createGenisis bool) error{
	var err error
	self.conn.Exec("DROP TABLE tiles")
	self.dbInit()
	self.conn.Exec("VACUUM")
	if !createGenisis{
		return nil
	}
	t:=NewTile()
	t.Id,err=self.GetNewTileID()
	if err!=nil{
		return err
	}
	var NW,NE,SE,SW Point

	NW.Lat=44.616896281347714
	NW.Lon=-123.07328249650676
	
	NE.Lat=44.61689628886895
	NE.Lon=-123.07196952462509

	SE.Lat=44.61596169059874
	SE.Lon=-123.07196952462509

	SW.Lat=44.61596169811975
	SW.Lon=-123.07328247537501
	
	t.NW=NW
	t.NE=NE
	t.SE=SE
	t.SW=SW
	err=t.Pickle()
	if err!=nil{
		return err
	}

	return self.updateTileToDB(*t,-1)
}


func (self *TileSet) GetOldestToCompress() (int,error){
	//get the tile ID used the longest ago that can be compressed, ie size>1
	rows,err:=self.conn.Query("SELECT id,lastUsed from tiles where onDisk=1 and comp>1 ")
	if err!=nil{
		return 0,err
	}
	oldest:=time.Now().Unix()+100
	var bestId int
	id:=0
	var lastUsed int64
	for rows.Next(){
		err=rows.Scan(&id,&lastUsed)
		if err!=nil{
			return 0,err
		}
		if lastUsed<oldest{
			bestId=id
			oldest=lastUsed
		}
	}
	return bestId,nil
}


