<div align="center">
  <a href="https://alist.nn.ci"><img width="100px" alt="logo" src="https://cloud.hive-net.cn/gallery-api/fs/show-gallery/2024_09_11_ukNhp1.png"/></a>
  <p><em>📷将 alist 作为图床使用</em></p>
  <a href="https://go.dev/dl/">
    <img src="https://img.shields.io/badge/Go-1.22.1-blue" />
  </a>
  <a href="https://github.com/ThinkerWen/alist-gallery/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/ThinkerWen/alist-gallery" alt="License" />
  </a>
  <a href="https://github.com/ThinkerWen/alist-gallery/releases">
    <img src="https://img.shields.io/github/v/release/ThinkerWen/alist-gallery.svg" alt="latest version" />
  </a>
</div>

<div align="center">
<p><em>Go语言开发，仅10MB极低内存占用，包含用户管理的图床最小可用方案</em></p>
</div>

---
**中文** | **[English](./README_en.md)**



## 安装

### 1.Docker(推荐)
```bash
mkdir /etc/alist-gallery
vim /etc/alist-gallery/config.yaml  # 配置文件⬇️
docker run -d --restart=unless-stopped -v /etc/alist-gallery/config.yaml:/app/config.yaml  -p 5243:5243 --name="alist-gallery" designerwang/alist-gallery:latest
```
### 2.可执行文件
在 [releases](https://github.com/ThinkerWen/alist-gallery/releases) 下载对应平台的可执行文件在本地运行

## 配置文件

```yaml
port: 5243  # alist-gallery 服务端口号
alist-host: https://assets.example.com # alist域名
gallery-location: https://assets.example.com:5243 # alist-gallery服务地址
storage-path: /Storage/Gallery # 图床在alist中的存储路径
alist-token: alist-4254afdc-1acg-1999-08aa-... # alist服务token
password: "" # 存储路径的文夹及密码(可选)
compression: 0 # 是否开启webp图片压缩(加快访问速度，0不开启，0~100)
redis:
  enable: true    # 是否开启Redis缓存
  host: 127.0.0.1 # Redis连接地址
  port: 6379      # Redis端口
  database: 0     # RedisDB序号
  password: ""    # Redis密码
  timeout: 60     # 缓存过期时间(分钟)
```

## 同步已有数据
配置好`config.yaml`后，浏览器加载图片时自动同步数据，无需人工导入

## 扩展

可以配置 nginx upstream 来隐藏端口并达到完美接入alist-api的效果，例如:
```conf
http {
    ...
    
    upstream gallery-api {
        server 127.0.0.1:5243;
    }
    ...
    
    server {
        location /gallery-api/ {
            proxy_pass http://gallery-api/; 
        }
    }
    ...
    
}
```
从而 gallery-location 由 `https://assets.example.com:5243` 变成 `https://assets.example.com/gallery-api`

## 使用
### PUT 表单上传文件
**PUT** `/fs/form-gallery`
> Body 请求参数
```json
{"file": "content"}
```
#### 请求参数
| 名称             | 位置   | 类型   | 必选 | 说明                               |
|----------------|--------|--------|------|----------------------------------|
| Authorization  | header | string | 是   | token                            |
| Content-Type   | header | string | 是   | 需要是multipart/form-data;          |
| File-Name      | header | string | 是   | 文件名(需要保证唯一)                      |
| As-Task        | header | string | 否   | 是否添加为任务                          |
| body           | body   | object | 否   |                                  |
| » file         | body   | string(binary)| 是 | 文件                               |
#### 返回示例
> 成功
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "animated_zoom.gif",
    "url": "https://assets.example.com:5243/fs/show-gallery/animated_zoom.gif",
    "task": {
      "id": "sdH2LbjyWRk",
      "name": "upload animated_zoom.gif to [/data](/alist)",
      "state": 0,
      "status": "uploading",
      "progress": 0,
      "error": ""
    }
  }
}

```

### PUT 流式上传文件
**PUT** `/fs/put-gallery`
> Body 请求参数
```text
string
```
#### 请求参数
| 名称             | 位置   | 类型            | 必选    | 说明          |
|----------------|--------|-----------------|--------|-------------|
| Authorization  | header | string          | 是     | token       |
| Content-Type   | header | string          | 是     |             |
| File-Name      | header | string          | 是     | 文件名(需要保证唯一) |
| As-Task        | header | string          | 否     | 是否添加为任务     |
| body           | body   | string(binary)  | 是     | 文件          |
#### 返回示例
> 成功
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "animated_zoom.gif",
    "url": "https://assets.example.com:5243/fs/show-gallery/animated_zoom.gif",
    "task": {
      "id": "sdH2LbjyWRk",
      "name": "upload animated_zoom.gif to [/data](/alist)",
      "state": 0,
      "status": "uploading",
      "progress": 0,
      "error": ""
    }
  }
}

```

### GET 展示图片
**GET** /fs/show-gallery
> Path 请求参数
#### 请求参数
```url
/fs/show-gallery/animated_zoom.gif
```
#### 返回示例
> 成功

展示图片
