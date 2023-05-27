package system

import (
	"DistributedFileSystem/common/lib/golog"
	"DistributedFileSystem/common/lib/system"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()
	info := system.GetInfo()

	marshal, err := json.Marshal(info)
	if err != nil {
		golog.Error.Println("json marshal err：", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(marshal)
}
