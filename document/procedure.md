## 処理場所と段落の関係
- frontend
  - backend
  - DB: conda_server_table, process_table
    - conda
    - DB: calc_process_table

## 初期化（websocket）
- @
  - conda_server_tableのstatusをすべてdisconnectに更新

## サーバ登録
- ipとportをbackendにPOST
  - 指定されたcondaに対してwebsocketのハンドシェイクを行う
    - connectが呼び出されたらcalc_process_tableの情報を送信
    - return success
  - 返されたcalc_process_tableからidを参照し、process_tableの情報を上書きする
  - ハンドシェイクできた場合はbackendでconda_server_tableにipなどをINSERT
  - read関数を非同期で実行
  - return success
- done

## サーバ登録解除

## プログラム登録
- zipファイルを登録
- サーバリストをbackendへGET
  - conda_server_tableからSELECTでipなどを取得
  - return server list
- ip一覧からipを選択し、condaへenv一覧をGETする
  - @
    - env情報を取得
    - return env list
- backendへファイルをPOST
  - ipをもとにさらにzipをcondaへPOST
    - zipを展開
    - calc_process_tableへINSERT
    - return success
  - プログラムをprocess_tableへ登録
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - 非同期でプロセス更新（プロセス実行へ）
  - return success
- done

## プロセス実行
- @
  - process_tableでreadyのidをcondaへGET
    - calc_process_tableの情報をもとに非同期で実行
    - calc_process_tableのstatusをrunningに更新
    - プログラムが終了した場合、channel経由でwebsocketにid情報を流す
    - // 正常に送信完了した場合calc_process_tableの情報をDROP（ディレクトリ削除する？）
    - return success
  - process_tableのstatusをrunningに更新
  - rsyncで同期
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - done

## プロセスkill
- プロセスリストからkillボタンをクリック
- backendへidをGET
  - condaへidをGET
    - calc_process_tableの情報をもとにPIDでkill
    - channel経由でwebsocketにkillしたプロセスのid情報を流す
    - // calc_process_tableの情報をDROP（ディレクトリ削除する？）
    - return success
  - process_tableのstatusをkilledに更新
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - return success
- done

## プロセス削除