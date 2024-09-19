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

---
**中文** | **[English](./README_en.md)**



## 安装

### 1.Docker(推荐)
```bash
mkdir /etc/alist-gallery
vim /etc/alist-gallery/config.yaml  # 配置文件⬇️
docker run -d --restart=unless-stopped -v /etc/alist-gallery/config.yaml:/app/config.yaml -p 5243:5243 --name="alist-gallery" designerwang/alist-gallery:latest
```
### 2.可执行文件
在 [releases](https://github.com/ThinkerWen/alist-gallery/releases) 下载对应平台的可执行文件在本地运行

## 配置文件

```yaml
port: 5243  # alist-gallery 服务端口号
alist-host: https://assets.example.com # alist域名
gallery-location: https://assets.example.com:5243 # alist-gallery服务地址
storage-path: /Storage/Gallery # 图床在alist中的存储路径
alist-token: alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr # alist服务token
password: "" # 存储路径的文夹及密码(可选)
```

## 同步已有数据
创建并修改完config.yaml后，在没生成gallery.db前，执行`sh sync.sh`同步当前`storage-path`下的图片数据到SQLite数据库中，完成同步

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
