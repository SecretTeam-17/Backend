package game

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist
var dist embed.FS

func Serve() http.FileSystem {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		fmt.Println(err)
	}

	return http.FS(sub)
}
