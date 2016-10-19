
# 目的

golang で BigQuery API を利用してみる。
さらに、そのプログラムを digdag + docker  で動作させる。

## この記事で伝えたい事

* docker の image は alpine を使うとサイズを小さくできる。

* golang コードを cross compile した結果をつかうなら、 docker の image に golang 環境のインストールは不要。

* digdag + docker でpるグラムを実行させるなら、 実行するプログラムの変更の度に docker image をつくりなおす必要は無い。


# ためした環境

Mac OS X 10.11.6
go 1.7.1 darwin/amd64
docker version 1.12.0
docker-machine version 0.8.0
digdag version 0.8.17
java version "1.8.0_101"

# 準備

* Googleサービスを使う準備が必要です。
[http://www.task-notes.com/entry/20151019/1445223600](http://www.task-notes.com/entry/20151019/1445223600) の
サービスアカウントを使った認証　の章などを参照して、サービスアカウントでの認証のためのファイルを download しておく。
(ファイル名は client.json とする)


* brew をつかって go をインストールする。
```
$ brew install go
```

* docker-machine をインストールする。
[docker-machine を使って boot2docker から脱却する](http://qiita.com/voluntas/items/7bcc9402b51a2ba99096) などを参照して
docker-machine　をインストールする。

作業：
　次のように徐々に動作確認して、組み上げていった。
　　1．golang で hello world プログラムを go run で実行する。 mac 用に compile して実行する。
　　2. docker + digdag で golang での hello world プログラムを cross compile したものを実行させる。
　　3．golang で BigQuery に query を発行し、グラフを生成するプログラムを go run で実行する。 mac 用に compile して実行する。  
　　4. 上のことを組み合わせて、 docker + digdag で BigQuery にquery を発行し、グラフを生成させる。

# 作業1

```
$ cd golang
$ go run hello.go
Hello world from go-lang.
こんにちは. こちらは GO 言語です。
amd64 darwin
```
macos の 64 ビット環境で動作していることがわかる。

golang コードを compile して、実行してみる。
```
$ ./makebin.sh

$ file bin/darwin386/hello
bin/darwin386/hello: Mach-O executable i386

$ bin/darwin386/hello
Hello world from go-lang.
こんにちは. こちらは GO 言語です。
386 darwin

$ file bin/darwin64/hello
bin/darwin64/hello: Mach-O 64-bit executable x86_64
$ bin/darwin64/hello
Hello world from go-lang.
こんにちは. こちらは GO 言語です。
amd64 darwin
```

MacOS 用の 32 bit, 64 bit の実行可能ファイルの作成、実行ができている。

# 作業2

まずは、docker image を作成します。(image ファイルのサイズ比較用に ubuntu と alpine の２つを作成する)
```
$ cd linux
$ docker build -t ubuntu:14.04 .
$ docker build -f Dockerfile-alpine -t alpine .
```

サイズを調べてみます。
```
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
ubuntu              14.04               b8e4a5a11df5        9 seconds ago       205.2 MB
alpine              latest              5b8ef0e43f03        33 minutes ago      20.39 MB
```

alpine のものは ubuntu の 1/10 のサイズであることがわかる。

それぞれの image でワークフローを実行できることを確認します。
(作業１で、 linux 用の実行ファイルがすでに作成済みのはずです。それを ./bin/hello として copy しておきます)

# ubuntu の docker イメージで corss compaile したプログラムを実行してみます。
```
$ digdag run sample-01.dig -g hello
 ... 省略 ...
 Hello world from go-lang.
こんにちは. こちらは GO 言語です。
amd64 linux

$ digdag run sample-01-alpine.dig -g hello    # alpine で実行
Hello world from go-lang.
こんにちは. こちらは GO 言語です。
amd64 linux
```

# 作業3

golang で BigQuery の piublicdata の出生数のデータに query を出して年度毎の出生数を得し、
そのデータを折れ線グラフにした png を生成する処理を golang で作成する。
(clinet.data を bigquery フォルダに copy しておく必要がある)

```
$ cd bigquery
$ go run graph.go
$ opne graph.png
```

graph.go を走らせると graph.png が生成される。open で画像を表示させて確認できる。
![生成された画像](./screenshots/graph.gif)

次の作業の為に, cross compile をしておく。
```
$ ./makebin.sh graph
```

# 作業4

作業3 で作成したプログラムを digdag + docker で実行してみる。

(client.json を このファイルがあるフォルダに copy して置くこと。
cp ../bigquery/bin/linux64/graph ./graph として、 コンパイル結果のファイルを copy して置くこと。)

```
$ cd all-together

# ubuntu の docker imaage で実行する。
$ digdag run sample-01.dig -s graph

# alpine の docker imaage で実行する。
$ digdag run sample-01-alpine.dig -s graph

# TODO: 2016-10-20 時点では、alpine の image では 次のようなエラーが出る。調査中...
# png 生成のための package が不足していると予想している。

2016-10-20 00:42:32 +0900 [ERROR] (0017@+sample-01-alpine+graph): Task failed with unexpected error: Command failed with code 1
java.lang.RuntimeException: Command failed with code 1
	at io.digdag.standards.operator.ShOperatorFactory$ShOperator.runTask(ShOperatorFactory.java:193)
	at io.digdag.util.BaseOperator.run(BaseOperator.java:51)
	at io.digdag.core.agent.OperatorManager.callExecutor(OperatorManager.java:300)
	at io.digdag.cli.Run$OperatorManagerWithSkip.callExecutor(Run.java:678)
	at io.digdag.core.agent.OperatorManager.runWithWorkspace(OperatorManager.java:244)
	at io.digdag.core.agent.OperatorManager.lambda$runWithHeartbeat$2(OperatorManager.java:138)
	at io.digdag.core.agent.LocalWorkspaceManager.withExtractedArchive(LocalWorkspaceManager.java:25)
	at io.digdag.core.agent.OperatorManager.runWithHeartbeat(OperatorManager.java:136)
	at io.digdag.core.agent.OperatorManager.run(OperatorManager.java:120)
	at io.digdag.cli.Run$OperatorManagerWithSkip.run(Run.java:660)
	at io.digdag.core.agent.MultiThreadAgent.lambda$run$0(MultiThreadAgent.java:95)
	at java.util.concurrent.Executors$RunnableAdapter.call(Executors.java:511)
	at java.util.concurrent.FutureTask.run(FutureTask.java:266)
	at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1142)
	at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:617)
	at java.lang.Thread.run(Thread.java:745)
2016-10-20 03:42:32 +0900 [INFO] (0017@+sample-01-alpine^failure-alert): type: notify
error:
  * +sample-01-alpine+graph:
    Command failed with code 1

Task state is saved at /Users/katoy/github/golang-docker-digdag/all-together/.digdag/status/20161019T000000+0000 directory.
  * Use --session <daily | hourly | "yyyy-MM-dd[ HH:mm:ss]"> to not reuse the last session time.
  * Use --rerun, --start +NAME, or --goal +NAME argument to rerun skipped tasks.
```

参考：
- http://blog.techium.jp/entry/2016/06/27/090000 ワークフローエンジンDigDagでDockerを使ってみる

- http://qiita.com/asakaguchi/items/484ba262965ef3823f61 Alpine Linux で Docker イメージを劇的に小さくする

- http://qiita.com/pottava/items/970d7b5cda565b995fe7 Alpine Linux で軽量な Docker イメージを作る

- http://d.hatena.ne.jp/taknb2nch/20140609/1402292641 Go言語でプログラムを実行中のアーキテクチャとOSを取得する。

- https://speakerdeck.com/stormcat24/oqian-falsedockerimezihamadazhong-i?slide=20 お前のDockerイメージはまだ重い

- https://cloud.google.com/bigquery/query-reference?hl=ja      BigQuery クエリ リファレンス ROLLUP 関数にある query サンプル

- http://takedajs.hatenablog.jp/entry/2016/04/03/094529    golangで折れ線グラフを作る

- https://github.com/gonum/plot It provides an API for building and drawing plots in Go.
