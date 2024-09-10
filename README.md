# alist-gallery
**中文** | **[English](https://github.com/ThinkerWen/alist-gallery/blob/main/README_en.md)**

将 alist 作为图床使用

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
alist-token: alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr # alist服务token(可查看图片)
password: "" # 存储路径的文夹及密码(可选)
```

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
file: []
```
#### 请求参数
| 名称             | 位置   | 类型   | 必选 | 说明                               |
|----------------|--------|--------|------|----------------------------------|
| Authorization  | header | string | 是   |    token                              |
| Content-Type   | header | string | 是   | 需要是multipart/form-data;          |
| Content-Length | header | string | 是   | 文件大小                             |
| File-Path      | header | string | 是   | 完整文件路径(与File-Name可选，优先File-Path) |
| File-Name      | header | string | 是   | 文件名(与File-Path可选，优先File-Path)    |
| As-Task        | header | string | 否   | 是否添加为任务                          |
| body           | body   | object | 否   |                              |
| » file         | body   | string(binary)| 是 | 文件                               |
#### 返回示例
> 成功
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "name": "animated_zoom.gif",
    "path": "/Storage/Gallery/animated_zoom.gif",
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