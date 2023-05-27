package objects

import (
	RedisMQ "DistributedFileSystem/common/lib/Redis"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Delete(ctx *gin.Context) {
	r := ctx.Request

	//获取hash
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := url.PathEscape(object)
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
	if len(files) != 1 {
		return
	}
	//locate.Del(hash)
	rdb := RedisMQ.Rds
	rdb.RemoveFile(hash, os.Getenv("LISTEN_ADDRESS"))

	err := os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0])) //直接将要删除的对象从objects目录下移到garbage下
	if err != nil {
		golog.Error.Println("rename err：", err)
	}
}
