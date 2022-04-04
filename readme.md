## 设计表结构
> 使用 https://dbdiagram.io 设计
## 安装docker
```bash
docker pull postgres:alpine
docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine
docker exec -it postgres psql -U root
docker logs postgres
```
docker rm 不能与-d一起用
## 安装TablePlus
## 迁移数据库
```bash
go get -u -d github.com/golang-migrate/migrate/cmd/migrate
migrate create -ext sql -dir .\db\migrations -seq InitDatabase
# postgres创建数据库
createdb --username=root --owner=root simple_bank
dropdb simple_bank
# driver postgres://user:password@host:port/dbname?query
migrate -path db/migrations -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" up
```
## 操作数据库
- database/sql
  
  原生操作数据库
- gorm

  关系映射操作
- sqlx
- sqlc
## 安装sqlc
```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
# sqlc 只能在linux上生成，可以采用wsl
sqlc init # 生成yaml配置文件
```
sqlc配置文件
```yaml
version: 1
packages:
  - path: "db/sqlc"
    name: "db"
    engine: "postgresql"
    schema: "db/migrations"
    queries: "db/query"
    emit_json_tags: true
    emit_prepared_queries: false
    emit_interface: false
    emit_exact_table_names: false
```
编写query sql, 运行sqlc generate
## 单元测试
测试的入口是TestMain, database/sql提供操作接口
```bash
go test simplebank-app/util -v -run TestPassword
go get -u github.com/lib/pq
go get -u github.com/strechr/testify
```
## 事务
```sql
BEGIN;
COMMIT;
BEGIN;
ROLLBACK;
SELECT * FROM account WHERE id = 1 FOR UPDATE; -- 事务之间会阻塞查询
-- 阻塞查询会导致无差别查询停止
-- 查询被阻塞的语句
SELECT blocked_locks.pid     AS blocked_pid,
       blocked_activity.usename  AS blocked_user,
       blocking_locks.pid     AS blocking_pid,
       blocking_activity.usename AS blocking_user,
       blocked_activity.query    AS blocked_statement,
       blocking_activity.query   AS current_statement_in_blocking_process
FROM  pg_catalog.pg_locks         blocked_locks
        JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
        JOIN pg_catalog.pg_locks         blocking_locks
             ON blocking_locks.locktype = blocked_locks.locktype
               AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
               AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
               AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
               AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
               AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
               AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
               AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
               AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
               AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
               AND blocking_locks.pid != blocked_locks.pid

        JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;
```
有外键关联的依赖表会阻塞事务下的主表查询
查询时，忽略键值更新
```sql
SELECT * FROM accounts WHERE id = 1 FOR NO KEY UPDATE;
```
无缓冲的channel需要保持入队与出队的顺序
造成数据库死锁：
1. 外键约束导致死锁
2. 对两个账户交替操作导致死锁
## 数据隔离级别
低级别的隔离导致读现象
- 脏读
  
  一个事务读取到其他事务未提交的数据
- 不可重复读

  一个事务读取相同行，但是两次获取数值不一样，值的不同是因为被其他事务提交后修改
- 幻读
  由于事务提交，导致新增，重复执行查询得到不同数据集
- 序列错误
## 数据隔离标准
- 读未提交的数据
  
  可以看到未提交事务的数据
- 读提交的数据
  
  只能看到事务提交后的数据
- 重复读

  相同查询返回相同的结果
- 可序列
  事务之间都是顺序的，不存在每个事务之间操作重叠
## 设置数据隔离级别
- mysql
mysql默认为不可重复读

```sql
SELECT @@transaction_isolation;
SELECT @@global.transaction_isolation;
set session transaction isolation level read uncommitted | read committted | repeatable read | serializable;
```
read uncommited  
其他会话可以读取到没有提交的修改数据  
read commited  
避免脏读，会话不读取读取其他会话未提交的数据  
repeatable read  
其他会话的修改提交，不影响当前会话，查询结果保持一致，但是如果在当前会话修改，会叠加其他事务修改的结果  
serializable  
事务之间的操作都是序列的，每次只允许执行一个操作，事务重试机制，容易发生死锁  
postgres  
```sql
show transaction isolation level; -- 默认 read commited
-- 只能在begin范围内设置
set session transaction isolation level read uncommitted | read committted | repeatable read | serializable;
```
postgres read uncommitted 不可脏读
postgres repeatable read 可以同时插入相同行
postgres 使用依赖关系检查冲突，避免重复新增，mysql 使用锁机制
## GitHub Actions
现在的github的仓库名称为main
## Web框架
- Gin
- Fiber

HTTP路由
- FastHttp
- Gorilla Mux
- HttpRouter
- Chi

安装
```bash
go get -u github.com/gin-gonic/gin
```
## 配置
```bash
go get github.com/spf13/viper
```
## 模拟
```bash
go get github.com/golang/mock/mockgen@v1.6.0
# mockgen [包名] [接口]
mockgen --destination db/mock/store.go simplebank-app/db/sqlc Store
```
## 添加Users表
## 哈希
bcrypt 相同加密产生hash会不一致，加入随机salt，但是hash保存有salt
## JWT
新的身份验证令牌： PASETO
Token认证流程
请求->access_token:JWT, PASETO->Authorization:Bear access_token

JWT：HMACSHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload).your-256-bit-secret)

HS256=HMAC + SHA256  
HMAC: Hash-based Message Authentication Code  
SHA: Secure Hash Algorithm  
Asymmetric 非对称签名算法  
私钥进行签名，公钥进行验证  
RS | PS | ES  
JWT 可以篡改加密算法，使用对称加密算法，然后使用公钥签名可以绕过服务器检查
## docker构建
多阶段构建  
```dockerfile
# Build stage
FROM golang:1.18-alpine3.15 AS base
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -o main main.go

# Run stage
FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=base /app/app.env /app/main ./
EXPOSE 9000
CMD ["/app/main"]
```
RUN和CMD区别  
RUN运行在构建时，可以是多条命令，每条命令都会创建一个layer，尽可能放在一条执行  
CMD运行执行时，只能有一条  
COPY和ADD区别  
ADD不仅支持本地文件拷贝，而且支持网络支援拷贝  
## docker网络
`docker inspect [container]` 查看容器信息
 ```bash
 docker network ls 
 docker network inspect bridge
 docker network create bank-network
 docker network connect bank-network postgres
 docker container inspect 
 docker run --name simplebank --network bank-network -e GIN_MODE=release -e DB_SOURCE="postgres://root:123456@postgres:5432/simple_bank?sslmode=disable" -d -p 9000:9000 simplebank
 ```
一个容器可以连接多个网络，加入同一个网络后，可以通过容器名作为主机名进行查找
