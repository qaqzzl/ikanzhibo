    v1.0
        完成基本功能
    v1.1
        优化代码
        解决内存泄漏 <-time.Tick() 导致的
    v1.2
        优化代码
        新增字段,直播间权重,分类权重






### 目录说明
    ```
    Spider/
    ├── parser
    │   ├── huya.go     解析器
    │   └── douyin.go   解析器
    ├── downloader.go      下载器
    ├── master.go          任务调度器
    ├── main.go            主程序
    │  
    └── README.md
    ```

    



vendor     --depth=1
```text
    git clone https://github.com/shirou/gopsutil.git

    git clone https://github.com/golang/sys.git
    https://github.com/golang/net.git
    https://github.com/golang/text.git
```    