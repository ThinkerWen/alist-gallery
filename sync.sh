#!/bin/bash

# 数据库路径
DB_PATH="./gallery.db"

# YAML配置文件路径
YAML_PATH="./config.yaml"

# 创建表的SQL语句
CREATE_TABLE_SQL="CREATE TABLE IF NOT EXISTS gallery_index (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  path VARCHAR(255) NOT NULL DEFAULT '',
  user VARCHAR(255) NOT NULL DEFAULT '',
  image_name VARCHAR(255) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);"

# 检查表是否存在
TABLE_EXISTS=$(sqlite3 $DB_PATH "SELECT name FROM sqlite_master WHERE type='table' AND name='gallery_index';")

if [[ -n $TABLE_EXISTS ]]; then
    echo "表已存在，停止插入"
    exit 0
fi

# 创建表
sqlite3 $DB_PATH "$CREATE_TABLE_SQL"

# 从YAML文件中提取配置信息
alist_token=$(grep 'alist-token' $YAML_PATH | awk '{print $2}')
alist_host=$(grep 'alist-host' $YAML_PATH | awk '{print $2}')
storage_path=$(grep 'storage-path' $YAML_PATH | awk '{print $2}')
password=$(grep 'password' $YAML_PATH | awk '{print $2}')

if [[ "$password" == \"\" ]]; then
  password=""
fi

# 设置请求URL和headers
url="${alist_host}/api/fs/list"
headers="Authorization: ${alist_token}"

# 分页获取数据并插入到SQLite数据库
page=1
per_page=50
while true; do
    # 请求API
    response=$(curl -s -X POST "$url" \
        -H "$headers" \
        -H "Content-Type: application/json" \
        -d "{\"path\": \"$storage_path\", \"password\": \"$password\", \"page\": $page, \"per_page\": $per_page, \"refresh\": false}")

    # 解析响应数据
    code=$(echo "$response" | jq -r '.code')
    if [[ "$code" != "200" ]]; then
        echo "alist server not running"
        exit 1
    fi

    # 获取文件内容列表
    content=$(echo "$response" | jq -c '.data.content[]')
    if [[ -z "$content" ]]; then
        break
    fi

    # 遍历每个项目并插入数据库
    for item in $content; do
        is_dir=$(echo "$item" | jq -r '.is_dir')
        name=$(echo "$item" | jq -r '.name')

        # 如果不是目录，插入到数据库
        if [[ "$is_dir" == "false" ]]; then
            sql_str="INSERT INTO gallery_index (path, user, image_name) VALUES ('$storage_path', '', '$name');"
            sqlite3 $DB_PATH "$sql_str"
        else
            # 处理子目录
            sub_page=1
            while true; do
                sub_storage_path="${storage_path}/${name}"
                sub_response=$(curl -s -X POST "$url" \
                    -H "$headers" \
                    -H "Content-Type: application/json" \
                    -d "{\"path\": \"$sub_storage_path\", \"password\": \"$password\", \"page\": $sub_page, \"per_page\": $per_page, \"refresh\": false}")

                sub_code=$(echo "$sub_response" | jq -r '.code')
                if [[ "$sub_code" != "200" ]]; then
                    echo "alist server not running"
                    break
                fi

                sub_content=$(echo "$sub_response" | jq -c '.data.content[]')
                if [[ -z "$sub_content" ]]; then
                    break
                fi

                for sub_item in $sub_content; do
                    sub_is_dir=$(echo "$sub_item" | jq -r '.is_dir')
                    sub_name=$(echo "$sub_item" | jq -r '.name')

                    if [[ "$sub_is_dir" == "false" ]]; then
                        sub_sql_str="INSERT INTO gallery_index (path, user, image_name) VALUES ('$sub_storage_path', '$name', '$sub_name');"
                        sqlite3 $DB_PATH "$sub_sql_str"
                    fi
                done

                sub_total=$(echo "$sub_response" | jq -r '.data.total')
                if [[ "$sub_total" -lt "$per_page" ]]; then
                    break
                fi
                ((sub_page++))
            done
        fi
    done

    total=$(echo "$response" | jq -r '.data.total')
    if [[ "$total" -lt "$per_page" ]]; then
        break
    fi
    ((page++))
done
