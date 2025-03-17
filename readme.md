### how to import project
- 配置环境变量
```shell
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
export GOPRIVATE=gitlab.xxx.com
```

- 下载依赖包
```shell
go mod tidy
```

### 启动应用
1. 在本地启动一个redis 服务，并且配置deploy/config_dameng.yml 中的redis选项
2. 启动方式1：command start: go run main.go server -c deploy/config_dameng.yml
3. 启动方式2：goland（IDE）方式
   1. 新建一个go build，Package path选择 github.com/davycun/eta
   2. 在本地启动一个redis 服务
   3. 在 program arguments 中添加 server -c deploy/config_dameng.yml

### 关于ES的fingerprint或者CaCert配置说明
- https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html
```shell
openssl s_client -connect 172.18.54.188:39001 -servername 172.18.54.188 -showcerts </dev/null 2>/dev/null | openssl x509 -fingerprint -sha256 -noout -in /dev/stdin
```


### 监控信息
- http://localhost:6060/debug/vars
- http://localhost:6060/debug/pprof/

### 测试说明
跑测试需要再项目目录下新建一个config_local.yml的配置文件，示例如下
```yaml
migrate: true
monitor:
   port: 6060
   id:
   node_id: 1
   epoch: 2023-01-01
server:
   host: 0.0.0.0
   port: 8080
   #  gin_mode: debug
   gin_mode: release
   api_doc_enable: true
   ignore_uri:
      - /oauth2/*
      - /storage/download/*
      - /storage/upload/*

redis:
   host: 127.0.0.1
   port: 16379
   password: abc@123
   db: 0
database:
   host: 127.0.0.1
   port: 15432
   user: postgres
   password: 123456
   dbname: postgres
   schema: eta
   type: postgres
   # 4:info，3:warn，2:error，1:silent
   log_level: 4
   # 打印慢sql，单位是毫秒
   slow_threshold: 200
```