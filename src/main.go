package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

var (
	masterAddr  = flag.String("master", "", "RPC master address")
	rpcEnabled  = flag.Bool("rpc", false, "enable RPC server")
	listenAddr  = flag.String("http", ":8080", "http listen address")
	dataFile    = flag.String("file", "store.json", "data store file name")
	hostname    = flag.String("host", "localhost:8080", "http host name")
	server_type = "master"
)

var store Store

func main() {
	flag.Parse()
	if *masterAddr != "" { // slave server
		server_type = "slave"
		store = NewProxyStore(*masterAddr)
	} else { // master server
		store = NewURLStore(*dataFile)
	}

	store = NewURLStore(*dataFile)
	if *rpcEnabled { // register RPC service in master server
		rpc.RegisterName("Store", store)
		rpc.HandleHTTP()
	}
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)
	http.ListenAndServe(*listenAddr, nil)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	log.Printf("%s redirecting to %s", server_type, key)
	var url string
	if err := store.Get(&key, &url); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.Printf("%s adding: %s", server_type, url)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if url == "" {
		fmt.Fprint(w, AddForm)
		return
	}
	var key string
	if err := store.Put(&url, &key); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "http://%s/%s", *hostname, key)
}

const AddForm = `
<form method="POST" action="/add">
URL: <input type="text" name="url">
<input type="submit" value="Add">
</form>
`
