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

## サーバ非アクティブ化
- 接続を一旦停止するipをbackendへGET
  - ipをもとにcondaへGET
    - channelを経由してwebsocketでdisconnectシグナルを送信
  - disconnectシグナルだった場合は関数脱出
  - return success
- done

## サーバ削除
- ipをbackendにGET
  - conda_server_tableからipをDROP
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
  - process_tableでreadyのidをcondaへGET
    - calc_process_tableの情報をもとに非同期で実行
    - calc_process_tableのstatusをrunningに更新
    - プログラムが終了した場合、channel経由でwebsocketにid情報を流す
    - return success
  - process_tableのstatusをrunningに更新
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - プログラムが終了した時にrsyncで同期、他に実行できるプロセスがあれば更新
  - done

## プロセスkill
- プロセスリストからkillボタンをクリック
- backendへidをGET
  - condaへidをGET
    - calc_process_tableの情報をもとにPIDでkill
    - channel経由でwebsocketにkillしたプロセスのid情報を流す
    - return success
  - process_tableのstatusをsyncに更新
  - rsyncで同期
  - process_tableのstatusをkilledに更新
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - return success
- done

## プロセス削除
- プロセスリストからdeleteボタンをクリック
- backendへidをGET
  - condaへidをGET
    - idをもとにディレクトリ削除
    - calc_process_tableのidをDROP
    - channel経由でwebsocketにdeleteしたプロセスのid情報を流す
    - return success
  - process_tableのidをDROP
  - idをもとにディレクトリがあれば削除
  - channel経由でwebsocketにprocess_tableの一覧を流す
  - return success
- done