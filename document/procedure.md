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
  - process_tableでready・running・syncingのプロセスのstatusをwebsocket経由で確認
  - read関数を非同期で実行
  - return success
- done

## サーバ非アクティブ化
- 接続を一旦停止するipをbackendへGET
  - ipをもとにwebsocket.CloseMessageをwebsocketへ送信
  - return success
- done

## サーバ削除
- ipをbackendにGET
  - conda_server_tableからipをDELETE
  - return success
- done

## プログラム登録
- zipファイルを登録
- conda情報をbackendへGET
  - conda_server_tableからSELECTでipなどを取得
  - ipをもとにcondaへenv情報をGET
  - return conda env + server ip
- condaへファイルをPUT
  - @
    - zipを展開
    - プロセスをcalc_process_tableへ登録
    - return process id
- 登録のためbackendへGET
  - プロセスをbackendのprocess_tableへ登録
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - 非同期でプロセス実行
  - return success
- done

## プロセス実行
- @
  - process_tableでreadyのidをwebsocketで送信
    - calc_process_tableの情報をもとに非同期で実行
    - calc_process_tableのstatusをrunningに更新
    - プログラムが終了した場合、channel経由でwebsocketにid情報を流す
    - return success
  - process_tableのstatusをrunningに更新
  - プログラムが終了した時にrsyncで同期、他に実行できるプロセスがあれば更新
  - done

## プロセスkill
- プロセスリストからkillボタンをクリック
- backendへidをGET
  - websocketでcondaへidを送信
    - calc_process_tableの情報をもとにPIDでkill
    - channel経由でwebsocketにkillしたプロセスのid情報を流す
    - return success
  - process_tableのstatusをsyncingに更新
  - rsyncで同期
  - process_tableのstatusをkilledに更新
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - return success
- done

## プロセス削除
- プロセスリストからdeleteボタンをクリック
- backendへidをDELETE
  - condaへidをDELETE
    - idをもとにディレクトリ削除
    - calc_process_tableのidをDELETE
    - channel経由でwebsocketにdeleteしたプロセスのid情報を流す
    - return success
  - process_tableのidをDELETE
  - idをもとにディレクトリがあれば削除
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - return success
- done