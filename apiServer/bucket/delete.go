package bucket

import (
	es "DistributedFileSystem/common/lib/ElasticSearch"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	// 获得桶名
	bucket := r.Header.Get("bucket")
	if bucket == "" {
		code := http.StatusBadRequest
		msg := "请求头缺少 bucket 字段"
		golog.Error.Println(msg)
		w.WriteHeader(code)
		return
	}

	err := es.DelBucket(bucket)
	if err != nil {
		code := http.StatusInternalServerError
		msg := "删除 bucket 时出错：" + err.Error()
		golog.Error.Println(msg)
		w.WriteHeader(code)
		return
	}
	code := http.StatusOK
	msg := "删除 bucket: " + bucket + " 成功"
	golog.Info.Println(msg)
	w.WriteHeader(code)
}
