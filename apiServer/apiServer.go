package main

import (
	"DistributedFileSystem/apiServer/bucket"
	"DistributedFileSystem/apiServer/heartbeat"
	"DistributedFileSystem/apiServer/locate"
	"DistributedFileSystem/apiServer/logs"
	"DistributedFileSystem/apiServer/objects"
	"DistributedFileSystem/apiServer/system"
	"DistributedFileSystem/apiServer/temp"
	"DistributedFileSystem/apiServer/tools"
	"DistributedFileSystem/apiServer/versions"
	RedisMQ "DistributedFileSystem/common/lib/Redis"
	"DistributedFileSystem/common/lib/golog"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		defer func() {
			if err := recover(); err != nil {
				golog.Error.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}

func InitRouter(r *gin.Engine) {
	r.Use(Cors())
	//heartbeat
	{
		r.GET("/heartbeat", heartbeat.Get)
	}
	//system
	{
		//这里的*id都是dataServer的ip地址,如:127.0.0.1:8081
		r.GET("/nodeSystemInfo/*id", system.NodeGet)
		r.GET("/systemInfo/*id", system.UseGet)
	}
	//bucket
	{
		r.GET("/bucket/*id", bucket.Get)
		r.PUT("/bucket/*id", bucket.Put)
		r.DELETE("/bucket/*id", bucket.Delete)
		r.HEAD("/bucket/*id", bucket.Head)
	}
	//logs
	{
		r.POST("/getLog/*id", logs.Post)
	}
	//locate
	{
		r.GET("/locate/*id", locate.Get)
	}
	//objects
	{
		//put 函数和 get 函数并不会访问本地磁盘上的对象，而是将 HTTP 请求转发给数据服务
		r.PUT("/objects/*id", objects.Put) //上传Size较小的对象文件
		r.GET("/objects/*id", objects.Get)
		r.DELETE("/objects/*id", objects.Delete)
		r.POST("/objects/*id", objects.Post) //获得分片上传Token
	}
	//versions
	{
		/*
			两种不同的接口
			http://apiServerIP/versions/ 查看所有对象的所有版本信息
			http://apiServerIP/versions/<xxx> 查看指定单个对象的所有版本信息
		*/
		r.GET("/versions/*id", versions.Get)
		r.GET("/allVersions/*id", versions.AllGet) //获得所有对象最新版本的分页查询列表，约定每页显示5条
	}
	//temp
	{
		r.HEAD("/temp/*id", temp.Head) //获得分片上传偏移量
		r.PUT("/temp/*id", temp.Put)   //分片上传: 适用于上传Size较大的对象文件，并且约定每次分片数据为5MB，使用需要先获得了分片上传的Token
	}
	//tools
	{
		//同一对象，最多留存5个历史版本，最早的版本会被删掉
		r.GET("/deleteOldMetadata/*id", tools.DeleteOldMetaDate) // 删除过期元数据的工具
		r.GET("/deleteOrphanServer/*id", tools.DeleteOrphan)     // 删除没有元数据引用的对象数据
		r.GET("/objectScanner/*id", tools.ObjectScanner)         //扫描所有节点对象文件，如果如果可以修复的执行修复操作
	}
}

func main() {
	// 实时读取日志
	go golog.ReadLog(time.Now().Format("2006-01-02"))

	RedisMQ.Rds = RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer RedisMQ.Rds.Client.Close()

	// 开始连接apiServers这个exchanges，将数据服务节点的地址保存起来
	go heartbeat.ListenHeartbeat()

	r := gin.Default()
	InitRouter(r)

	listenAddress := os.Getenv("LISTEN_ADDRESS")
	golog.Info.Println(os.Getenv("LISTEN_ADDRESS"), "===> apiServer Start running <===")

	//监听并启动 ip在tools中规划好了
	golog.Info.Println(r.Run(listenAddress))
}
