# Process Manager
ZIPファイルでプログラムを予めアップロードして実行するシンプルななジョブスケジューラアプリケーション。Webアプリなのでスマホなどでも閲覧可能。

# DEMO
![image](https://user-images.githubusercontent.com/33301907/121660471-41c65700-cade-11eb-9a58-2ee5776beaac.png)

# Features
特に機械学習の研究目的で作成しており、Condaのカスタム環境を導入して、その上でプログラムを自動で実行できる。

# Requirement
* Ubuntu18.04以降（セットアップ用スクリプトがBashなため

※設定ファイルを見ながらDockerが動かせるならOSは問わない

# Installation
```bash
cd <process-manager repository>
mkdir log
cd ./setup
bash docker-install.sh
sudo reboot #必ず！！
```

# Usage
## とりあえず実行
```bash
cd <process-manager repository>
docker-compose up --build -d
```
ブラウザで http://localhost:8080 にアクセス

## 実際に使う際に必要な設定と実行
### 1. Condaのパッケージリストの追加
`./conda/docker/conda_packages`に、condaのパッケージリストをエクスポートして得られるymlファイルを入れる。ただし、同じ環境名になるとエラーになるため、必ず自分がわかる環境名にする必要あり。
### 2. データセット用のVolumeパス設定
dockerの中と外をつなぐ設定。`docker-compose.yml`の45行目以降に自分のデータセットが存在するディレクトリパスを追記する。
### 3. 起動
`docker-compose up --build -d`

## 停止
```bash
cd <process-manager repository>
docker-compose down
```
* PCシャットダウン時は勝手にdownされる。
* `docker-compose up`は何事もなく重複して立ち上がるため、その状態では内部で挙動が不安定になる。なので、再度立ち上げたい場合はまずこのコマンドを実行。
# Update
```bash
docker-compose down
git pull origin main
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

***
# Develop
開発用
```shell
docker-compose -f docker-compose.dev.yml up --build
docker-compose -f docker-compose.dev.yml up --detach --build conda backend
```
