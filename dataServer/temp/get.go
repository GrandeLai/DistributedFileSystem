package temp

import (
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
)

func Get(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
	if err != nil {
		golog.Error.Println("open file err: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()
	io.Copy(w, file)
}
