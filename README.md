# 一.初始化 后端

## 0.安装go语言环境
     https://www.yuque.com/docs/share/640ece10-1927-4e06-9c92-c3cab0352490

    没有科学上网的同学，配置go代理，打开以下链接参考配置
    https://goproxy.cn/

## 1.初始化项目并安装依赖
    git clone https://gitee.com/Berners/file.git
    cd file
    go install 

## 2.配置参数，包括数据库,OSS,...
    lib/config/index.go

## 3.初始化表
    cd dao/
    go test -v -run TestDaoInit_CreateFileTable DaoInit_test.go

## 4.启动项目,命令行根目录执行
    方式一：当前命令行启动,关闭命令行退出，用于开发环境
        go run main.go

    方式二：后台启动，用于生产环境
        1.编译可执行文件
        go build -o file-server main.go

        2.后台启动
        sudo nohup ./file-server > nohup_bluebell.log 2>&1 &
            备注
            ./file-server 是我们应用程序的启动命令
            nohup ... & 表示在后台不挂断的执行上述应用程序的启动命令
            > nohup_bluebell.log 表示将命令的标准输出重定向到 nohup_bluebell.log 文件
            2>&1 表示将标准错误输出也重定向到标准输出中，结合上一条就是把执行命令的输出都定向到 nohup_bluebell.log 文件
    
        3.停止
            3.1 查看软件进程
                ps -ef|grep file-server
            3.2 找到 file-server 运行进程id

            3.3 根据进程id执行杀掉
                sudo kill PID

# 二.初始化前端,打开地址查看
    https://gitee.com/Berners/file-front
