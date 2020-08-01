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
bash setup.sh
```

# Usage

```bash
cd ./setup
docker-compose up -d
```
# Update
```bash
cd ./setup
docker-compose down
docker-compose up -d
```

# Note

# Author

* Hirorittsu
* 遠里研究室１期生

# License


