# alist-gallery

将 alist 作为图床使用

## Install

### Docker(推荐)
```bash
mkdir /etc/alist-gallery
vim /etc/alist-gallery/config.yaml  # 配置文件(Configuration)⬇️
docker run -d --restart=unless-stopped -v /etc/alist-gallery/config.yaml:/app/config.yaml -p 5243:5243 --name="alist-gallery" designerwang/alist-gallery:latest
```
### 可执行文件
在 [releases](https://github.com/ThinkerWen/alist-gallery/releases) 下载对应平台的可执行文件在本地运行

## Configuration

```yaml
port: 5243  # alist-gallery 服务端口号
alist-host: https://assets.example.com # alist域名
gallery-location: https://assets.example.com:5243 # alist-gallery服务地址
storage-path: /Storage/Gallery # 图床在alist中的存储地址
alist-token: alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr # alist服务token，可查看图片
password: "" # 存储地址的文件及密码(可选)
```

## Extension

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