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
