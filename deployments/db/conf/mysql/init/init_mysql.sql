CREATE USER IF NOT EXISTS 'admin'@'%' IDENTIFIED WITH caching_sha2_password BY 'admin123';
CREATE DATABASE IF NOT EXISTS `king_storage` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_storage.* TO 'admin'@'%';

-- create table quote_day
drop table if exists `king_storage`.`quote_day`;
create table `king_storage`.`quote_day` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    `yesterday_closed` DECIMAL(10,2) NOT NULL COMMENT '昨日收盘价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `num_of_year` INT NOT NULL COMMENT '天数',
    `xd` DOUBLE NOT NULL COMMENT '前复权比例',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_code_date on `king_storage`.`quote_day`(`code`,`date`);

drop table if exists `king_storage`.`quote_week`;
create table `king_storage`.`quote_week` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(10,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(10,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(10,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(10,2) NOT NULL COMMENT '最低价',
    `yesterday_closed` DECIMAL(10,2) NOT NULL COMMENT '昨日收盘价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '开始时期',
    `num_of_year` INT NOT NULL COMMENT '周数',
    `xd` DOUBLE NOT NULL COMMENT '前复权比例',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_code_date on `king_storage`.`quote_week`(`code`,`date`);

-- create table stock
drop table if exists `king_storage`.`stock`;
create table `king_storage`.`stock` (
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `name` VARCHAR(32) NOT NULL COMMENT '名称',
    `suspend` VARCHAR(32) NOT NULL COMMENT '停牌状态',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间',
     PRIMARY KEY(`code`)
);

CREATE DATABASE IF NOT EXISTS `king_account` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_account.* TO 'admin'@'%';

-- create table account
drop table if exists `king_account`.`account`;
create table `king_account`.`account` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `username` VARCHAR(32) COMMENT '用户名',
    `password` VARCHAR(64) NOT NULL COMMENT '密码',
    `nick_name` VARCHAR(32) COMMENT '密码', 
    `phone` VARCHAR(15) COMMENT '电话',
    `email` VARCHAR(32) COMMENT 'email',
    `status` TINYINT NOT NULL COMMENT '状态',
    `register_datetime` DATETIME NOT NULL COMMENT '注册时间',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);

create index idx_username on `king_account`.`account`(`username`);
create index idx_phone on `king_account`.`account`(`phone`);
create index idx_email on `king_account`.`account`(`email`);
create unique index idx_username_phone_email on `king_account`.`account`(`username`,`phone`,`email`);

-- create table assets
drop table if exists `king_account`.`assets`;
create table `king_account`.`assets` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `fund_id` CHAR(19) NOT NULL COMMENT 'fund 表 id',
    `user_id` CHAR(19) NOT NULL COMMENT 'account 表 id',
    `type` TINYINT NOT NULL COMMENT '类型', 
    `cash_position` DECIMAL(10,2) NOT NULL COMMENT '头寸',
    `code` VARCHAR(8) NOT NULL COMMENT '代码',
    `open_interest` INT NOT NULL COMMENT '持仓量',
    `first_buy_datetime` DATETIME NOT NULL COMMENT '第一次购买时间',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_fund_id on `king_account`.`assets`(`fund_id`);
create index idx_user_id on `king_account`.`assets`(`user_id`);

-- create table fund
drop table if exists `king_account`.`fund`;
create table `king_account`.`fund` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `user_id` CHAR(19) NOT NULL COMMENT 'account 表 id',
    `opening_cash` DECIMAL(10,2) NOT NULL COMMENT '初始金额',
    `end_cash` DECIMAL(10,2) COMMENT '结算金额',
    `status` TINYINT NOT NULL COMMENT '状态(1:启用,2:冻结)',
    `init_datetime` DATETIME NOT NULL COMMENT '初始化时间',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_user_id on `king_account`.`fund`(`user_id`);


-- create table fund
drop table if exists `king_account`.`transaction_record`;
create table `king_account`.`transaction_record` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `action` TINYINT NOT NULL COMMENT '动作(buy/sell)',
    `assert_type` TINYINT NOT NULL COMMENT '资产类型',
    `assert_code` CHAR(8) NOT NULL COMMENT '资产代码',
    `close_price` DECIMAL(10,2) NOT NULL COMMENT '成交价',
    `volume` INT NOT NULL COMMENT '成交量',
    `datetime` DATETIME NOT NULL COMMENT '成交时间',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_datetime on `king_account`.`transaction_record`(`datetime`);
create index idx_assert_code on `king_account`.`transaction_record`(`assert_code`);

-- create table fund
drop table if exists `king_account`.`transaction_fee`;
create table `king_account`.`transaction_fee` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `record_id` CHAR(19) NOT NULL COMMENT '交易 id',
    `item` VARCHAR(16) NOT NULL COMMENT '项目(费用)',
    `money` DECIMAL(10,2) NOT NULL COMMENT '金额',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
create index idx_record_id on `king_account`.`transaction_fee`(`record_id`);

