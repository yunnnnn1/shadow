# shadow

#### 介绍
基于golang实现的C/S架构agent,使用gRPC协议,实现执行远程命令就像在本地执行一样,能够保证命令一边执行,一边返回执行结果,支持远程shell与python命令执行，支持并发执行。

#### 安装说明
编译Linux,Windows和Mac环境下可执行程序

```go
go get -u github.com/jingmingyu/shadow
go mod tidy 
```

######linux
```go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/Sserver server/main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/Sclient client/main.go
```
######windows
```go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/Sserver server/main.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/Sclient client/main.go
```
######mac
```go
CGO_ENABLED=0 GOOS=mac GOARCH=amd64 go build -o bin/mac/Sserver server/main.go

CGO_ENABLED=0 GOOS=mac GOARCH=amd64 go build -o bin/mac/Sclient client/main.go
```

#### 使用说明
目前shadow支持远程命令执行和文件下载的功能

远程执行：   
1.登录远程服务器 (默认使用3721端口)

    shell>./Sserver -port 3721 

2.登录本地服务器，拷贝本地测试脚本到Server端，测试客户端,此处以192.168.2.178为远程server端举例子

    shell> scp ./script/test.py 192.168.2.178:/tmp/test.py
    shell> ./Sclient EXEC -Rserver 192.168.2.178:3721 -user oracle -cmd "python -u /tmp/test.py"  

3.查看本地日志返回

    shell>print_test
          2022-01-25 21:52:39,671 - test - INFO     - logger_test
          print_test
          2022-01-25 21:52:40,676 - test - INFO     - logger_test
          print_test
          2022-01-25 21:52:41,680 - test - INFO     - logger_test
          print_test
          2022-01-25 21:52:42,686 - test - INFO     - logger_test

4.查看远程日志

    shell>2022/01/25 21:50:31 [Shadow Server] Start at Pid:5966 , Listen at Port:3721              
          2022/01/25 21:52:39 [Info]-->Command: python -u /tmp/test.py

文件下发：

目前文件下载只支持http的方式 暂不支持https

1.启动server

    shell>./Sserver -port 3721 

2.启动http服务,以python举例

    shell> python -m SimpleHTTPServer 8000 
           Serving HTTP on 0.0.0.0 port 8000 ... 

3.客户端进行文件下发

    shell> ./Sclient FILE -Rserver 192.168.2.178:3721 -Cip 192.168.2.178:8000 -FileName test -Path /temp/ -IsRecover N -timeout 60 -limitspeed 60

* 参数详解：
  ##### -Cip 192.168.2.178:8000 软件包所在的主机ip/port
  ##### -FileName               软件名,若有多级目录则按照 Yourdir1/Yourdir2/Yourfile;若需要传输多个文件 Yourdir1/Yourdir2/Yourfile*
  ##### -Path                   下载文件到远程主机的文件路径
  ##### -IsRecover              是否覆盖 Y/N
  ##### -timeout                软件服务器响应时间 单位 s
  ##### -limitspeed             文件下载速度 单位 kb


#### NOTE：
1. 执行python程序时，如果程序中使用到了print，在调用shadow时，必须要加入 -u 参数。原因：print默认输出是带有缓冲区,即只有当标准输出写满整个缓冲区时，才能将结果进行返回，-u 即不使用缓冲区，直接将结果返回标准输出。
2. 远程server以root启动(涉及到指定用户执行命令)，默认以127.0.0.1启动,目前已支持指定用户进行命令执行。
3. 目前远程执行命令都是以同步的方式进行运行，即不使用 nohup sh +x xxxx.sh & 方式。

#### ToDo ：
1. 密码验证
2. 服务端的守护进程

####  提问 ：
如果使用过程中有问题,欢迎提交问题至 `958200673@qq.com`