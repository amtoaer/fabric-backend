package web

import (
	"net/http"
	// static files
	_ "github.com/amtoaer/fabric-backend/statik"
	"github.com/gin-contrib/static"
	"github.com/rakyll/statik/fs"
)

type ginFS struct {
	FS http.FileSystem
}

// Open 打开文件
func (b *ginFS) Open(name string) (http.File, error) {
	return b.FS.Open(name)
}

// Exists 文件是否存在
func (b *ginFS) Exists(prefix string, filepath string) bool {
	if _, err := b.FS.Open(filepath); err != nil {
		return false
	}
	return true
}

func getFileSystem() static.ServeFileSystem {
	var StaticFS static.ServeFileSystem
	StaticFS = &ginFS{}
	StaticFS.(*ginFS).FS, _ = fs.New()
	return StaticFS
}
