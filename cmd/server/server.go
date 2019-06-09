package server

import (
	"fmt"
	"os"

	"github.com/comomac/shin-kamishibai/server/pkg/config"
	"github.com/comomac/shin-kamishibai/server/pkg/fdb"
	svr "github.com/comomac/shin-kamishibai/server/server"
)

func main() {
	// fmt.Println(genChar(3))

	// convert format
	// jfile := userHome("etc/kamishibai-kai/db.json")
	// tfile := userHome("etc/shin-kamishibai/db.txt")
	// convJtoF(jfile, tfile) // json to txt
	// convFtoJ(tfile, jfile) // txt to json

	// // load db
	// db := NewFlatDB(userHome("etc/shin-kamishibai/db.txt"))
	// db.Load()

	// fmt.Println(db.BookIDs())
	// fmt.Println(db.GetBookByID("7IL"))

	// // export database, check if it goes generate proper flat db
	// db.Export(userHome("etc/shin-kamishibai/db2.txt"))
	// ibook := db.IBooks[100]
	// fmt.Printf("%+v %+v\n", ibook, ibook.Book)

	// # test if page update works and only update 4 bytes instead of everything
	// x, err := db.UpdatePage("7IL", 9876)
	// check(err)
	// fmt.Println(x)

	// use config on local dir by default if no param given
	cfgFilePath := "config.json"
	if len(os.Args) > 1 {
		cfgFilePath = os.Args[1]
	}

	cfg, err := config.Read(cfgFilePath)
	if err != nil {
		fmt.Println("faile to read config file")
		panic(err)
	}

	// new db
	// db := NewFlatDB(userHome("etc/shin-kamishibai/db.txt"))
	db := fdb.New(cfg.DBPath)
	db.Load()
	for _, dir := range cfg.AllowedDirs {
		fdb.AddDir(db, dir)
	}

	svr.Start(cfg, db)
}

// 22849