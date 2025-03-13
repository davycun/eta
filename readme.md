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