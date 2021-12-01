## 注意事项

- 对比之前 [python](https://github.com/thank243/StairUnlocker) 版效率极大提升，实测200+节点只需要5秒完成检测
- 测速及解锁测试仅供参考，不代表实际使用情况，由于网络情况变化、Netflix封锁及ip更换，测速具有时效性
- 本项目使用 [Golang](https://go.dev/) 编写
- Netflix 解锁测速结果说明:

~~~~text
Full Unlock             全解锁
None                    未解锁
~~~~

## 特性

- 使用 clash-core，原生支持 Shadowsocks，Trojan，V2Ray
- 支持极速测试，可以反映解锁状态
- 支持 Netflix 解锁测试，分为 全解锁 / 无解锁 两档
- 可上传至 Gist

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
#subconverter 服务器地址
converterAPI: https://api.dler.io 
#订阅地址
subURL: 
#true导出到本地文件， false上传到gist
localFile: true
#github token
token: 
#最大同时测试数
maxConn: 32
~~~~

2. 运行程序：
~~~~bash
./main
~~~~

4. 命令参数：

~~~~bash
usage: main [-h] [-u SUBURL] [-t TOKEN] [-g GISTURL]

optional arguments:
	-h    this help
	-v    show current version of StairUnlock
	-u    Load config from subscription url
	-t    The github token
	-g    The gist api URL

~~~~

### 使用actions自动更新

待定