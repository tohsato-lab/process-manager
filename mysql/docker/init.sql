CREATE DATABASE IF NOT EXISTS process_manager_db;
USE process_manager_db;

CREATE TABLE IF NOT EXISTS `calc_server_table`
(
    `ip`        VARCHAR(15)  NOT NULL PRIMARY KEY COMMENT 'IPアドレス',
    `port`      VARCHAR(15)  NOT NULL COMMENT 'ポート',
    `status`    VARCHAR(100) NOT NULL COMMENT '状態',
    `num_limit` INT DEFAULT 1 COMMENT 'プロセス実行制限数'
) DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `process_table`
(
    `id`            VARCHAR(100) NOT NULL PRIMARY KEY COMMENT 'ID',
    `process_name`  VARCHAR(100) NOT NULL COMMENT 'ファイル名',
    `env_name`      VARCHAR(200) NOT NULL COMMENT '環境名',
    `server_ip`     VARCHAR(15)  NOT NULL COMMENT 'プロセス実行IP',
    `comment`       VARCHAR(7281) DEFAULT '' COMMENT 'コメント',
    `status`        VARCHAR(100)  DEFAULT 'ready' COMMENT '状態',
    `upload_date`   datetime      DEFAULT NULL COMMENT '投稿日時',
    `start_date`    datetime      DEFAULT NULL COMMENT '実行開始日時',
    `complete_date` datetime      DEFAULT NULL COMMENT '実行完了日時',
    `in_trash`      Bool          DEFAULT FALSE COMMENT 'ゴミ箱',
    `args`          VARCHAR(3641) DEFAULT '' COMMENT '引数'
) DEFAULT CHARACTER SET = utf8;

CREATE TABLE IF NOT EXISTS `calc_process_table`
(
    `id`          VARCHAR(100) NOT NULL PRIMARY KEY COMMENT 'ID',
    `target_file` VARCHAR(200) NOT NULL COMMENT '起動ファイル名',
    `env_name`    VARCHAR(200) NOT NULL COMMENT '環境名',
    `status`      VARCHAR(100)  DEFAULT 'ready' COMMENT '状態',
    `pid`         INT           DEFAULT NULL COMMENT 'PID',
    `args`        VARCHAR(3641) DEFAULT '' COMMENT '引数'
) DEFAULT CHARACTER SET = utf8;

CREATE USER IF NOT EXISTS 'golang'@'%' IDENTIFIED BY 'process_manager';
GRANT ALL PRIVILEGES ON process_manager_db.* to 'golang'@'%';
