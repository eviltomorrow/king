#!/bin/bash

mkdir -p $(pwd)/data/mongo/{db,log,init}
mkdir -p $(pwd)/data/mysql/{db,log,init}
mkdir -p $(pwd)/data/etcd
chmod -R 777 $(pwd)/data


# init_mongo.js
cat > $(pwd)/data/mongo/init/init_mongo.js <<EOF
db = db.getSiblingDB('transaction_db');
db.createUser({"user":"admin","pwd":"admin123","roles":[{"role":"dbOwner","db":"transaction_db"}]});
db.createCollection('metadata');
db.metadata.createIndex({date: 1, code: 1},{background: true});
EOF
chmod a+x $(pwd)/data/mongo/init/init_mongo.js

# init.sql
cat > $(pwd)/data/mysql/init/init_mysql.sql <<EOF
CREATE USER 'admin'@'%' IDENTIFIED BY 'admin123';
CREATE DATABASE \`king_storage\` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_storage.* TO 'admin'@'%';

-- create table quote_day
drop table if exists \`king_storage\`.\`quote_day\`;
create table \`king_storage\`.\`quote_day\` (
    \`id\` CHAR(19) NOT NULL PRIMARY KEY,
    \`code\` CHAR(8) NOT NULL COMMENT '股票代码',
    \`open\` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    \`close\` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    \`high\` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    \`low\` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    \`yesterday_closed\` DECIMAL(10,2) NOT NULL COMMENT '昨日收盘价',
    \`volume\` BIGINT NOT NULL COMMENT '交易量',
    \`account\` DECIMAL(18,2) NOT NULL COMMENT '金额',
    \`date\` TIMESTAMP NOT NULL COMMENT '日期',
    \`num_of_year\` INT NOT NULL COMMENT '天数',
    \`xd\` DOUBLE NOT NULL COMMENT '前复权比例',
    \`create_timestamp\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    \`modify_timestamp\` TIMESTAMP COMMENT '修改时间'
);
create index idx_code_date on \`king_storage\`.\`quote_day\`(\`code\`,\`date\`);

drop table if exists \`king_storage\`.\`quote_week\`;
create table \`king_storage\`.\`quote_week\` (
    \`id\` CHAR(19) NOT NULL PRIMARY KEY,
    \`code\` CHAR(8) NOT NULL COMMENT '股票代码',
    \`open\` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    \`close\` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    \`high\` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    \`low\` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    \`yesterday_closed\` DECIMAL(10,2) NOT NULL COMMENT '昨日收盘价',
    \`volume\` BIGINT NOT NULL COMMENT '交易量',
    \`account\` DECIMAL(18,2) NOT NULL COMMENT '金额',
    \`date\` TIMESTAMP NOT NULL COMMENT '开始时期',
    \`num_of_year\` INT NOT NULL COMMENT '周数',
    \`xd\` DOUBLE NOT NULL COMMENT '前复权比例',
    \`create_timestamp\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    \`modify_timestamp\` TIMESTAMP COMMENT '修改时间'
);
create index idx_code_date_end on \`king_storage\`.\`quote_week\`(\`code\`,\`date\`);

-- create table stock
drop table if exists \`king_storage\`.\`stock\`;
create table \`king_storage\`.\`stock\` (
    \`code\` CHAR(8) NOT NULL COMMENT '股票代码',
    \`name\` VARCHAR(32) NOT NULL COMMENT '名称',
    \`suspend\` VARCHAR(32) NOT NULL COMMENT '停牌状态',
    \`create_timestamp\` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    \`modify_timestamp\` TIMESTAMP COMMENT '修改时间',
     PRIMARY KEY(\`code\`)
);
EOF

touch $(pwd)/data/.ready