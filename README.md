# Process Manager
プログラムを予めアップロードしておけば、自動で実行してくれるWebアプリケーション。  
Webアプリなのでスマホなどでも閲覧可能。

# DEMO
![image](https://user-images.githubusercontent.com/33301907/121660471-41c65700-cade-11eb-9a58-2ee5776beaac.png)

# Features
特に機械学習の研究目的で作成しており、Pytorch、Keras用のプログラムも自動で実行できる。

# Requirement
* Ubuntu18.04以降（セットアップ用スクリプトがBashなため

※設定ファイルを見ながらDockerが動かせるならOSは問わない

# Installation
```bash
cd process-manager/setup
bash docker-install.sh
sudo reboot #必ず！！
```

# Usage
```bash
cd process-manager
docker-compose up -d
```
ブラウザで http://localhost:8080 にアクセス

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

# Develop
開発用
```shell
docker-compose -f docker-compose.dev.yml up --build
docker-compose -f docker-compose.dev.yml up --detach --build conda backend
```

# Author
* Hirorittsu
* 遠里研究室１期生
