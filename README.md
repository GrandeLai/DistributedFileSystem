# Distributed File system

分布式文件系统



启动方法：

- 根目录创建/log/n，n表示你apiServer和dataServer的总个数

- 根目录创建/tmp/n/garbage，/tmp/n/objects，/tmp/n/temp，n表示你dataServer的总个数

- 启动：

  ```
  export LOG_DIRECTORY=/log/1 LISTEN_ADDRESS=127.0.0.1:8091 STORAGE_ROOT=C:/tmp/1 go run dataServer/dataServer.go &
  export LOG_DIRECTORY=/log/2 LISTEN_ADDRESS=127.0.0.1:8092 STORAGE_ROOT=C:/tmp/2 go run dataServer/dataServer.go &
  export LOG_DIRECTORY=/log/3 LISTEN_ADDRESS=127.0.0.1:8093 STORAGE_ROOT=C:/tmp/3 go run dataServer/dataServer.go &
  export LOG_DIRECTORY=/log/4 LISTEN_ADDRESS=127.0.0.1:8094 STORAGE_ROOT=C:/tmp/4 go run dataServer/dataServer.go &
  export LOG_DIRECTORY=/log/5 LISTEN_ADDRESS=127.0.0.1:8095 STORAGE_ROOT=C:/tmp/5 go run dataServer/dataServer.go &
  export LOG_DIRECTORY=/log/6 LISTEN_ADDRESS=127.0.0.1:8096 STORAGE_ROOT=C:/tmp/6 go run dataServer/dataServer.go &
  
  export LOG_DIRECTORY=/log/7 LISTEN_ADDRESS=127.0.0.1:8081 go run apiServer/apiServer.go &
  ```

  