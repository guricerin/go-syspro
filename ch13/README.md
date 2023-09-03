# 第13章 シグナルによるプロセス間通信

## シグナルとは

- プロセス間通信の一種
- プロセスに対して送られるなんらかの通知
- 単純なデータ（シグナルの種類）しか通知できない
- 対象プロセスが停止している場合、強制処理されるシグナルもあれば、再開までキューイング（待機）するシグナルもある

## シグナルの用途

- プロセス間通信
  - プロセスA -> カーネル -> プロセスB
  - （宛先はプロセスA自身のこともある）
- ソフトウェア割り込み
  - シグナルを受け取ったプロセスは現在実行中のタスクを中断して、あらかじめ登録しておいたルーチン（シグナルハンドラ）を実行
  - ex. 不正なメモリアクセス、0除算などソフトウェア起因で発生
  - 図13.1にはハードウェア割り込みである`Ctrl + C`もこれに属するように描かれてるが......
    - ソフトウェア、ハードウェア関係なく`割り込みが発生したらシグナルは生成される`と考えればよいか？

- 参考: ハードウェア割り込み
  - ハードウェア -> CPU -> カーネル -> プロセス
  - 割り込みを受けたカーネルがシグナルを生成する
  - ex. `Ctrl + C`によるプロセス終了、タイマーなど
    - ref: https://access.redhat.com/documentation/ja-jp/red_hat_enterprise_linux_for_real_time/7/html/reference_guide/chap-hardware_interrupts

- 参考: システムコール
  - プロセス -> カーネル
  - シグナルと違って、最大7個程度の引数を指定可能

## 13.1 シグナルのライフサイクル

1. raise
2. send
3. handle

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

実演。  

- SIGKILL
  - プロセスを強制終了
  - `top`コマンドで実演する

```sh
# ターミナルAで実行
top

# ターミナルBで実行
kill -KILL {process_id}
```

- SIGSTOP
  - プロセスを一時停止・バックグラウンドジョブ化
  - `top`コマンドで実演する

```sh
# ターミナルAで実行
top

# ターミナルBで実行
pkill -STOP {process_name}

# ターミナルAで実行
# フォアグラウンドジョブとして戻す
fg {process_name}
```

### 13.2.2 サーバーアプリケーションでハンドルするシグナル

- SIGTERM
  - `kill()`システムコール、`kill`コマンドがデフォで送信
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

## 13.4 シグナルのハンドラを書く

- C言語の場合
  - `signal()`や`sigaction()`の第一引数にフックしたいシグナル種、第二引数にユーザ定義のシグナルハンドラを指定
  - ref: [シグナルハンドラ内では非同期安全な関数のみを呼び出す - JPCERT/CC](https://www.jpcert.or.jp/sc-rules/c-sig30-c.html)
- Go言語の場合
  - チャネルを使う
    - `signal.NotifyContext()`
    - `signal.Notify()`
- みかんOSの場合
  - シグナルハンドラテーブルとシグナルハンドラ登録システムコールを作成し、それらを各プロセスにコピーしてた（気がする）

実演。  

- コンテナ時代とシグナル
  - k8sやdockerでは外からタスクを終了させるとき、SIGTERMをコンテナ内プロセスに送信
  - コンテナで動作させるシステムを作るときは、シグナルを受け取ってお片付けしてからプロセスを終了するようなシステムを作りましょう


### 13.4.1 シグナルを無視する

`signal.Ignore()`を使う。  
実演。  

### 13.4.1 シグナルのハンドラをデフォルトに戻す

`signal.Reset()`を使う。  

### 13.4.1 シグナルの送付を停止させる

シグナル受信を停止。  
`signal.Stop()`を使う。  

### 13.4.4 シグナルを他のプロセスに送る

`os.Process`構造体の`Signal()`メソッドを使う。  
実演。  

- プロセスを外部から停止するお作法
  - SIGKILLは子プロセスまでは殺せない
    - まずSIGTERMを送信して、プロセス側に自分で終了処理させるのがよい
  - SIGSTOPで停止状態のプロセスはSIGKILL以外には反応しない
    - まずSIGCONTでプロセスを再起動後、SIGTERMを送信するのがよい

## 13.5 シグナルの応用例 (Server::Starter)

- Server::Starter
  - 新しいサーバを起動して新しいリクエストをそちらに流しつつ、古いサーバのリクエストが完了したら正しく終了させる
  - これを利用できるようにサーバを作れば、サービス停止時間ゼロでサーバ再起動が可能
  - https://github.com/lestrrat-go/server-starter

### 13.5.1 Server::Starterの使い方
