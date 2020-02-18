package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Server holds link to database and configuration
type Server struct {
	Database *FlatDB
	Config   *Config
}

// Start launches http server
func (svr *Server) Start() {
	cfg := svr.Config
	db := svr.Database

	// setup session
	httpSession := &SessionStore{}

	h := http.NewServeMux()

	// public folder access

	// debug
	fserv := http.FileServer(http.Dir(svr.Config.PathDir + "/web"))

	// generated and packed
	// fs := fileSystem{__binmapName}
	// fserv := http.FileServer(fs)

	h.HandleFunc("/", handlerFS(fserv))

	// public api, page
	h.HandleFunc("/login", loginPOST(httpSession, cfg))
	h.HandleFunc("/login.html", loginGet(cfg, db))

	// private api, page
	h.HandleFunc("/api/thumbnail/", renderThumbnail(db, cfg)) // /thumbnail/{bookID}    get book cover thumbnail
	h.HandleFunc("/api/cbz/", getPageOnly(db))                // /cbz/{bookID}/{page}   get book page
	h.HandleFunc("/browse.html", browseGet(cfg, db, "ssp/browse.ghtml"))
	h.HandleFunc("/legacy.html", browseGet(cfg, db, "ssp/legacy.ghtml"))
	h.HandleFunc("/read.html", readGet(cfg, db))

	// middleware
	slog := svrLogging(h, httpSession, cfg)
	h1 := CheckAuthHandler(slog, httpSession, cfg)

	port := cfg.IP + ":" + strconv.Itoa(cfg.Port)
	fmt.Println("listening on", port)
	fmt.Println("allowed dirs: " + strings.Join(cfg.AllowedDirs, ", "))
	log.Fatal(http.ListenAndServe(port, h1))
}
