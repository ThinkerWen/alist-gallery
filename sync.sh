#!/bin/bash

# 数据库路径
DB_PATH="./gallery.db"

# YAML配置文件路径
YAML_PATH="./config.yaml"

# 创建表的SQL语句
CREATE_TABLE_SQL="CREATE TABLE gallery_index
(
    id         INTEGER primary key autoincrement,
    path       VARCHAR(255)  default ''                not null,
    user       VARCHAR(255)  default ''                not null,
    image_name VARCHAR(255)  default ''                not null,
    image_url  varchar(2000) default ''                not null,
    created_at TIMESTAMP     default CURRENT_TIMESTAMP not null
);

create unique index idx_image_name on gallery_index (image_name);"

# 检查表是否存在
TABLE_EXISTS=$(sqlite3 $DB_PATH "SELECT name FROM sqlite_master WHERE type='table' AND name='gallery_index';")

if [[ -n $TABLE_EXISTS ]]; then
    echo "表已存在，停止插入"
    exit 0
fi

# 创建表
sqlite3 $DB_PATH "$CREATE_TABLE_SQL"
echo "建表完成！"