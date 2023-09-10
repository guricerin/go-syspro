# 第13章 シグナルによるプロセス間通信

間違えてたらごめんなさい。  

## 割り込みとシグナル

似て非なるもの。  

### 割り込み

- aka: ハードウェア割り込み、外部割り込み
  - CPU以外のハードウェアが起因のイベント
- flow:
  1. 外部機器がCPUに割り込み信号を送信
  1. CPUは実行中のプロセスを一時中断し、対応する割り込みハンドラを実行
- ex. キーボード、マウス、インターネットI/O、タイマーなど
- ref: 
  - [割り込み処理 - 筑波大学](http://www.coins.tsukuba.ac.jp/~yas/coins/os2-2012/2013-01-29/index.html)
  - [ハードウェア割り込み処理 - OSDN](https://ja.osdn.net/projects/linux-kernel-docs/wiki/2.3%E3%80%80%E3%83%8F%E3%83%BC%E3%83%89%E3%82%A6%E3%82%A7%E3%82%A2%E5%89%B2%E3%82%8A%E8%BE%BC%E3%81%BF%E5%87%A6%E7%90%86)

### シグナル

- aka: ソフトウェア割り込み、内部割り込み
  - CPU内部（ソフトウェア）で完結する
- プロセス間通信の一種
  - プロセスA -> カーネル -> プロセスB
  - （宛先はプロセスA自身のこともある）
- flow:
  1. カーネル or プロセスがシグナルを作成、対象プロセスに送信
     1. シグナルを作成したのがプロセスならカーネルが中継 
  2. 対象プロセスは実行中の処理を一時中断し、対応するシグナルハンドラを実行
- 単純なデータ（シグナルの種類）しか通知できない
- 対象プロセスが停止している場合、強制処理されるシグナルもあれば、再開までキューイング（待機）するシグナルもある
- ex. 0除算、`kill`コマンド、`Ctrl + C`など
- ref:
  - [Ctrl+Cとkill -SIGINTの違いからLinuxプロセスグループを理解する - ギークを目指して](http://equj65.net/tech/linuxprocessgroup/)

### `Ctrl + C`がシグナル？キー入力は割り込みでは？

- キー入力の度に割り込みは発生する
- が、`Ctrl + C`が確定したら、そこから先はCPU内部つまりシグナルで成り立つ仕組み（たぶん）
- flow（推測）:
  1. `Ctrl + C` -> （sshとかでつないでたらpty ->）tty
  2. ttyからカーネルに制御文字を渡す（`termios cc_c INTR`で検索検索ゥ！！！）
  3. カーネルが`SIGINT`シグナルをカレントプロセスに送信
  4. カレントプロセスがシグナルハンドラを実行
  5. カレントプロセスが終了（シグナルハンドラを弄っていたら、この限りではない）

### 参考: システムコール
- プロセス -> カーネル
- シグナルと違って、最大7個程度の引数を指定可能

## 13.1 シグナルのライフサイクル

1. raise
1. generate
1. send
1. handle

## 13.2 シグナルの種類

ubuntu22.04 on WSL2で`man 7 signal`実行した結果（抜粋）  

```
       Signal      Standard   Action   Comment
       ────────────────────────────────────────────────────────────────────────
       SIGABRT      P1990      Core    Abort signal from abort(3)
       SIGALRM      P1990      Term    Timer signal from alarm(2)
       SIGBUS       P2001      Core    Bus error (bad memory access)
       SIGCHLD      P1990      Ign     Child stopped or terminated
       SIGCLD         -        Ign     A synonym for SIGCHLD
       SIGCONT      P1990      Cont    Continue if stopped
       SIGEMT         -        Term    Emulator trap
       SIGFPE       P1990      Core    Floating-point exception

       SIGHUP       P1990      Term    Hangup detected on controlling terminal
                                       or death of controlling process
       SIGILL       P1990      Core    Illegal Instruction
       SIGINFO        -                A synonym for SIGPWR
       SIGINT       P1990      Term    Interrupt from keyboard
       SIGIO          -        Term    I/O now possible (4.2BSD)
       SIGIOT         -        Core    IOT trap. A synonym for SIGABRT
       SIGKILL      P1990      Term    Kill signal
       SIGLOST        -        Term    File lock lost (unused)
       SIGPIPE      P1990      Term    Broken pipe: write to pipe with no
                                       readers; see pipe(7)
       SIGPOLL      P2001      Term    Pollable event (Sys V);
                                       synonym for SIGIO
       SIGPROF      P2001      Term    Profiling timer expired
       SIGPWR         -        Term    Power failure (System V)
       SIGQUIT      P1990      Core    Quit from keyboard
       SIGSEGV      P1990      Core    Invalid memory reference
       SIGSTKFLT      -        Term    Stack fault on coprocessor (unused)
       SIGSTOP      P1990      Stop    Stop process
       SIGTSTP      P1990      Stop    Stop typed at terminal
       SIGSYS       P2001      Core    Bad system call (SVr4);
                                       see also seccomp(2)
       SIGTERM      P1990      Term    Termination signal
       SIGTRAP      P2001      Core    Trace/breakpoint trap
       SIGTTIN      P1990      Stop    Terminal input for background process
       SIGTTOU      P1990      Stop    Terminal output for background process
       SIGUNUSED      -        Core    Synonymous with SIGSYS
       SIGURG       P2001      Ign     Urgent condition on socket (4.2BSD)
       SIGUSR1      P1990      Term    User-defined signal 1
       SIGUSR2      P1990      Term    User-defined signal 2
       SIGVTALRM    P2001      Term    Virtual alarm clock (4.2BSD)
       SIGXCPU      P2001      Core    CPU time limit exceeded (4.2BSD);
                                       see setrlimit(2)
       SIGXFSZ      P2001      Core    File size limit exceeded (4.2BSD);
                                       see setrlimit(2)
       SIGWINCH       -        Ign     Window resize signal (4.3BSD, Sun)

       The signals SIGKILL and SIGSTOP cannot be caught, blocked, or ignored.

       Up to and including Linux 2.2, the default behavior for SIGSYS, SIGXCPU, SIGXFSZ, and (on architectures  other
       than  SPARC  and MIPS) SIGBUS was to terminate the process (without a core dump).  (On some other UNIX systems
       the default action for SIGXCPU and SIGXFSZ is to terminate the process without a core dump.)  Linux  2.4  con‐
       forms to the POSIX.1-2001 requirements for these signals, terminating the process with a core dump.

       SIGEMT  is  not  specified in POSIX.1-2001, but nevertheless appears on most other UNIX systems, where its de‐
       fault action is typically to terminate the process with a core dump.

       SIGPWR (which is not specified in POSIX.1-2001) is typically ignored by default on those  other  UNIX  systems
       where it appears.

       SIGIO (which is not specified in POSIX.1-2001) is ignored by default on several other UNIX systems.
```

### 13.2.1 ハンドルできないグナル

#### 実演

- SIGKILL
  - プロセスを強制終了

```sh
# ターミナルAで実行
top

# ターミナルBで実行
ps a | grep 'top'
kill -KILL {process_id}
```

- SIGSTOP
  - プロセスを一時停止・バックグラウンドジョブ化

```sh
# ターミナルAで実行
top

# ターミナルBで実行
ps a | grep 'top'
pkill -STOP {process_name}

# ターミナルAで実行
# フォアグラウンドジョブとして戻す
fg {process_name}
```

### 13.2.2 サーバーアプリケーションでハンドルするシグナル

- SIGTERM
  - `kill()`システムコール、`kill`コマンドがデフォで送信
  - プロセスの終了
- SIGHUP
  - コンソールアプリ: プロセスの終了
  - サーバアプリ: 設定ファイルの再読み込みを外部から指示 

### 13.2.3 コンソールアプリケーションでハンドルするシグナル

- SIGINT
  - `Ctrl + C`で終了
  - ハンドル可能なSIGKILL
- SIGQUIT
  - `Ctrl + \`でコアダンプを生成し終了
- SIGTSTP
  - `Ctrl + Z`で停止、バックグラウンドジョブ化
  - ハンドル可能なSIGSTOP
- SIGCONT
  - フォアグラウンドジョブ化
- SIGWINCH
  - ウィンドウサイズ変更
- SIGHUP
  - 疑似ターミナルから切断されるときに呼ばれるシグナル

### 13.2.4 たまに使うかもしれない、その他のシグナル

- SIGUSR1, SIGUSR2
  - ユーザ定義のシグナル
  - どういう意味のシグナルか、アプリ側で自由に決めて良い
- SIGPWR
  - 外部電源が切断し、無停電電源装置が使われたがバッテリー残量が低下したときにOSから送信される

## 13.3 Go言語によるシグナルの種類

```go
// syscall_unix.go

// A Signal is a number describing a process signal.
// It implements the os.Signal interface.
type Signal int
```

```go
// zerrors_linux_amd64.go

// Signals
const (
	SIGABRT   = Signal(0x6)
	SIGALRM   = Signal(0xe)
	SIGBUS    = Signal(0x7)
	SIGCHLD   = Signal(0x11)
	SIGCLD    = Signal(0x11)
	SIGCONT   = Signal(0x12)
	SIGFPE    = Signal(0x8)
	SIGHUP    = Signal(0x1)
	SIGILL    = Signal(0x4)
	SIGINT    = Signal(0x2)
	SIGIO     = Signal(0x1d)
	SIGIOT    = Signal(0x6)
	SIGKILL   = Signal(0x9)
	SIGPIPE   = Signal(0xd)
	SIGPOLL   = Signal(0x1d)
	SIGPROF   = Signal(0x1b)
	SIGPWR    = Signal(0x1e)
	SIGQUIT   = Signal(0x3)
	SIGSEGV   = Signal(0xb)
	SIGSTKFLT = Signal(0x10)
	SIGSTOP   = Signal(0x13)
	SIGSYS    = Signal(0x1f)
	SIGTERM   = Signal(0xf)
	SIGTRAP   = Signal(0x5)
	SIGTSTP   = Signal(0x14)
	SIGTTIN   = Signal(0x15)
	SIGTTOU   = Signal(0x16)
	SIGUNUSED = Signal(0x1f)
	SIGURG    = Signal(0x17)
	SIGUSR1   = Signal(0xa)
	SIGUSR2   = Signal(0xc)
	SIGVTALRM = Signal(0x1a)
	SIGWINCH  = Signal(0x1c)
	SIGXCPU   = Signal(0x18)
	SIGXFSZ   = Signal(0x19)
)
```

- SIGINTとSIGKILLは全OSで使用可能

```go
var (
  Interrupt Signal = syscall.SIGINT
  Kill      Signal = syscall.SIGKILL
)
```

## 13.4 シグナルのハンドラを書く

- C言語の場合
  - `signal()`や`sigaction()`の第一引数にフックしたいシグナル種、第二引数にユーザ定義のシグナルハンドラを指定
  - ref: [シグナルハンドラ内では非同期安全な関数のみを呼び出す - JPCERT/CC](https://www.jpcert.or.jp/sc-rules/c-sig30-c.html)
- Go言語の場合
  - チャネルを使う
    - `signal.NotifyContext()`
    - `signal.Notify()`
- みかんOSの場合
  - 割り込みとシグナルを厳密に区別してなかった（気がする）
  - 割り込みハンドラテーブルと割り込みハンドラ登録システムコールを作成し、それらを各プロセスに配ってた（気がする）

#### 実演

```sh
cd ch13/13.4_signal-handler
go run main.go
# Ctrl + C
```

```sh
cd ch13/13.4_signal-handler-2
go run main.go
# Ctrl + C
```

- コンテナ時代とシグナル
  - k8sやdockerでは外からタスクを終了させるとき、SIGTERMをコンテナ内プロセスに送信
  - コンテナで動作させるシステムを作るときは、シグナルを受け取ってお片付けしてからプロセスを終了するようなシステムを作りましょう

### 13.4.1 シグナルを無視する

`signal.Ignore()`を使う。  

#### 実演

```sh
cd ch13//13.4.1_ignore-signal
go run main.go
# 最初の5秒の間にCtrl + C

go run main.go
# 次の5秒の間にCtrl + C
```

### 13.4.2 シグナルのハンドラをデフォルトに戻す

`signal.Reset()`を使う。  

### 13.4.3 シグナルの送付を停止させる

シグナル受信を停止。  
`signal.Stop()`を使う。  

### 13.4.4 シグナルを他のプロセスに送る

`os.Process`構造体の`Signal()`メソッドを使う。  

#### 実演

```sh
# ターミナルAで実行
top

# ターミナルBで実行
cd ch13/13.4.4_send-signal
ps a | grep 'top'
go run main.go {process_id}
# Ctrl + C
```

- プロセスを外部から停止するお作法
  - SIGKILLは子プロセスまでは殺せない
    - まずSIGTERMを送信して、プロセス側に自分で終了処理させるのがよい
  - SIGSTOPで停止状態のプロセスはSIGKILL以外には反応しない
    - まずSIGCONTでプロセスを再起動後、SIGTERMを送信するのがよい

## 13.5 シグナルの応用例 (Server::Starter)

- `Server::Starter`
  - 新しいサーバを起動して新しいリクエストをそちらに流しつつ、古いサーバのリクエストが完了したら正しく終了させる
  - これを利用できるようにサーバを作れば、サービス停止時間ゼロでサーバ再起動が可能
  - https://github.com/lestrrat-go/server-starter

```sh
# start_serverコマンドを使用可能にする

# バイナリをインストール
go install github.com/lestrrat/go-server-starter/cmd/start_server

# .zshrcに以下を追記
export PATH="$HOME/go/bin:$PATH"
export GOPATH="$(go env GOPATH)"
```

### 13.5.1 Server::Starterの使い方

```sh
start_server --port {port_no} --pid-file {pid_file} -- ./{server_app}
```

### 13.5.2 Server::Starterが子プロセスを再起動する仕組み

- `start_server`にSIGHUPを送信して再起動できる

```sh
kill -HUP `cat app.pid`
```

- SIGHUPを受信した`Server::Starter`は、新しいプロセスを起動し、起動済みの子プロセスにはSIGTERMを送信する
- 子プロセスであるサーバが「SIGTERMを受信したら新規のリクエスト受付を停止し、処理中のリクエストが完了するまで待機してから終了する」という実装になっていれば、ダウンタイム無しでサービスを更新可能

### 13.5.3 Server::Starter対応のサーバーの実装例

#### 実演 

```sh
# ターミナルAで実行
cd ch13/13.5.3_graceful-restart
go build -o a.out
start_server --port 9999 --pid-file app.pid -- ./a.out

# ターミナルBで実行
# プロセスの親子関係を確認
ps aufx | grep a.out
# サーバが動いてるか確認
curl localhost:9999
# start_serverを再起動させてみる
kill -HUP `cat app.pid`
# もう一回アクセスしてみる
curl localhost:9999

# ターミナルAで実行
# Ctrl + C
```

## 13.6 Go言語ランタイムにおけるシグナルの内部実装

まず基礎知識を紹介。  

### 非同期シグナル安全な関数 (Async-safe)
- 非同期シグナル安全なシステムコール、およびそれらのみを使用する関数
- シグナルハンドラ内で実行してもよい
- シグナルハンドラ実行中に、同じシグナルを受信して同じシグナルハンドラが多段に実行されても問題ないと言える関数

```sh
# ubuntu22.04 on WSL2
man 7 signal-safety
```

### 非同期シグナル安全でない関数 (Async-unsafe)
- 非同期シグナル安全でないシステムコール、およびそれらを使用する関数
  - `malloc()`, `free()`、I/O系（`printf()` <- 内部で`malloc()`を使ってる）、その他いっぱい
  - 共有オブジェクトへのアクセス（ただし`volatile sig_atomic_t`型変数へのアクセスは大丈夫）
- シグナルハンドラ内で実行すると、解析困難なバグが生まれる
  - デッドロックの可能性がある
  - コンパイラによる意図しない最適化
- ただし、実行中にシグナルを受信しないことが保証されていたら、シグナルハンドラ内で使用してもよい
  - 最初にsignal-maskしておいて、ある段階でsignal-unmask（シグナル無視を解除）する

#### 補足
- シグナル無視（signal-ignore）
  - シグナルを受信してもシグナルハンドラを実行しない
- シグナルマスク（signal-mask）
  - そもそもシグナルを受信しない
  - どのスレッドにも受信されないシグナルは、いずれかのスレッドがunmaskするまでpending

#### 実演

```sh
# ターミナルAで実行
cd ch13/13.6.1_signal-unsafe
./run.sh
# Ctrl + Cしても終了しないことがある

# ターミナルBで強制終了
ps a | grep 'a.out'
kill -KILL {process_id}
```

```sh
cd ch13/13.6.2_signal-safe
./run.sh
# Ctrl + Cすると終了
```

### スレッドとシグナル
- スレッドごとに固有のもの
  - シグナルマスク（受信しないシグナルを設定）
- スレッド間で共有のもの
  - シグナルハンドラ

### マルチスレッド
- マルチスレッドなプロセスが受信したシグナルは、すべてのスレッド or いずれかのスレッドに送信される
  - kill(2) : すべてのスレッドに送信
  - pthread_kill(2) : 特定のスレッドに送信
- シグナルを受信したスレッドはシグナルハンドラを実行
  - ただし、システムや言語によっては異なる。例えば[Pythonではシグナルハンドラは常にメインスレッドで実行される](https://docs.python.org/ja/3.8/library/signal.html#signals-and-threads)
- もちろんシグナルハンドラ内では非同期安全な関数のみを使うこと

### 結論
- シグナルを簡単かつ安全に処理するなら、
  - 専用スレッドでシグナルを待機・処理する
  - その他すべてのスレッドはシグナルをマスクして受信しないようにする
- 専用スレッド内では非同期シグナル安全でない関数を使っても良い
  - シグナルハンドラではないため

#### 実演

```sh
cd ch13/13.6.3_multi-thread
./run.sh
# Ctrl + Cすると終了
```

### 参考
- https://codezine.jp/article/detail/4700
- https://www.jpcert.or.jp/sc-rules/c-sig30-c.html
- https://www.jpcert.or.jp/sc-rules/c-sig31-c.html
- http://www.hpcs.cs.tsukuba.ac.jp/~tatebe/lecture/h21/syspro/l7-thread.pdf
- https://kazmax.zpp.jp/cmd/s/signal.7.html

### Go言語の場合
- シグナルハンドラ使ってませんでしたね？
  - [signal.NotifyContext()](./13.4_signal-handler/main.go)
  - [signal.Notify()](./13.4_signal-handler-2/main.go)
  - 実は言語ランタイム側では一部で使ってる（後述）

#### マルチスレッド
1. シグナル処理専用スレッドを用意
1. 各スレッド初期化時に`minitSignalMask()`を実行
   1. すべてのシグナルをマスク
   1. 内部では`sigprocmask()`を実行
   1. さらにその内部では`pthread_sigmask()`を実行 <-- [実演用コード](./13.6.3_multi-thread/main.c)にも登場
1. `ensureSigM()`でシグナル処理専用スレッドにシグナル受信を許可
   1. goroutineが必ず特定のOSスレッドで実行されることを保証する`runtime.LockOSThread()`を実行
1. いずれかのスレッドが`signal`パッケージの関数を実行
   1. `ensureSigM()`が監視しているチャネル（シグナル処理専用スレッド）に更新情報を届けるため、`runtime`パッケージの`signal_enable()`や`signal_disable()`を実行
   1. `sigprocmask()`を実行し、シグナル処理専用スレッドのシグナルマスクを更新
1. シグナル処理専用スレッドは
   1. シグナル送信先のハンドラを指定する際、`sigaction()` -> `sighandler()`を実行（シグナルハンドラ）
   1. `sigsend()`が実行され、共有メモリ領域（queue）にシグナルを書き出す
      1. 内部では`atomic`パッケージ使ってるから、シグナルハンドラで実行されて共有メモリ領域を読み書きしても大丈夫なはず。しらんけど。
1. `signal`パッケージは`signal_recv()`を使って共有メモリ領域からシグナルを取り出し、`Notify()`で渡されたチャネルに伝達

## 13.7 Windowsとシグナル

- GUI用のメッセージング・ループの例  

```cpp
// ref: https://github.com/microsoft/Windows-classic-samples/blob/1d363ff4bd17d8e20415b92e2ee989d615cc0d91/Samples/RadialController/cpp/RadialController.cpp  
// VC++

    MSG msg;
    while (GetMessage(&msg, NULL, 0, 0))
    {
        TranslateMessage(&msg);
        DispatchMessage(&msg);  // Dispatch message to WindowProc

        if (msg.message == WM_QUIT)
        {
            Windows::Foundation::Uninitialize();
            break;
        }
    }
```
