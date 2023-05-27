package bucket

import (
	es "DistributedFileSystem/common/lib/ElasticSearch"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Put(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	bucket := r.Header.Get("bucket")
	if bucket == "" {
		msg := "请求头缺少bucket字段"
		code := http.StatusBadRequest
		golog.Error.Println(msg)
		w.WriteHeader(code)
		return
	}
	err := es.AddBucket(bucket)
	if err != nil {
		golog.Error.Println("增加bucket出错：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg := "增加bucket: " + bucket + " 成功"
	code := http.StatusCreated
	golog.Info.Println(msg)
	w.WriteHeader(code)
}
