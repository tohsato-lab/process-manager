CREATE DATABASE IF NOT EXISTS process_manager_db;
USE process_manager_db;

CREATE TABLE IF NOT EXISTS `main_processes` (
  `id` VARCHAR(100) NOT NULL PRIMARY KEY COMMENT 'ID',
  `status` VARCHAR(100) NOT NULL COMMENT '状態',
  `filename` VARCHAR(200) NOT NULL COMMENT 'ファイル名',
  `target_file` VARCHAR(200) NOT NULL COMMENT '起動ファイル名',
  `env_name` VARCHAR(200) NOT NULL COMMENT '環境名',
  `start_date` datetime DEFAULT NULL COMMENT '実行開始日時',
  `complete_date` datetime DEFAULT NULL COMMENT '実行完了日時',
  `pid` int DEFAULT NULL COMMENT 'PID',
  `comment` VARCHAR(10922) DEFAULT NULL COMMENT 'コメント',
  `upload_date` datetime DEFAULT NULL COMMENT '投稿日時',
  `in_trash` Bool DEFAULT NULL COMMENT 'ゴミ箱'
) DEFAULT CHARACTER SET=utf8;

CREATE USER IF NOT EXISTS 'golang'@'%' IDENTIFIED BY 'process_manager';
GRANT ALL PRIVILEGES ON process_manager_db.* to 'golang'@'%';
