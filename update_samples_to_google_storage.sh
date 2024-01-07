#!/bin/bash

# このスクリプトは、次の手順で使用します：
# 1. `YOUR_BUCKET_NAME` を実際のGoogle Cloud Storageバケット名に置き換えます。
# 2. ローカルのHugoプロジェクトでスクリプトを保存し、実行可能にします。
# 3. 必要に応じてHugoのビルドコマンドにオプションを追加します。
# 4. スクリプトを実行してHugoサイトをビルドし、Google Cloud Storageに同期します。

# **注意：**
# 初回セットアップ時には、バケットの権限設定を行う必要があります。以下のコマンドで、すべてのユーザーにバケット内のオブジェクトの閲覧権限を付与することができますが、セキュリティ上のリスクを理解し、適切な権限管理が行えるようにしてください。
# gsutil iam ch allUsers:objectViewer gs://${BUCKET_NAME}

# 変更があった実際のバケット名に置き換えてください。
BUCKET_NAME="YOUR_BUCKET_NAME"

# Hugoのビルド
hugo

# Google Cloud Storageにアップロード (ディレクトリ同期)
gsutil -m rsync -r public gs://${BUCKET_NAME}

# パブリックアクセスを有効にする（初回のみでOK、更新時に行ってもOK、一部非公開などがある場合はこの限りではない）
gsutil iam ch allUsers:objectViewer gs://YOUR_BUCKET_NAME

# ウェブサイトのURLを表示 (YOUR_BUCKET_NAMEを適切な値に置き換えてください)
echo "Your website is now available at https://${BUCKET_NAME}.storage.googleapis.com"
