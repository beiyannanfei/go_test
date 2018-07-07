# loveworld crontab job

- A. 配置文件相关

    config目录，根据 default.json.bak 文件生成 default.json 文件，具体释义如下:
    ```json
      {
        "platform": "qq-android",                 //平台信息
        "mysql": {                                //mysql配置信息(需要和gs相同)
          "host": "127.0.0.1",
          "port": 3306,
          "database": "loveworld_develop",
          "username": "root",
          "password": "",
          "tableName": {                          //用到的表名称
            "table_mail": "user_mail",
            "table_avgBook": "user_avg_book",
            "table_avgData": "user_avg_data",
            "table_device": "user_device"
          }
        },
        "redis": {                                //redis配置信息(需要和gs相同)
          "host": "127.0.0.1",
          "port": 6379,
          "database": 0,
          "password": ""
        },
        "logger": {
          "logDir": "/logs",
          "level": "debug",
          "output": true
        },
        "crontab": {                              //job相关信息
          "comment": "lovejob-android-wechat",    //job组信息    
          "nodeCommand": "/usr/local/bin/node",   //which node node命令所在目录
          "shellCommand": "/bin/bash",            //which bash  bash命令所在目录
          //定时任务相关，格式：fileName: * * * * *
          "sendRankReward": "0 0 1 * *",          //每月初发奖
          "initCharacterRank": "0 0 1 * *",       //初始化每个月的新排行榜
          "sendXingeMsg": "0 20 * * *"            //每晚8点发放信鸽推送
        },
        "xinGe": {
          //信鸽相关
          "accessId": 2100003306,
          "secretKey": "46cfadd7101f1dc4d1e2eee5e07632fe",
          "ios_env": "env"                        //推送ios信鸽消息区分测试环境和正式环境(测试："env", 正式: "pro")
        }
      }
    ```
    
- B. nodejs安装，要求版本 >8.0.0 建议：8.10.0
    ``` bash
    wget https://nodejs.org/dist/v8.10.0/node-v8.10.0-linux-x86.tar.gz
    tar -zxf node-v8.10.0-linux-x86.tar.gz -C /usr/local/
    cd /usr/local/
    ln -s node-v8.10.0-linux-x86/ node
    cat>/etc/profile<<EOF
    
    # set for nodejs
    export NODE_HOME=/usr/local/node
    export PATH=$NODE_HOME/bin:$PATH
    EOF
    source  /etc/profile
    node -v
    ```
- C. 安装cnpm
    
        npm i cnpm -g
    
- D. 安装依赖库

        cd /home/cronjob/ 
        cnpm i
        
- E. 执行脚本创建cronjob
    
        cd /home/cronjob/
        node app
        
- F. 检查 crontab 任务(类似下表)

        @monthly cd /Users/wyq/workspace/lovejob/ && /usr/local/bin/node /Users/wyq/workspace/lovejob/app/initCharacterRank.js #lovejob-android-wechat
        @monthly cd /Users/wyq/workspace/lovejob/ && /usr/local/bin/node /Users/wyq/workspace/lovejob/app/sendRankReward.js #lovejob-android-wechat
        0 20 * * * cd /Users/wyq/workspace/lovejob/ && /usr/local/bin/node /Users/wyq/workspace/lovejob/app/sendXingeMsg.js #lovejob-android-wechat

- G. 关于log

    crontab脚本会在执行时产出log，保存于logs目录，log文件名称格式: 
        
        YYYYMMDD-filename.log
        eg.
            20180421-sendRankReward.log
            20180421-sendRankReward.log