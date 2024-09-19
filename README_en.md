<div align="center">
  <a href="https://alist.nn.ci"><img width="100px" alt="logo" src="https://cloud.hive-net.cn/gallery-api/fs/show-gallery/2024_09_11_ukNhp1.png"/></a>
  <p><em>📷Use alist as a graph bed</em></p>
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
<p><em>Developed in Go language, it has a very low memory of only 10MB and includes the minimum available solution for user-managed graph beds</em></p>
</div>

---
**English** | **[中文](https://github.com/ThinkerWen/alist-gallery/blob/main/README.md)**

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
alist-token: alist-4254afdc-1acg-1999-08aa-b6134kx4kv63FdkHJFPeaFDdEGYmSe29KETy4fdsareKM8fdsagfdsgfdgfdagdfgr # alist service token
password: "" # Folder password for storage path (optional)
```

## Synchronize existing data
After creating and modifying the config.yaml, run the `sh sync.sh` command to synchronize the image data in the current storage-path to the SQLite database before the gallery.db is generated

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
{"file": "content"}
```
#### Request parameters
| Name           | location | type           | must-have | introduce                             |
|----------------|----------|----------------|-----------|---------------------------------------|
| Authorization  | header   | string         | 是         | Token                                 |
| Content-Type   | header   | string         | 是         | Must be multipart/form-data;          |
| File-Name      | header   | string         | 是         | File name (guarantee unique required) |
| As-Task        | header   | string         | 否         | Whether to add it as a task           |
| body           | body     | object         | 否         |                                       |
| » file         | body     | string(binary) | 是         | File                                  |
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

### PUT stream upload file
**PUT** `/fs/put-gallery`
> Body request parameters
```text
string
```
#### Request parameters
| Name           | location | type           | must-have | introduce                             |
|----------------|----------|----------------|-----------|---------------------------------------|
| Authorization  | header   | string         | 是         | Token                                 |
| Content-Type   | header   | string         | 是         |                                       |
| File-Name      | header   | string         | 是         | File name (guarantee unique required) |
| As-Task        | header   | string         | 否         | Whether to add it as a task           |
| body           | body     | string(binary) | 否         | File                                  |
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
