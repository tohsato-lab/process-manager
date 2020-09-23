CREATE DATABASE IF NOT EXISTS process_manager_db;
CONNECT process_manager_db;

CREATE TABLE IF NOT EXISTS `process_table` (
  `id` VARCHAR(100) NOT NULL PRIMARY KEY COMMENT "ID",
  `use_vram` FLOAT NOT NULL COMMENT "使用VRAM容量",
  `status` VARCHAR(100) NOT NULL COMMENT "状態",
  `filename` VARCHAR(200) NOT NULL COMMENT "ファイル名",
  `targetfile` VARCHAR(200) NOT NULL COMMENT "起動ファイル名",
  `env_name` VARCHAR(200) NOT NULL COMMENT "環境名",
  `start_date` datetime DEFAULT NULL COMMENT "実行開始日時",
  `complete_date` datetime DEFAULT NULL COMMENT "実行完了日時",
  `pid` int DEFAULT NULL COMMENT "PID"
) DEFAULT CHARACTER SET=utf8;

CREATE user `golang`@'%' IDENTIFIED BY 'golang';
GRANT ALL PRIVILEGES ON process_manager_db.* to `golang`@'%' IDENTIFIED BY 'golang'
