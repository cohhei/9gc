# note for 9gc

## Visual Studio Code のセットアップ
- `GO: Install/Update Tools` でツール類をインストール
- GOROOTとGOPATHの設定

コンテナ内じゃなくて**ホスト側**で環境変数に**追加**してあげる必要がる
```bash
export GOPATH=$GOPATH:/root/go
export GOROOT=$GOROOT:/usr/local/go
```

[fish-shell](https://fishshell.com)の場合
```bash
set -U GOPATH $GOPATH:/root/go
set -U GOROOT $GOROOT:/usr/local/go
```