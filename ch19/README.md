# 第19章 Go言語とコンテナ

## 19.1 仮想化

ハードウェアとOSの間にレイヤーを設けることで、1台のハードウェア上に複数のOSやシステムを安全に共存させること。  
クラウドコンピューティングを支える重要技術。  

### 仮想化の種類
- エミュレーション
  - CPUを完全にエミュレーションすることで別のハードウェアのソフトも動かせる
  ![a](hyper-visor.jpg)
  - ref: [仮想化入門（ハイパーバイザー）](https://qiita.com/inasato/items/8489f2234160707c0110)
  - Type 1
    - = 準仮想化？ゲストOSはハイパーバイザーの上にいることを知っている。
    - ex. KVM, Hyper-V, VMWare vSphere
  - Type 2
    - = 完全仮想化？ゲストOSはハイパーバイザーの上にいることを知らない。
    - ex. QEMU, VirtualBox, VMWare vSphere
  - ref: https://www.redhat.com/ja/topics/virtualization/what-is-a-hypervisor

- ネイティブ仮想化
  - 同じアーキテクチャのCPUに限定
  - エミュレーション不要なので高速

### 19.1.1 仮想化は低レイヤーの技術の組み合わせ

- [PopekとGoldbergの仮想化要件](https://ja.wikipedia.org/wiki/Popek%E3%81%A8Goldberg%E3%81%AE%E4%BB%AE%E6%83%B3%E5%8C%96%E8%A6%81%E4%BB%B6)
  - ゲストOSで実行したシステムコールなどの特権命令が、ホストOSや他のゲストOSに影響を起こしちゃだめ
  - あといろいろ
- Intel系CPUで「PopekとGoldbergの仮想化要件」を満たすための機能
  - VT-x
    - ユーザーモード、特権モード、ハイパーバイザー用OSのモード（<- new）
  - 拡張ページテーブル
    - ホストOSとゲストOSのメモリアドレスを変換
  - VT-d
    - PCIパススルー（ホストOSを介さずに外部ハードウェアにアクセス）
  - VT-c
    - ネットワーク仮想化

### 19.1.2 準仮想化

- センシティブ命令（ホストOSや他のゲストOSに影響を与える命令）の扱い
  - 準仮想化
    - センシティブ命令のかわりにhypercallを呼び出す
    - そのためゲストOSをこれに乗せるには、そのOSのカスタマイズが必要
    - GPUなどの新しいハードウェア対応も難しい
  - 完全仮想化
    - 割り込みが発生し、ハイパーバイザーが処理を代行する
    - 処理が遅い
    - CPUの仮想化支援機能が充実してきたため、速くなった
    - Amazon EC2はほぼこっち

## 19.2 コンテナ

- コンテナ
  - OSカーネルはホストOSのものをそのまま使わせ、その他リソースを分離
  - aka. OSレベル仮想化
- コントロールグループ (cgroups)
  - 以下の使用量とアクセスを制限するカーネル機能
    - CPU
    - メモリ
    - ブロックデバイス
    - ネットワーク
    - `/dev`以下のデバイスファイル
- 名前空間 (Namespace)
  - 以下の名前空間を分離するカーネル機能
  - コンテナ側からは、ホストOSの一部のリソースしか見えなくなる
    - プロセスID
    - ネットワーク
    - マウント（ファイルシステム）
    - UTS（ホスト名）
    - IPC : Inter Process Communication
      - セマフォ、MQ（Message Queuing）、共有メモリなど
    - ユーザー（UID、GID）

## 19.3 Windows Subsystem for Linux 2（WSL2）

- WSL1
  - Linuxのシステムコールを逐一Windowsの命令に翻訳する仕組み
  - 互換性はいまいち
  - Windows - Linux間のファイルIOは速い
- WSL2
  - 完全なLinuxカーネルを軽量なVM環境で動作させる仕組み
  - WSL2が入っている環境では、WindowsもLinuxカーネルもハイパーバイザー（Hyper-V）上に乗っている
    - Windows Serverの話になるが、ハードウェアの上のレイヤーがOSだったとしても、Hyper-Vを有効化して再起動したら、ハードウェアとOSの間にHyper-Vが割り込むようにして機能するらしい
      - https://tooljp.com/windows/chigai/html/HyperVisor/HyperVisor-type1-vs-type2.html
  - Windowsとメモリ空間を共有している
    - Windowsにメモリを要求したり返還したりするメモリバルーニング機能があるが、Linuxはメモリをキャッシュに使いがちなため、うまく動作していない
    - `.wslconfig`などでメモリ量を設定する
      - https://learn.microsoft.com/ja-jp/windows/wsl/wsl-config#configure-global-options-with-wslconfig
  - initプロセスが特殊な実装に置き換えられている
    - 今はsystemdも設定したら使える（ https://learn.microsoft.com/ja-jp/windows/wsl/systemd ）
  - Linuxカーネル以外のOS部分はコンテナのような仕組み
    - ref: [ここが変だよ「WSL2」](https://logmi.jp/tech/articles/326106)
    > なので実はWSLの各ディストロは、WSL自体は1個のVMなのですが、その中の各ディストロはコンテナというおもしろい小ネタです。
  - Windows - WSL2間のファイルIOはネットワークファイルシステムのプロトコルを使っているので遅い
  - WindowsとWSL2は別々のネットワーク
    - 双方間の通信も、ルーティングされることで可能

## 19.4 libcontainerでコンテナを自作する

何の成果も得られませんでした  
