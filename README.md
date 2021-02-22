# Process Manager

プログラムを予めアップロードしておけば、自動で実行してくれるWebアプリケーション。  

# DEMO

![image](https://user-images.githubusercontent.com/33301907/88943540-926a0280-d2c6-11ea-8418-4411e00177bc.png)

# Features

特に機械学習の研究目的で作成しており、Pytorch、Keras用のプログラムも自動で実行できる。

なお、予めGPUメモリの使用容量を設定すれば、容量オーバーしない複数のプログラムを同時に実行する。

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


