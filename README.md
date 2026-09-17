# Plachta

**Plachta** 是一个面向 Linux VPS 的代理基础设施管理工具，提供代理核心安装、Reality 节点配置、配置导出以及基础运行状态管理。

当前版本：**0.1.0**

> 当前版本主要围绕 Xray Reality 节点进行管理。

---

## 功能

当前已实现：

* [x] Go CLI
* [x] Linux x86_64 / amd64 支持
* [x] Linux arm64 / aarch64 支持
* [x] 一键安装
* [x] Xray Reality 安装
* [x] Reality 配置自动生成
* [x] VLESS Reality URI 生成
* [x] Quantumult X 配置导出
* [x] 服务状态检查
* [x] 系统信息查看
* [x] 配置验证
* [x] 配置备份
* [x] 模块信息查看
* [ ] Firewall 自动管理
* [ ] Firewall 白名单自动维护
* [ ] 重启后 Firewall 自动恢复
* [ ] 更多代理核心与客户端支持

未完成的功能不会影响当前 Reality 节点的正常使用。

---

## 系统要求

### 支持的系统

目前主要测试环境：

* Debian 12
* Linux x86_64 / amd64
* Linux arm64 / aarch64

理论上支持其他 systemd Linux 发行版，但未全部测试。

### 权限

安装 Plachta 和 Xray Reality 通常需要 `root` 权限。

建议直接使用 root 用户执行：

```bash
sudo -i
```

或者：

```bash
su -
```

---

## 一键安装

在 Linux VPS 上执行：

```bash
curl -fsSL https://raw.githubusercontent.com/KirisameMarishark/Plachta/main/install.sh | bash
```

安装程序会自动：

1. 检测系统架构
2. 下载对应的 Plachta Linux 二进制
3. 安装到系统路径
4. 设置可执行权限

安装完成后检查：

```bash
plachta version
```

预期输出：

```text
Plachta 0.1.0
```

---

## 基本命令

查看帮助：

```bash
plachta help
```

查看版本：

```bash
plachta version
```

查看系统信息：

```bash
plachta system
```

查看当前服务状态：

```bash
plachta status
```

例如：

```text
Plachta Status
------------------------------
xray          : active
```

---

## 安装 Reality

执行：

```bash
plachta install reality
```

Plachta 会自动完成 Reality 的基础配置，包括：

* 安装 Xray
* 生成 Reality 配置
* 生成 UUID
* 生成 Reality 密钥
* 生成 Short ID
* 配置 Reality SNI
* 配置 Xray 服务
* 重启 Xray

安装完成后会直接输出 VLESS Reality URI。

示例：

```text
vless://UUID@SERVER:443?type=tcp&security=reality&pbk=PUBLIC_KEY&sid=SHORT_ID&sni=www.cloudflare.com&fp=chrome&flow=xtls-rprx-vision&encryption=none#Plachta-Reality
```

> 示例中的 UUID、Public Key、Short ID 和服务器地址仅用于说明格式，实际配置由 Plachta 自动生成。

---

## Quantumult X 导出

Reality 安装完成后，可以直接导出 Quantumult X 配置：

```bash
plachta export quantumultx
```

Plachta 会根据当前 Reality 配置生成 Quantumult X 可使用的节点配置。

示例格式：

```text
vless=SERVER:443, method=none, password=UUID, obfs=over-tls, obfs-host=SNI, reality-base64-pubkey=PUBLIC_KEY, reality-hex-shortid=SHORT_ID, vless-flow=xtls-rprx-vision, tag=Plachta-Reality
```

将输出内容复制到 Quantumult X 对应配置位置即可。

---

## 查看 Reality 配置

可以使用：

```bash
plachta show reality
```

查看当前 Reality 相关信息。

---

## 配置验证

执行：

```bash
plachta verify
```

用于检查当前 Plachta 配置状态。

---

## 配置备份

执行：

```bash
plachta backup
```

用于备份当前配置。

---

## 模块

查看当前模块：

```bash
plachta modules
```

查看指定模块信息：

```bash
plachta module-info reality
```

---

## 更新

Plachta 使用 GitHub Release 发布 Linux 二进制。

重新执行一键安装命令即可获取当前 Release 中对应架构的版本：

```bash
curl -fsSL https://raw.githubusercontent.com/KirisameMarishark/Plachta/main/install.sh | bash
```

安装后检查：

```bash
plachta version
```

---

## 当前 Reality 使用流程

一个全新的 Debian VPS 可以按照以下步骤使用：

### 1. 安装 Plachta

```bash
curl -fsSL https://raw.githubusercontent.com/KirisameMarishark/Plachta/main/install.sh | bash
```

### 2. 检查版本

```bash
plachta version
```

### 3. 安装 Reality

```bash
plachta install reality
```

### 4. 检查 Xray

```bash
plachta status
```

如果显示：

```text
xray          : active
```

说明 Xray 服务已经正常运行。

### 5. 导出 Quantumult X

```bash
plachta export quantumultx
```

---

## 手动检查 Xray

Plachta 安装 Reality 后，Xray 使用 systemd 管理。

查看状态：

```bash
systemctl status xray --no-pager -l
```

重启：

```bash
systemctl restart xray
```

检查是否运行：

```bash
systemctl is-active xray
```

预期：

```text
active
```

---

## 项目结构

项目主要结构：

```text
Plachta/
├── cmd/
│   └── plachta-go/
│       └── main.go
│
├── internal/
│   └── core/
│       ├── cli/
│       ├── config/
│       ├── firewall/
│       ├── module/
│       ├── reality/
│       ├── subscription/
│       ├── system/
│       └── version/
│
├── install.sh
├── VERSION
├── README.md
└── go.mod
```

---

## Reality 配置

当前 Reality 模块负责生成和管理基础 Xray Reality 配置。

典型 Reality 参数包括：

| 参数          | 说明                  |
| ----------- | ------------------- |
| UUID        | VLESS 用户 UUID       |
| Public Key  | Reality 公钥          |
| Private Key | Reality 私钥          |
| Short ID    | Reality Short ID    |
| SNI         | Reality Server Name |
| Fingerprint | TLS Fingerprint     |
| Flow        | VLESS XTLS Vision   |
| Port        | Reality 服务端口        |

默认使用 TCP +
