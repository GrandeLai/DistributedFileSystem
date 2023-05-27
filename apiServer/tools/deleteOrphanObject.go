package tools

import (
	es "DistributedFileSystem/common/lib/ElasticSearch"
	RedisMQ "DistributedFileSystem/common/lib/Redis"
	"DistributedFileSystem/common/lib/golog"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func delete(hash string) {

	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	res, err := client.Do(request)
	if res.StatusCode == http.StatusOK {
		fmt.Println("delete:", hash)
	} else if err != nil {
		fmt.Println(err)
	}
}

func DeleteOrphan(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	defer r.Body.Close()

	buckets := es.GetAllBucket()
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		realhash, _ := url.PathUnescape(hash)
		flag := false
		for _, bucket := range buckets {
			hashInMetadata, err := es.HasHash(bucket, realhash)
			if err != nil {
				log.Println(err)
				return
			}
			if hashInMetadata {
				flag = true
				break
			}
		}
		if !flag {
			delete(realhash)
		}
	}
	w.WriteHeader(http.StatusOK)
	rdb := RedisMQ.Rds
	rdb.Incr("DFSUpHold")
	op := fmt.Sprintf("进行了删除无元数据引用的文件的操作")
	date := time.Now().Format("2006-01-02")
	time := time.Now().Format("15:04:05")
	rdb.InsertOp(op, date, time)
	golog.Trace.Println("删除无元数据引用的对象数据")
}
