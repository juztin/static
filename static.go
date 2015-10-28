package main

import (
	"flag"
	"fmt"
	"net/http"
)

var portFlag = flag.Int("port", 8080, "Port number")
var folderFlag = flag.String("path", "./", "Static root")

func serve(root string, port int) {
	addr := fmt.Sprintf(":%d", port)
	fileServer := http.FileServer(http.Dir(root))
	fmt.Printf("Service static for \"%s\" on :%d.\n", root, port)
	http.ListenAndServe(addr, fileServer)
}

func main() {
	flag.Parse()
	serve(*folderFlag, *portFlag)
}
