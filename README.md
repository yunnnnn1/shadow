# shadow

#### 介绍
基于golang实现的轻量级agent，使用RPC与http协议构建内部API，实现执行远程命令就像在本地执行一样，能够保证命令一边执行，一边返回执行结果，
支持远程shell与python命令执行，支持并发执行。

#### 使用说明

   1.登录远程服务器 (默认使用3721端口)


    shell>./Sserver -port 3721 

   2.登录本地服务器，拷贝本地测试脚本到Server端，测试客户端,此处以192.168.2.178为远程server端举例子

        
    shell> scp ./script/test.py 192.168.2.178:/tmp/test.py
    shell> ./Sclient  -Rserver 192.168.2.178:3721 -cmd "python -u /tmp/test.py"  

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

#### NOTE： 
* BASE版
1. 执行python程序时，如果程序中使用到了print，在调用shadow时，必须要加入 -u 参数。原因：print默认输出是带有缓冲区，
   即只有当标准输出写满整个缓冲区时，才能将结果进行返回，-u 即不使用缓冲区，直接将结果返回标准输出。
2. 远程server启动时，默认以127.0.0.1启动
3. 目前远程执行命令都是以同步的方式进行运行，即不使用 nohup sh +x xxxx.sh & 方式。
4. 默认以root启动，目前已支持指定用户进行命令执行

