# [Golang]Testcontainers-for-elasticsearch
[Testcontainers](https://java.testcontainers.org/modules/elasticsearch/)を使ったElasticsearchのテスト環境の構築と、[go-elasticsearch](https://github.com/elastic/go-elasticsearch)を使った単体テストのサンプルです。

## 開発環境
```bash
$ go version
go version go1.22.2 darwin/amd64
```

## 単体テストの実行
Elasticsearchへの検索リクエスト(`Search`メソッド)に対する単体テストを実行します。

```bash
$ go test -count=1 -v -run TestSearch ./elasticsearch/elasticsearch_test.go 
=== RUN   TestSearch
2024/05/11 00:27:34 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 20.10.7
  API Version: 1.41
  Operating System: Docker Desktop
  Total Memory: 7961 MB
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 5a114182e2c75129cfdcc16b9c2b5bfa3ac06f8ace09c02d1cb2e17e25f54ee1
  Test ProcessID: d875f463-4b69-4fac-a1b1-d4b22eb023d1
2024/05/11 00:27:34 🐳 Creating container for image testcontainers/ryuk:0.7.0
2024/05/11 00:27:34 ✅ Container created: 63a385e183a3
2024/05/11 00:27:34 🐳 Starting container: 63a385e183a3
2024/05/11 00:27:34 ✅ Container started: 63a385e183a3
2024/05/11 00:27:34 🚧 Waiting for container id 63a385e183a3 image: testcontainers/ryuk:0.7.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2024/05/11 00:27:35 🔔 Container is ready: 63a385e183a3
2024/05/11 00:27:35 🐳 Creating container for image docker.elastic.co/elasticsearch/elasticsearch:8.7.1
2024/05/11 00:27:35 ✅ Container created: 3af53534b0ac
2024/05/11 00:27:35 🐳 Starting container: 3af53534b0ac
2024/05/11 00:27:35 ✅ Container started: 3af53534b0ac
2024/05/11 00:27:35 🚧 Waiting for container id 3af53534b0ac image: docker.elastic.co/elasticsearch/elasticsearch:8.7.1. Waiting for: &{timeout:<nil> Log:.*("message":\s?"started(\s|")?.*|]\sstarted\n) IsRegexp:true Occurrence:1 PollInterval:100ms}
2024/05/11 00:27:55 🔔 Container is ready: 3af53534b0ac
2024/05/11 00:27:55 index test_index not found
=== RUN   TestSearch/正常系(200)
2024/05/11 00:27:56 index test_index is already exits
=== RUN   TestSearch/異常系(400-BadRequest)
2024/05/11 00:27:56 index test_index is already exits
=== RUN   TestSearch/異常系(404-IndexNotFound)
2024/05/11 00:27:57 🐳 Terminating container: 3af53534b0ac
2024/05/11 00:27:57 🚫 Container terminated: 3af53534b0ac
--- PASS: TestSearch (23.47s)
    --- PASS: TestSearch/正常系(200) (0.15s)
    --- PASS: TestSearch/異常系(400-BadRequest) (0.03s)
    --- PASS: TestSearch/異常系(404-IndexNotFound) (0.03s)
PASS
ok  	command-line-arguments	24.165s
```
