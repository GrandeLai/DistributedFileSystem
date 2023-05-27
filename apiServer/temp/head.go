package temp

import (
	"DistributedFileSystem/common/lib/golog"
	"DistributedFileSystem/common/lib/rs"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Head(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	token := strings.Split(r.URL.EscapedPath(), "/")[3] //token为获得分片上传Token中的Token
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		golog.Error.Println("new rs put stream err：", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//获取数据节点已经储存该对象多少数据了
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current)) //current是已经上传的字节
	golog.Info.Println("获取数据节点已经储存该对象多少数据成功")
}
