# post-master

※ playwrightのお試しスクレイピングプログラム

投稿や外部更新データを定期的に取得し、データベースに保存
保存されたデータで投稿用ファイルなどを生成してアウトプットする


## Concept:
- 量産化するために普遍的な構造にする
- 更新可能な単位に切り分ける(Hugo Themeを変えればデザ‐んがかわり、記事ができようされるなど)
- 


## Usage: 
基本的には実行関数をクライアントに登録して実行開始

1. 定期スクレイピング
2. データ整形、保存
3. 整形して記事化（カテゴリ、タグ、地図など）
4. 生成記事の表示（Hugo）



## For example:
1. 任意でドメインを取得
2. Google Cloud DNSでゾーン作成
3. レジストラでNSレコードを設定
4. Google Search Console で新規サイトを登録
5. レジストラ側のDNSレコードにTXTでセキュリティトークンを追加してサイト作成
6. Google Cloud StorageでCNAMEを登録
7. GCSでドメイン名（www.study-hugo.com）のバケットを生成
8. 画面またはgsutilでビルドしたファイルをアップロード
9. ウェブの設定でメインページをindex.htmlを設定
10. アップロードしたファイル・フォルダを公開に設定
11. 基本状態完了
12. [Google Functions] スクレイピング -> [GCS].mdファイル保存
13. [Google Cloud Schedule] 定時実行指定 -> スクレイピング関数指定


この他、Next.js&Vercelなども便利です。
1. [Google Functions] スクレイピング -> [GCS].mdファイルをgithubにpush
2. [Google Cloud Schedule] 定時実行指定 -> スクレイピング関数指定


## TODO:
TODO.goを参照
//  必要な場所を指定して親構造ごと持ってくる

//  適切な情報を適切な項目に代入する

//  「この関数外」データの被りをチェック DBを正常に更新できたかで判別可能

//  「この関数外」HUGO設置、自動デプロイ