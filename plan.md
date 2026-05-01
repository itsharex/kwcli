# KWCLI 项目计划

## 1. 项目目标

KWCLI 是 KWDB 生态的命令行工具，采用组件化设计理念，提供统一的组件管理、本地安装和运维能力。

**开发语言**：Go (v1.21+)
**CLI 框架**：Cobra + Viper

第一阶段目标：
- 实现 `kwcli playground start` 命令，一键启动 KWDB Playground 交互式学习平台
- 支持本地安装 KWDB 单机版
- 建立可扩展的组件管理框架

## 2. 设计理念

| 设计原则 | 说明 |
|---------|------|
| **组件化架构** | 每个功能（playground、kwdb、cluster 等）都是独立组件，可独立安装、升级、卸载 |
| **按需获取** | 组件首次使用时自动从远程仓库下载/克隆，无需预先全量安装 |
| **本地隔离** | 所有组件安装在 `~/.kwcli/components/` 下，与系统环境隔离 |
| **多版本共存** | 支持同一组件多个版本并存，可灵活切换 |
| **统一入口** | 通过 `kwcli <component>` 统一调用，降低用户学习成本 |
| **镜像加速** | Docker 镜像默认优先阿里云镜像仓库，提升国内访问速度 |

## 3. 架构设计

```
kwcli/
├── main.go                  # 程序入口
├── Makefile                 # 构建脚本
├── go.mod                   # Go 模块定义
├── cmd/                     # CLI 命令入口
│   ├── root.go              # 根命令（全局配置、版本信息）
│   ├── playground.go        # playground 子命令
│   ├── kwdb.go              # kwdb 子命令
│   ├── sql.go               # SQL 连接命令
│   └── general.go           # 通用命令 (list, install, update, uninstall, source)
└── pkg/                     # 核心包
    ├── component/
    │   └── manager.go       # 组件生命周期管理
    ├── config/
    │   └── config.go        # 配置管理 (Viper)
    └── utils/
        ├── docker.go        # Docker 工具
        └── http.go          # HTTP 下载工具
```

### 目录约定

```
~/.kwcli/
├── bin/                     # 组件二进制软链接（未来支持）
├── components/              # 组件安装目录
│   ├── playground/          # Playground 子项目
│   │   └── versions/        # 多版本共存目录
│   └── kwdb/                # KWDB 安装目录
├── data/                    # 运行数据
└── config.yaml              # 全局配置
```

## 4. 功能规划

### Phase 1 - Playground 支持 ✅ 已完成

- [x] 项目初始化与 CLI 框架搭建（Go + Cobra）
- [x] `kwcli playground start` 命令实现
- [x] `kwcli playground install` 安装指定版本
- [x] `kwcli playground versions` 查看已安装版本
- [x] `kwcli playground uninstall` 卸载（支持 `--remove-files`）
- [x] 自动克隆 KWDB/playground 子项目
- [x] 通过 Docker Compose 启动 Playground 服务
- [x] 服务健康检查与端口冲突检测
- [x] `kwcli playground stop` 停止服务
- [x] `kwcli playground status` 查看状态
- [x] `kwcli playground restart` 重启服务
- [x] `kwcli playground upgrade` 升级（支持 `--source` 代码源、`--registry` 镜像源）
- [x] `kwcli playground logs` 查看日志
- [x] 版本安装支持（`-v v1.1.0`）
- [x] 代码源切换（GitHub / AtomGit / auto）
- [x] Docker 镜像源切换（Docker Hub / 阿里云 / auto）

### Phase 2 - KWDB 单机版安装 ✅ 已完成

- [x] `kwcli kwdb install` 命令
- [x] 二进制下载支持（预留接口，Gitee API）
- [x] 无匹配包时自动回退到 Docker
- [x] `kwcli kwdb start` 本地启动单机版 KWDB
- [x] `kwcli kwdb stop` 停止 KWDB 服务
- [x] `kwcli kwdb restart` 重启 KWDB 服务
- [x] `kwcli kwdb status` 查看 KWDB 状态（精确匹配，不混入其他容器）
- [x] `kwcli kwdb logs` 查看 KWDB 日志
- [x] 配置文件生成与管理
  - 默认配置模板
  - 交互式配置向导 (`kwcli kwdb config edit`)
  - 查看配置 (`kwcli kwdb config show`)
  - 设置配置项 (`kwcli kwdb config set <key> <value>`)
  - 从配置启动 KWDB
- [x] SQL 连接 (`kwcli sql`)
  - 交互式 SQL Shell
  - 单条 SQL 执行 (`-e` 参数)
  - 支持 Docker 和 Binary 模式

### Phase 2.5 - 体验优化 ✅ 已完成

- [x] 全局默认代码源配置（`kwcli source`，默认 atomgit）
- [x] 所有 install / update 命令支持 `--source` flag
- [x] 构建时注入编译时间（`kwcli -v` 显示版本+编译时间）
- [x] Makefile 构建脚本
- [x] 阿里云镜像仓库优先（`registry.cn-hangzhou.aliyuncs.com/kwdb`）
- [x] Docker 镜像拉取失败自动回退机制

### Phase 3 - 组件仓库与版本管理（待开发）

- [ ] 组件清单与版本索引
- [ ] 离线镜像支持
- [ ] 组件升级与回滚

## 5. 命令手册

### 全局配置

| 命令 | 说明 |
|------|------|
| `kwcli source` | 查看当前默认代码源 |
| `kwcli source github` | 设置默认代码源为 GitHub |
| `kwcli source atomgit` | 设置默认代码源为 AtomGit |

### SQL 连接

| 命令 | 说明 |
|------|------|
| `kwcli sql` | 交互式连接 KWDB |
| `kwcli sql -e "SELECT 1"` | 执行单条 SQL 并退出 |
| `kwcli sql -u <user>` | 指定用户名 |
| `kwcli sql -d <database>` | 指定数据库 |
| `kwcli sql --host <host>` | 指定连接主机 |

### Playground 相关

| 命令 | 说明 |
|------|------|
| `kwcli playground install` | 安装 Playground（支持 `-v` 指定版本，`--source` 指定代码源） |
| `kwcli playground start` | 启动 Playground 服务（首次会自动下载） |
| `kwcli playground stop` | 停止 Playground 服务 |
| `kwcli playground status` | 查看 Playground 运行状态 |
| `kwcli playground restart` | 重启 Playground 服务 |
| `kwcli playground upgrade` | 升级 Playground（支持 `--source` 代码源，`--registry` 镜像源） |
| `kwcli playground logs` | 查看 Playground 日志 |
| `kwcli playground versions` | 查看已安装的 Playground 版本 |
| `kwcli playground uninstall` | 卸载 Playground（加 `--remove-files` 彻底删除） |

### KWDB 相关

| 命令 | 说明 |
|------|------|
| `kwcli kwdb install` | 安装 KWDB（自动选择 binary 或 Docker） |
| `kwcli kwdb start` | 启动 KWDB 服务 |
| `kwcli kwdb stop` | 停止 KWDB 服务 |
| `kwcli kwdb restart` | 重启 KWDB 服务 |
| `kwcli kwdb status` | 查看 KWDB 运行状态 |
| `kwcli kwdb logs` | 查看 KWDB 日志 |
| `kwcli kwdb config show` | 查看当前配置 |
| `kwcli kwdb config edit` | 交互式编辑配置 |
| `kwcli kwdb config set <key> <value>` | 设置配置项 |
| `kwcli kwdb config path` | 显示配置文件路径 |

### 通用命令

| 命令 | 说明 |
|------|------|
| `kwcli list` | 查看可用组件列表 |
| `kwcli install <component>` | 安装指定组件（支持 `--source`） |
| `kwcli update <component>` | 升级指定组件（支持 `--source`） |
| `kwcli uninstall <component>` | 卸载指定组件 |
| `kwcli status` | 查看所有运行中的组件 |

## 6. 技术选型

| 项目 | 选型 | 理由 |
|-----|------|------|
| 开发语言 | Go 1.21+ | 跨平台编译，单二进制分发 |
| CLI 框架 | Cobra | Go 生态最成熟的 CLI 框架 |
| 配置管理 | Viper | 支持多格式、环境变量覆盖 |
| Playground 部署 | Docker Compose | playground 原生支持 |
| KWDB 安装 | Binary + Docker | 优先二进制，无匹配包时回退 Docker |
| 镜像仓库 | 阿里云 + Docker Hub | 国内用户优先阿里云，失败自动回退 |

## 7. 构建与安装

```bash
# 构建（自动注入编译时间）
make build

# 运行测试
make test

# 安装到系统
make install

# 或安装到用户目录
make install-local

# 查看所有可用命令
make help
```

### 手动构建

```bash
# 构建
make build

# 或手动 go build
go build -ldflags "-X 'github.com/KWDB/kwcli/cmd.buildTime=$(date '+%Y-%m-%d %H:%M:%S')'" -o kwcli .

# 安装到系统
sudo mv kwcli /usr/local/bin/

# 或添加到 PATH
export PATH=$PATH:/path/to/kwcli
```

## 8. 关键设计决策

### 8.1 Playground 作为组件还是内置子模块？

**决策：作为可管理组件，运行时克隆。**

理由：
- 采用组件化设计，playground 与 kwcli 解耦
- playground 可独立更新，无需重新编译 kwcli
- 减小 kwcli 二进制体积

### 8.2 KWDB 安装策略

**决策：优先二进制下载，失败则回退 Docker。**

理由：
- 二进制运行效率更高，无容器开销
- Docker 作为兜底方案，确保在任何系统上都能运行
- 自动检测系统架构，选择合适的安装包

### 8.3 配置管理

**决策：YAML 配置文件 + 交互式向导。**

理由：
- YAML 格式易于阅读和编辑
- 交互式向导降低配置门槛
- 支持通过命令行参数覆盖

### 8.4 全局默认代码源

**决策：默认 AtomGit，可一键切换。**

理由：
- 国内网络环境下 AtomGit 访问更稳定
- 用户可通过 `kwcli source github` 一键切换
- 所有命令的 `--source auto` 默认读取全局配置

### 8.5 Docker 镜像源策略

**决策：默认优先阿里云镜像，失败回退 Docker Hub。**

理由：
- 国内拉取 Docker Hub 镜像经常超时或失败
- 阿里云镜像仓库 `registry.cn-hangzhou.aliyuncs.com/kwdb` 提供加速
- auto 模式下自动尝试阿里云，失败后无感回退 Docker Hub
