## 処理場所と段落の関係
- frontend
  - backend
  - DB: conda_server_table, process_table
    - conda
    - DB: calc_process_table

## 初期化（websocket）
- - conda_server_tableのstatusをすべてdisconnectに更新

## サーバ登録
- ipとportをbackendにPOST
  - 指定されたcondaに対してwebsocketのハンドシェイクを行う
    - connectが呼び出されたらcalc_process_tableの情報を送信
    - return success
  - ハンドシェイクできた場合はbackendでconda_server_tableにipなどをINSERT
  - read関数を非同期で実行
  - return success
- done.

## プログラム登録
- zipファイルを登録
- サーバリストをbackendへGET
  - conda_server_tableからSELECTでipなどを取得
  - return server list
- ip一覧からipを選択し、condaへenv一覧をGETする
- - env情報を取得
- - return env list
- backendへファイルをPOST
  - ipをもとにさらにzipをcondaへPOST
    - zipを展開
    - return success
  - プログラムをprocess_tableへ登録
  - return success
- done.
