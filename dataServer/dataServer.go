package main

import (
	RedisMQ "DistributedFileSystem/common/lib/Redis"
	"DistributedFileSystem/common/lib/golog"
	"DistributedFileSystem/dataServer/heartbeat"
	"DistributedFileSystem/dataServer/locate"
	"DistributedFileSystem/dataServer/objects"
	"DistributedFileSystem/dataServer/system"
	"DistributedFileSystem/dataServer/temp"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func main() {
	// 实时读取日志
	go golog.ReadLog(time.Now().Format("2006-01-02"))

	RedisMQ.Rds = RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer RedisMQ.Rds.Client.Close()

	//第一次启动时，将所有对象存储到map
	locate.CollectObjects()
	//开始发送心跳包
	go heartbeat.StartHeartbeat()

	r := gin.Default()
	InitRouter(r)

	golog.Info.Println(os.Getenv("LISTEN_ADDRESS"), "===>  dataServer Start running  <===")
	//监听并启动 ip在tools中规划好了
	//目前是10.29.1.1和10.29.1.6
	golog.Info.Println(r.Run(os.Getenv("LISTEN_ADDRESS")))
}

func InitRouter(r *gin.Engine) {
	//system
	{
		r.GET("/systemInfo", system.Get)
	}
	//temp
	{
		r.POST("/temp/*id", temp.Post)
		r.PATCH("/temp/*id", temp.Patch) //将数据暂存下来，等待转正，并进行数据校验
		r.PUT("/temp/*id", temp.Put)     //转正，将$STORAGE_ROOT/temp/t.Uuid.dat 改为 $STORAGE_ROOT/objects/hash
		r.DELETE("/temp/*id", temp.Delete)
		r.HEAD("/temp/*id", temp.Head) //返回对象的大小
		r.GET("/temp/*id", temp.Get)   //获取对象本身
	}
	//objects
	{
		r.GET("/objects/*id", objects.Get)
		r.DELETE("/objects/*id", objects.Delete)
	}
}
