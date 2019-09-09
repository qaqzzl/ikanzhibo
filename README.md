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

    

#### 住控制器
    任务调度器
        分配任务,去重
        
        
        
        
        
        
        


#### Redis队列
    * 关注不在线
    未关注不在线
    在线
    全平台