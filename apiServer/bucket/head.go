package bucket

import (
	es "DistributedFileSystem/common/lib/ElasticSearch"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Head(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 获得桶名
	bucket := r.Header.Get("bucket")

	if bucket == "" {
		golog.Error.Println("请求头缺少bucket字段")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 查找bucket
	msg := "查询到bucket：" + bucket + "存在"
	code := es.SearchBucket(bucket)
	golog.Info.Println(msg)
	w.WriteHeader(code)
}
