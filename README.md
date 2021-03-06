# Process Manager

プログラムを予めアップロードしておけば、自動で実行してくれるWebアプリケーション。  
Webアプリなのでスマホなどでも閲覧可能。

# Dvelop
```
RENAME TABLE process_table TO main_processes;
ALTER TABLE main_processes ADD COLUMN upload_date datetime;
ALTER TABLE main_processes ADD COLUMN in_trash bool DEFAULT false;
ALTER TABLE main_processes DROP COLUMN use_vram;
ALTER TABLE main_processes DROP COLUMN exec_count;
ALTER TABLE main_processes CHANGE COLUMN targetfile target_file VARCHAR(200) NOT NULL;

CREATE TABLE main_processes_bak SELECT * FROM main_processes;
UPDATE main_processes set upload_date = complete_date;

```

# DEMO

![image](https://user-images.githubusercontent.com/33301907/121660471-41c65700-cade-11eb-9a58-2ee5776beaac.png)

# Features

特に機械学習の研究目的で作成しており、Pytorch、Keras用のプログラムも自動で実行できる。

# Requirement

* Ubuntu18.04以降（セットアップ用スクリプトがBashなため

※設定ファイルを見ながらDockerが動かせるならOSは問わない

# Installation

```bash
cd ./setup
bash docker-install.sh
cd ../
```

# Usage

```bash
docker-compose up -d
```
# Update
```bash
docker-compose down
git pull origin master
docker-compose up --build -d
```

# Note
もしローカルでanaconda環境を利用している場合は、以下のコマンドを利用してパッケージリストを出力し、`./server/docker/anaconda_packages`の中に入れると、build時に導入することが可能。
```bash
conda activate 'hogehoge'
conda env export | grep -v "^prefix: " > pytorch.yml
```

# Author

* Hirorittsu
* 遠里研究室１期生

# License


