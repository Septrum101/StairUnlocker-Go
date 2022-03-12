## 注意事项

- 本项目使用 [Golang](https://go.dev/) 编写
- 对比之前 [python](https://github.com/thank243/StairUnlocker) 版效率极大提升，实测200+节点仅需5秒
- 测速及解锁测试仅供参考，不代表实际使用情况，由于网络情况变化、Netflix封锁及ip更换，测速具有时效性
- Netflix 解锁测速结果说明:

~~~~text
Full Unlock             全解锁
None                    未解锁
~~~~

## 特性

- 使用 clash-core，原生支持 Shadowsocks，Trojan，V2Ray
- 支持极速测试，可以反映解锁状态，分为 全解锁 / 无解锁 两档
- 支持周期性测试，并将 clash provider 文件上传至 Gist

## 支持平台

### 已测试平台

1. Windows 10 x64
2. Debian 10 x64
3. Ubuntu 20 x64
4. macOS 10.15.7 x64

## 致谢

- [Dreamacro](https://github.com/Dreamacro/clash)

## 使用说明

1. config.yaml 配置文件说明

~~~~yaml
# info / warning / error / debug / silent
log-level: info

# subconverter 服务器地址
converterAPI: https://api.dler.io

# 订阅地址
subURL:

# true：使用本地proxies.yaml文件，导出结果到netflix.yaml
# false：上传到gist
localFile: false

# github token, localFile 为 true 时设置
token:

# 最大同时测试数
maxConn: 32

# Daemon mode 检测周期
internal: 3600

# 启用telegramBot
enableTelegram: false
telegram:
  # bot的token
  telegramToken:
  # telegram用户id
  chatID:
~~~~

2. 运行程序：

- 单次运行:

~~~~bash
./StairUnlocker-Go
~~~~

- 服务模式:

~~~~bash
./StairUnlocker-Go -D
~~~~

4. 命令参数：

~~~~bash
usage: StairUnlocker-Go [-h] [-u SUBURL] [-t TOKEN] [-g GISTURL] [-D]

optional arguments:
    -h    this help
    -v    show current version of StairUnlock
    -u    Load node from subscription url
    -t    The github token
    -g    The gist api URL
    -D    Daemon mode 

~~~~

### 使用 systemd 添加系统服务

StairUnlcoker-Go 可以使用 systemd 在后台运行，以保证当节点不能解锁流媒体时可以重新检测节点并更新 gist

1. 将 StairUnlcoker-Go 的二进制文件拷贝到 `/usr/local/bin/StairUnlocker`:

~~~~bash
# mkdir /usr/local/bin/StairUnlocker/
# cp StairUnlocker-Go /usr/local/bin/
# mkdir /etc/StairUnlocker/
# cp config.yaml /etc/StairUnlocker/
~~~~

2. 在 `/etc/systemd/system/stairunlocker.service` 建立 systemd 配置文件:

~~~~systemd
[Unit]
Description=StairUnlcoker-Go daemon, A fast stream media test tool in Go.
After=network.target

[Service]
Type=simple
Restart=always
ExecStart=/usr/local/bin/StairUnlocker-Go -D -f /etc/StairUnlocker/config.yaml

[Install]
WantedBy=multi-user.target
~~~~

3. 设置系统启动项

~~~~bash
# systemctl enable stairunlocker
~~~~

4. 立即启动 StairUnlcoker-Go:

~~~~bash
# systemctl start stairunlocker
~~~~

5. 检查运行状态和日志:

~~~~bash
# systemctl status stairunlocker
# journalctl -u stairunlocker -f
~~~~
