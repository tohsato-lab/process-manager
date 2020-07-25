# Process Manager

プログラムを予めアップロードしておけば、自動で実行してくれるWebアプリケーション。  

# DEMO

![Screenshot from 2020-07-25 23-56-09](https://user-images.githubusercontent.com/33301907/88459781-c910d900-ced2-11ea-9b49-86ba85aba54e.png)

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
docker exec -it server /bin/bash
cd process-manager
bash launch.sh
```

# Note

# Author

* Hirorittsu
* 遠里研究室１期生

# License


