CREATE USER IF NOT EXISTS 'admin'@'%' IDENTIFIED WITH caching_sha2_password BY 'admin123';
CREATE DATABASE IF NOT EXISTS `king_storage` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_storage.* TO 'admin'@'%';

-- CREATE TABLE quote_day
DROP TABLE IF EXISTS `king_storage`.`quote_day`;
CREATE TABLE `king_storage`.`quote_day` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(11,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(11,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(11,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(11,2) NOT NULL COMMENT '最低价',
    `yesterday_closed` DECIMAL(11,2) NOT NULL COMMENT '昨日收盘价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `num_of_year` INT NOT NULL COMMENT '天数',
    `xd` DOUBLE NOT NULL COMMENT '前复权比例',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_code_date ON `king_storage`.`quote_day`(`code`,`date`);

DROP TABLE IF EXISTS `king_storage`.`quote_week`;
CREATE TABLE `king_storage`.`quote_week` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `open` DECIMAL(11,2) NOT NULL COMMENT '开盘价',
    `close` DECIMAL(11,2) NOT NULL COMMENT '收盘价',
    `high` DECIMAL(11,2) NOT NULL COMMENT '最高价',
    `low` DECIMAL(11,2) NOT NULL COMMENT '最低价',
    `yesterday_closed` DECIMAL(11,2) NOT NULL COMMENT '昨日收盘价',
    `volume` BIGINT NOT NULL COMMENT '交易量',
    `account` DECIMAL(18,2) NOT NULL COMMENT '金额',
    `date` TIMESTAMP NOT NULL COMMENT '开始时期',
    `num_of_year` INT NOT NULL COMMENT '周数',
    `xd` DOUBLE NOT NULL COMMENT '前复权比例',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_code_date ON `king_storage`.`quote_week`(`code`,`date`);

-- CREATE TABLE stock
DROP TABLE IF EXISTS `king_storage`.`stock`;
CREATE TABLE `king_storage`.`stock` (
    `code` CHAR(8) NOT NULL COMMENT '股票代码',
    `name` VARCHAR(32) NOT NULL COMMENT '名称',
    `suspend` VARCHAR(32) NOT NULL COMMENT '停牌状态',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY(`code`)
);


CREATE DATABASE IF NOT EXISTS `king_auth` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_auth.* TO 'admin'@'%';

-- CREATE TABLE passport
DROP TABLE IF EXISTS `king_auth`.`passport`;
CREATE TABLE `king_auth`.`passport` (
    `id` CHAR(21) NOT NULL PRIMARY KEY,
    `account` VARCHAR(32) NOT NULL COMMENT '账户',
    `code` VARCHAR(32) COMMENT 'code',
    `salt` VARCHAR(32) NOT NULL COMMENT '盐',
    `salt_password` VARCHAR(64) NOT NULL COMMENT '盐_密码',
    `email` VARCHAR(32) COMMENT '邮箱',
    `phone` VARCHAR(15) COMMENT '电话',
    `status` TINYINT NOT NULL COMMENT '状态',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);

CREATE UNIQUE INDEX idx_account ON `king_auth`.`passport`(`account`);
CREATE UNIQUE INDEX idx_code ON `king_auth`.`passport`(`code`);
CREATE UNIQUE INDEX idx_email ON `king_auth`.`passport`(`email`);
CREATE UNIQUE INDEX idx_phone ON `king_auth`.`passport`(`phone`);


CREATE DATABASE IF NOT EXISTS `king_account` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_account.* TO 'admin'@'%';

-- CREATE TABLE assets
DROP TABLE IF EXISTS `king_account`.`assets`;
CREATE TABLE `king_account`.`assets` (
    `user_id` CHAR(19) NOT NULL COMMENT 'passport 表 id',
    `fund_no` CHAR(19) NOT NULL COMMENT 'fund no',
    `type` TINYINT NOT NULL COMMENT '类型', 
    `cash_position` DECIMAL(11,2) NOT NULL COMMENT '头寸',
    `code` VARCHAR(8) NOT NULL COMMENT '代码',
    `name` VARCHAR(32) NOT NULL COMMENT '名称',
    `open_interest` INT NOT NULL COMMENT '持仓量',
    `open_id` CHAR(19) NOT NULL COMMENT '开仓 id',
    `first_buy_datetime` DATETIME NOT NULL COMMENT '第一次购买时间',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_fund_no ON `king_account`.`assets`(`fund_no`);
CREATE INDEX idx_user_id ON `king_account`.`assets`(`user_id`);
CREATE UNIQUE INDEX idx_user_id_fund_no_code ON `king_account`.`assets`(`user_id`, `fund_no`, `code`);

-- CREATE TABLE fund
DROP TABLE IF EXISTS `king_account`.`fund`;
CREATE TABLE `king_account`.`fund` (
    `alias_name` VARCHAR(32) NOT NULL COMMENT '资金账户别名',
    `user_id` CHAR(19) NOT NULL COMMENT 'passport 表 id',
    `fund_no` CHAR(19) NOT NULL PRIMARY KEY,
    `opening_cash` DECIMAL(11,2) NOT NULL COMMENT '初始金额',
    `end_cash` DECIMAL(11,2) COMMENT '剩余金额',
    `yesterday_end_cash` DECIMAL(11,2) COMMENT '昨日剩余金额',
    `status` TINYINT NOT NULL COMMENT '状态(0:启用,1:冻结)',
    `init_datetime` DATETIME NOT NULL COMMENT '初始化时间',
    `version` INT NOT NULL COMMENT '版本号',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_user_id_fund_no ON `king_account`.`fund`(`user_id`, `fund_no`);


-- CREATE TABLE transaction_record
DROP TABLE IF EXISTS `king_account`.`transaction_record`;
CREATE TABLE `king_account`.`transaction_record` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `user_id` CHAR(19) NOT NULL COMMENT 'passport 表 id',
    `fund_no` CHAR(19) NOT NULL COMMENT 'fund no',
    `action` TINYINT NOT NULL COMMENT '动作(buy/sell)',
    `assets_type` TINYINT NOT NULL COMMENT '资产类型',
    `assets_code` CHAR(8) NOT NULL COMMENT '资产代码',
    `assets_name` VARCHAR(32) NOT NULL COMMENT '资产名称',
    `close_price` DECIMAL(11,2) NOT NULL COMMENT '成交价',
    `volume` INT NOT NULL COMMENT '成交量',
    `datetime` DATETIME NOT NULL COMMENT '成交时间',
    `status` TINYINT NOT NULL COMMENT '状态',
    `open_id` CHAR(19) NOT NULL COMMENT '开仓 id',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_datetime ON `king_account`.`transaction_record`(`datetime`);
CREATE INDEX idx_user_id_fund_no ON `king_account`.`transaction_record`(`user_id`, `fund_no`);

-- CREATE TABLE transaction_fee
DROP TABLE IF EXISTS `king_account`.`transaction_fee`;
CREATE TABLE `king_account`.`transaction_fee` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `record_id` CHAR(19) NOT NULL COMMENT '交易 id',
    `item` VARCHAR(16) NOT NULL COMMENT '项目(费用)名称',
    `money` DECIMAL(11,2) NOT NULL COMMENT '金额',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE INDEX idx_record_id ON `king_account`.`transaction_fee`(`record_id`);


CREATE DATABASE IF NOT EXISTS `king_cron` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL ON king_cron.* TO 'admin'@'%';

-- CREATE TABLE scheduler_record
DROP TABLE IF EXISTS `king_cron`.`scheduler_record`;
CREATE TABLE `king_cron`.`scheduler_record` (
    `id` CHAR(19) NOT NULL PRIMARY KEY,
    `name` VARCHAR(64) NOT NULL COMMENT '名称',
    `date` TIMESTAMP NOT NULL COMMENT '日期',
    `service_name` VARCHAR(512) NOT NULL COMMENT '服务名',
    `func_name` VARCHAR(512) NOT NULL COMMENT '方法名',
    `status` VARCHAR(16) NOT NULL COMMENT '状态', 
    `code` VARCHAR(16) COMMENT '终态',
    `error_msg` TEXT COMMENT '错误信息',
    `parent_id` CHAR(19) COMMENT '父键',
    `create_timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `modify_timestamp` TIMESTAMP COMMENT '修改时间'
);
CREATE UNIQUE INDEX idx_name_date ON `king_cron`.`scheduler_record`(`name`, `date`);
