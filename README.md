# KWCLI

KWCLI 是 KWDB 生态的命令行工具，采用组件化架构设计，帮助你一键安装、运行和管理 KWDB 相关组件。

**开发语言**：Go (v1.21+)

## 特性

- 🧩 **组件化管理**：每个功能都是独立组件，按需安装、灵活扩展
- 🚀 **一键启动**：`kwcli playground start` 秒级启动 Playground
- 🐳 **开箱即用**：自动处理依赖与环境检查
- 💻 **跨平台**：基于 Go 构建，支持 Linux / macOS / Windows
- 🌐 **国内加速**：默认优先阿里云镜像仓库，代码源默认 AtomGit
- ⚙️ **全局配置**：支持一键切换默认代码源和镜像源
- 📊 **SampleDB**：内置智能电表模型，一键初始化 schema、生成数据、运行场景查询

## 安装

### 使用 Makefile（推荐）

```bash
# 克隆仓库
git clone https://github.com/shawn0915/kwcli.git
cd kwcli

# 构建（自动注入编译时间）
make build

# 安装到系统
make install

# 或安装到用户目录
make install-local
```

### 手动构建

```bash
git clone https://github.com/shawn0915/kwcli.git
cd kwcli
go build -ldflags "-X 'github.com/shawn0915/kwcli/cmd.buildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'github.com/shawn0915/kwcli/cmd.commitHash=$(git rev-parse --short HEAD)'" -o kwcli .
sudo mv kwcli /usr/local/bin/
```

## 快速开始

### 全局配置

```bash
# 查看当前默认代码源（默认 atomgit）
kwcli source

# 切换到 GitHub
kwcli source github
```

### 启动 Playground

```bash
# 一键启动 KWDB Playground
kwcli playground start

# 指定版本安装
kwcli playground install -v v1.1.0

# 升级 Playground（代码源和镜像源可分别指定）
kwcli playground upgrade --source atomgit --registry auto
```

启动完成后，打开浏览器访问 http://localhost:3006 即可开始学习。

### 安装 KWDB

```bash
# 安装 KWDB（自动选择 Docker 方式）
kwcli kwdb install

# 启动 KWDB
kwcli kwdb start
```

### 体验 SampleDB（智能电表模型）

```bash
# 初始化数据库和表结构
kwcli sampledb init

# 生成示例数据（100 个电表，10000 条读数）
kwcli sampledb generate

# 查看所有场景
kwcli sampledb list

# 运行指定场景
kwcli sampledb run top10-area-energy
kwcli sampledb run fault-meters

# 运行所有场景
kwcli sampledb run --all

# 清理所有 SampleDB 数据
kwcli sampledb clean
```

## 命令手册

### 全局配置

| 命令 | 说明 |
|------|------|
| `kwcli source` | 查看当前默认代码源 |
| `kwcli source github` | 设置默认代码源为 GitHub |
| `kwcli source atomgit` | 设置默认代码源为 AtomGit |

### Playground 相关

| 命令 | 说明 |
|------|------|
| `kwcli playground install` | 安装 Playground（支持 `-v` 指定版本，`--source` 代码源） |
| `kwcli playground start` | 启动 Playground 服务（首次会自动下载） |
| `kwcli playground stop` | 停止 Playground 服务 |
| `kwcli playground status` | 查看 Playground 运行状态 |
| `kwcli playground restart` | 重启 Playground 服务 |
| `kwcli playground upgrade` | 升级 Playground（支持 `--source` 代码源，`--registry` 镜像源） |
| `kwcli playground logs` | 查看 Playground 日志 |
| `kwcli playground versions` | 查看已安装的 Playground 版本 |
| `kwcli playground uninstall` | 卸载 Playground（加 `--remove-files` 彻底删除） |

### SampleDB 相关

| 命令 | 说明 |
|------|------|
| `kwcli sampledb init` | 创建智能电表数据库和表结构 |
| `kwcli sampledb generate` | 生成示例数据 |
| `kwcli sampledb list` | 列出所有场景查询 |
| `kwcli sampledb run <name>` | 运行指定场景查询 |
| `kwcli sampledb run --all` | 运行所有场景查询 |
| `kwcli sampledb clean` | 清理所有 SampleDB 数据 |
| `kwcli sampledb status` | 检查 SampleDB 是否存在 |

### SQL 连接

| 命令 | 说明 |
|------|------|
| `kwcli sql` | 交互式连接 KWDB |
| `kwcli sql -e "SELECT 1"` | 执行单条 SQL |
| `kwcli sql -u <user> -d <db>` | 指定用户和数据库 |

### KWDB 相关

| 命令 | 说明 |
|------|------|
| `kwcli kwdb install` | 安装 KWDB（Docker 方式） |
| `kwcli kwdb start` | 启动 KWDB 服务 |
| `kwcli kwdb stop` | 停止 KWDB 服务 |
| `kwcli kwdb restart` | 重启 KWDB 服务 |
| `kwcli kwdb status` | 查看 KWDB 状态 |
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

## 架构简介

KWCLI 采用组件化架构设计：

1. **组件隔离**：所有组件安装在 `~/.kwcli/components/` 目录下，不影响系统环境
2. **按需获取**：组件在首次使用时自动从远程仓库拉取，无需预装
3. **统一入口**：通过 `kwcli <component>` 统一调度，简化操作
4. **多版本共存**：支持同一组件多个版本并存，灵活切换

```
~/.kwcli/
├── components/          # 组件安装目录
│   ├── playground/
│   │   └── versions/    # 多版本共存
│   └── kwdb/            # KWDB 安装目录
├── data/                # 运行时数据
└── config.yaml          # 全局配置（代码源等）
```

## 依赖要求

- [Go](https://golang.org/) 1.21+（仅构建时需要）
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose（Playground 必需）
- Git（用于克隆组件仓库）

## Roadmap

- [x] Playground 一键启动 (`kwcli playground start`)
- [x] KWDB 单机版安装 (`kwcli kwdb install`)
  - 支持 Docker 部署
  - 二进制下载接口已预留（待 Gitee API 集成）
- [x] KWDB 服务管理 (`kwcli kwdb start/stop/restart/status/logs`)
- [x] 配置管理 (`kwcli kwdb config show/edit/set/path`)
- [x] SQL 连接 (`kwcli sql`)
- [x] Playground 版本管理 (`kwcli playground install/versions/uninstall`)
- [x] 全局代码源配置 (`kwcli source`)
- [x] 阿里云镜像加速 (`--registry auto`)
- [x] Makefile 构建脚本
- [x] SampleDB 智能电表模型 (`kwcli sampledb`)
- [ ] 组件清单与版本索引
- [ ] 离线镜像与私有化部署支持
- [ ] Homebrew / install.sh 一键安装

## 作者

**Shawn Yan**

- 公众号「少安事务所」主笔
- KWDB MVP
- 个人主页：[shawnyan.cn](https://shawnyan.cn)

## 贡献

欢迎提交 Issue 和 PR！

## License

Apache License 2.0
