# HMDCバックエンドサーバ

## 環境構築手順

### (1) Goのセットアップ

1. goのバージョンマネージャであるgvmをインストールする。  
   ```bash
   bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
   ```

   実行後に表示される `source ...` のコマンドを叩く。

1. 他のgoバージョンを入れる前に必要なgo1.4をインストールする。  
   ```bash
   gvm install go1.4
   ```

1. その後、このリポジトリで用いるgo1.17をバイナリからインストールする。  
   ```bash
   gvm install go1.17 -B
   ```

1. 使用するバージョンを1.17に設定する。  
   ```bash
   gvm use go1.17
   ```

### (2) PostgreSQL関連のセットアップ

1. aptから必要なパッケージをインストールする。  
   ```bash
   sudo apt update
   sudo apt install -y postgresql postgresql-contrib
   sudo systemctl start postgresql.service
   ```

1. PostgreSQLのインストールによってUbuntuシステムに作成されたpostgresユーザの権限でDBを作成する。  
   ここではDB名は `hmdc` とする。
   ```bash
   sudo -u postgres createdb hmdc
   ```

1. DBにアクセスするためのPostgreSQLユーザを作成する。  
   ここではユーザ名 `jonouchi`、パスワード `shizuka` とする。
   1. psqlを呼び出す。  
      ```bash
      sudo -u postgres psql
      ```
    
   1. psqlの対話プロンプトが立ち上がるので、ユーザ名とパスワードを設定してユーザを作成する。  
      ```sql
      CREATE USER jonouchi WITH PASSWORD 'shizuka';
      ```
    
   1. 対話プロンプトを開いたまま、作成したDBへのすべての権限を上記ユーザにわたす。  
      ```sql
      GRANT ALL PRIVILEGES ON DATABASE hmdc TO jonouchi;
      ```

      実行後、対話プロンプトを `\q` で閉じる。
   
1. `.bashrc` にDBホスト名、ポート番号、DB名、ユーザ名、パスワードを環境変数に保存する処理を追加する。  
   ホスト名はlocalhost。  
   ```bash
   echo 'export DB_HOST=localhost
   export DB_PORT=5432
   export DB_NAME=hmdc
   export DB_USER=jonouchi
   export DB_PASSWORD=shizuka' >> ~/.bashrc
   ```

   その後、`.bashrc` をsourceする。  
   ```bash
   source ~/.bashrc
   ```

1. 当ファイルと同じディレクトリから、sql-migrateをインストールする。  
   ```bash
   go get -v github.com/rubenv/sql-migrate/...
   ```

1. sql-migrateコマンドを実行し、規定の構造をDBに反映する。  
   ```bash
   sql-migrate up
   ```
