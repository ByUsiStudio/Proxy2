<p align="center">
    <a >
        <img src="https://img.shields.io/badge/hp_lite-6.0-red.svg?style=flat" />
    </a>
    <a >
        <img src="https://img.shields.io/badge/Licence-MIT-green.svg?style=flat" />
    </a>
    <a >
        <img src="https://img.shields.io/badge/Licence-AGPL-green.svg?style=flat" />
    </a>

<p align="center">    
    <b>如果对您有帮助，您可以点右上角 "Star" 支持一下 谢谢！</b>
</p>

# Proxy2 6.0内网穿透

#### Proxy2介绍
Proxy2 6.0是一个单机方案
我们采用的是数据转发实现 稳定性可靠性是有保证的即便是极端的环境只要能上网就能实现穿透。
我们是首个支持双通道模式的工具及通道可以选TCP(多路复用)或者QUIC传输，他们都是目前比较高效，高性能的传输方式
我们支持TCP和UDP协议，针对 http/https ws/wss 协议做了大量的优化工作可以更加灵活的控制。让用户使用更佳舒服简单。

### 服务端

#### 二进制文件下载方式
- 下载6.0的二进制文件运行即可
- 配置说明 app.yml
- app.yml文件放在和二进制同目录下即可，不然会默认启动配置
- 建议：部署时关闭所有防火墙，云厂的安全组，注意UDP端口放开，还有TCP，

- app.yml文件
```yaml
admin:
  username: 'admin' #后台账号
  password: '123456' #后台密码
  port: 9090 #管理后台监听的端口（TCP传输方式）

cmd:
  port: 6666 #控制指令端口，所有Proxy2客户端需要连接这个端口（TCP传输方式） 

tunnel:
  ip: '192.168.0.217' #隧道监听服务器外网的IP（记得改成你的服务器IP或者解析的域名也可以）
  port: 9091 #隧道传输数据端口，这个端口用来传输数据的，注意这个是UDP协议，如果是安全组设置记得UDP的放开
  open-domain: false #true 开启80，443端口域名转发（如果你的服务有宝塔或者nginx等，端口多半是被用了），false 关闭

acme:
  email: '23232003@qq.com' #申请证书必须写一个邮箱可以随便写
  http-port: '5634' #证书验证会访问http接口，会通过80转发过来，所以这个端口不用暴露外网

system:
  site-title: 'Proxy2内网穿透' #站点标题，显示在登录页和后台顶部
  open-register: true #是否开启公开注册
  register-review: false #注册后是否需要审核
  smtp:
    enabled: false #是否启用SMTP邮件服务
    host: 'smtp.example.com' #SMTP服务器地址
    port: 587 #SMTP端口
    username: 'email@example.com' #SMTP账号
    password: 'password' #SMTP密码
    from: 'email@example.com' #发件人邮箱
    from-name: 'Proxy2' #发件人名称
    enable-ssl: false #是否启用SSL
```


### 客户端运行方式
##### docker
```shell
#直接连接码
sudo docker run --name hp-lite --restart=always -d  -e c=连接码 registry.cn-shenzhen.aliyuncs.com/heixiaoma/hp-lite:latest
#直接连接码
sudo docker run --name hp-lite --restart=always -d -e  c=连接码 heixiaoma/hp-lite:latest
```

##### Linux或者win
```shell
chmod -R 777 ./hp-lite-amd64
#方式1
./hp-lite-amd64 -server=xxx.com穿透服务:16666 -deviceId=32位的设备ID
#方式2
./hp-lite-amd64 -c=连接码
```