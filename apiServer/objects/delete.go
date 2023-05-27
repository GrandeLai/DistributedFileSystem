package objects

import (
	es "DistributedFileSystem/common/lib/ElasticSearch"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Delete(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	// 获得桶名
	bucket := strings.Split(r.URL.EscapedPath(), "/")[2]
	if bucket == "" {
		golog.Error.Println("url 缺少 bucket 字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 获取对象名称
	name := strings.Split(r.URL.EscapedPath(), "/")[3]

	//从es中获取object的最新版本
	version, err := es.SearchLatestVersion(bucket, name)
	if err != nil {
		golog.Error.Println("es search latest version err：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if version.Version == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//向es中插入object一个版本，size=0,hash=""，因为在get时如果哈希值为空则说明被标记为删除
	err = es.PutMetadata(bucket, name, version.Version+1, 0, "")
	if err != nil {
		golog.Error.Println("es put metadata err：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	golog.Info.Println("删除 object： 成功", name)
	w.WriteHeader(http.StatusOK)
}
