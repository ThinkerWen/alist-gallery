# alist-gallery
**English** | **[中文](https://github.com/ThinkerWen/alist-gallery/blob/main/README.md)**

Use alist as a graph bed

## Install

### 1.Docker(recommend)
```bash
mkdir /etc/alist-gallery
vim /etc/alist-gallery/config.yaml  # Configuration⬇️
docker run -d --restart=unless-stopped -v /etc/alist-gallery/config.yaml:/app/config.yaml -p 5243:5243 --name="alist-gallery" designerwang/alist-gallery:latest
```
### 2.Executable file
Download the executable file of the corresponding platform from [releases](https://github.com/ThinkerWen/alist-gallery/releases) and run it locally

## Configuration

```yaml
port: 5243  # alist-gallery service port
alist-host: https://assets.example.com # alist domain
gallery-location: https://assets.example.com:5243 # alist-gallery service location
storage-path: /Storage/Gallery # The path where the graph bed is stored in the alist
alist-token: alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr # alist service token (can view image)
password: "" # Folder password for storage path (optional)
```

## Extension

nginx upstream can be configured to hide the port and achieve perfect access to the alist-api, for example:
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
Thus the gallery-location changes from `https://assets.example.com:5243` to `https://assets.example.com/gallery-api`

## Usage
### PUT form upload file
**PUT** `/fs/form-gallery`
> Body request parameters
```json
file: []
```
#### Request parameters
| Name           | location | type           | must-have | introduce                                                 |
|----------------|----------|----------------|-----------|-----------------------------------------------------------|
| Authorization  | header   | string         | 是         | Token                                                     |
| Content-Type   | header   | string         | 是         | Must be multipart/form-data;                              |
| Content-Length | header   | string         | 是         | File size                                                 |
| File-Path      | header   | string         | 是         | Full file path (optional with File-Name, File-Path first) |
| File-Name      | header   | string         | 是         | File name (optionally with File-Path, File-Path first)    |
| As-Task        | header   | string         | 否         | Whether to add it as a task                               |
| body           | body     | object         | 否         |                                                           |
| » file         | body     | string(binary) | 是         | File                                                      |
#### Response example
> Success
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

### GET display image
**GET** /fs/show-gallery
> Path request parameters
#### Request parameters
```url
/fs/show-gallery/animated_zoom.gif
```
#### Response example
> Success

Display image