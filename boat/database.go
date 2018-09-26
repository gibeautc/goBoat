package boat

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
	_,err:=self.conn.Exec("CREATE TABLE IF NOT EXISTS tiles(id INTEGER PRIMARY KEY AUTOINCREMENT ,onDisk INTEGER,comp INTEGER,inRam INTEGER,lastUsed INTEGER)")
	if err!=nil{
		return err
	}
	return nil
}


func (self *TileSet) getNewTileID() (uint32,error){
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
	onDisk:=Exists("tiles/"+strconv.Itoa(int(tile.Id)))
	_,err:=self.conn.Exec("REPLACE INTO tiles(id,onDisk,comp,inRam,lastUsed) VALUES($1,$2,$3,$4,$5)",tile.Id,onDisk,len(tile.Data),index,time.Now().Unix())
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



