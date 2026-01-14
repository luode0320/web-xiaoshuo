# 小说阅读系统部署文档

## 部署方式概述

本系统支持多种部署方式，包括传统部署和容器化部署。

## Docker容器化部署（推荐）

### 环境要求

- Docker 20.10+
- Docker Compose v2.0+

### 部署步骤

1. **克隆项目代码**
   ```bash
   git clone https://github.com/your-username/web-xiaoshuo.git
   cd web-xiaoshuo
   ```

2. **构建并启动服务**
   ```bash
   docker-compose up -d
   ```

3. **初始化数据库**
   ```bash
   # 进入后端容器执行数据库迁移
   docker exec -it xiaoshuo-backend go run main.go
   # 或者执行数据库迁移脚本
   docker exec -it xiaoshuo-backend go run migrations/migrate.go
   ```

4. **访问系统**
   - 前端访问地址：http://localhost:3000
   - 后端API地址：http://localhost:8888/api/v1

### Docker Compose 服务说明

- `mysql`: MySQL数据库服务（端口3306）
- `redis`: Redis缓存服务（端口6379）
- `xiaoshuo-backend`: 后端API服务（端口8888）
- `xiaoshuo-frontend`: 前端Web服务（端口3000）

## 传统部署方式

### 后端部署

1. **安装Go环境**
   - 安装Go 1.24+
   - 配置GOPATH

2. **配置数据库**
   - 安装MySQL 5.7+
   - 创建数据库和用户
   - 执行数据库迁移脚本

3. **配置Redis**
   - 安装Redis
   - 确保Redis服务正常运行

4. **编译并运行**
   ```bash
   cd xiaoshuo-backend
   go mod tidy
   go build -o server .
   ./server
   ```

### 前端部署

1. **安装Node.js环境**
   - 安装Node.js 18+
   - 安装npm

2. **编译并部署**
   ```bash
   cd xiaoshuo-frontend
   npm install
   npm run build
   ```

3. **使用Web服务器部署**
   - 将`dist`目录内容部署到Nginx/Apache等Web服务器

## 生产环境配置

### 配置文件调整

修改 `xiaoshuo-backend/config/config.yaml`：

```yaml
server:
  port: "8888"
  mode: "release"  # 生产环境使用release模式
  base_path: "/api/v1"

database:
  host: "mysql"    # 如果使用Docker，使用服务名
  port: "3306"
  user: "root"
  password: "Ld588588"
  dbname: "xiaoshuo"
  charset: "utf8mb4"

redis:
  addr: "redis:6379"  # 如果使用Docker，使用服务名
  password: ""
  db: 0

jwt:
  secret: "your-production-secret-key"  # 生产环境请使用强密钥
  expires: 3600
```

### 环境变量

推荐使用环境变量覆盖配置：

- `DB_HOST`: 数据库主机
- `DB_PORT`: 数据库端口
- `DB_USER`: 数据库用户名
- `DB_PASSWORD`: 数据库密码
- `DB_NAME`: 数据库名
- `REDIS_ADDR`: Redis地址
- `SERVER_PORT`: 服务端口
- `JWT_SECRET`: JWT密钥

## 监控和日志

### 日志配置

系统使用Zap日志库，日志输出到标准输出，Docker环境下可通过以下命令查看：

```bash
# 查看后端日志
docker logs xiaoshuo-backend

# 查看前端日志
docker logs xiaoshuo-frontend

# 实时查看日志
docker logs -f xiaoshuo-backend
```

### 健康检查

- 后端健康检查：`GET /api/v1/health`
- 前端可通过访问根路径检查服务状态

## 性能优化

### 数据库优化

- 启用MySQL慢查询日志
- 定期优化表结构
- 配置合适的索引

### Redis优化

- 配置Redis持久化策略
- 设置合适的内存淘汰策略
- 监控Redis内存使用情况

### 前端优化

- 启用Nginx Gzip压缩
- 配置静态资源缓存策略
- 使用CDN加速静态资源访问

## 安全配置

### HTTPS配置

在生产环境中建议配置HTTPS：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;

    # 前端配置
    location / {
        root /usr/share/nginx/html;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api {
        proxy_pass http://xiaoshuo-backend:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 访问控制

- 限制API访问频率
- 配置防火墙规则
- 定期更新依赖包

## 备份和恢复

### 数据库备份

```bash
# 备份数据库
docker exec xiaoshuo-mysql mysqldump -u root -pLd588588 xiaoshuo > backup.sql

# 恢复数据库
docker exec -i xiaoshuo-mysql mysql -u root -pLd588588 xiaoshuo < backup.sql
```

### 定期备份

建议配置定时任务进行定期备份：

```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec xiaoshuo-mysql mysqldump -u root -pLd588588 xiaoshuo > /backup/xiaoshuo_${DATE}.sql
# 删除7天前的备份
find /backup -name "xiaoshuo_*.sql" -mtime +7 -delete
```

## 故障排除

### 常见问题

1. **服务无法启动**
   - 检查端口是否被占用
   - 检查数据库连接是否正常
   - 查看日志输出获取错误信息

2. **数据库连接失败**
   - 检查数据库服务是否正常运行
   - 检查数据库配置是否正确
   - 检查网络连接是否正常

3. **Redis连接失败**
   - 检查Redis服务是否正常运行
   - 检查Redis配置是否正确

### 日志分析

- 后端错误日志：关注5xx错误和异常堆栈
- 前端错误日志：关注API调用失败和JavaScript错误
- 数据库慢查询日志：优化慢查询

## 扩展和维护

### 水平扩展

- 使用负载均衡器部署多个前端实例
- 数据库读写分离
- Redis集群配置

### 版本升级

1. 备份当前系统
2. 拉取最新代码
3. 执行数据库迁移
4. 重新构建镜像
5. 重启服务

```bash
# 更新代码
git pull origin main

# 重新构建并启动服务
docker-compose down
docker-compose build
docker-compose up -d
```

## 联系和支持

如需技术支持或遇到部署问题，请联系开发团队。