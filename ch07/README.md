# 7 UDPソケットを使ったマルチキャスト通信

## 7.1 UDP と TCP の用途の違い

- from: https://www.reddit.com/r/ProgrammerHumor/comments/krft3r/tcp_vs_udp/
  ![tcp udp](image/a.jpg)

- UDP
  - 複数のコンピューターに同時にメッセージを送ることが可能なマルチキャストとブロードキャストをサポート
- UDPを使用しているプロトコル
  - DNS
  - NTP : 時計合わせのプロトコル
  - ストリーミング動画・音声プロトコル
  - WebRTC : ブラウザ上で行う P2P のための動画・音声通信プロトコル
- TCPとUDPの中間的なプロトコル
  - RUDP : 再送処理とウインドウ制御だけを TCP から UDP に移植
  - DTLS : UDP に TLS による暗号化を載せつつハンドシェイクだけは行う


### 7.1.1 UDP が使われる場面は昔と今で変わってきている

- いまのVPN
  - SSL-VPNが増えた。SSL-VPNのうち、パケットをHTTPSで包む方式あり
- いまのDNS
  - レスポンスが512バイトを超えたらTCPにフォールバック
- いまの独自プロトコル
  - エラー処理を作り込む必要あり
- TCPのバージョンがあがり、高性能になった

> アプリケーション開発という視点で見れば、「ロスしてもよい、ロスしても順序が変わってもアプリケーションレイヤーでカバーできる、マルチキャストが必要、ハンドシェイクの手間すら惜しいなど、いくつかの特別な条件に合致する場合以外はTCP」という選択でよいでしょう。

- QUIC
  - トランスポート層のプロトコル
  - TCPのレイヤーを軽量化、UDP使用、TLSの暗号化を合体
  - QUIC上ではHTTP/3が動作する
  - エラー処理はQUICのレイヤーで行う

## 7.2 UDP と TCP の処理の流れの違い

### 7.2.1 サーバー側の実装例

[server](01-unicast/server/main.go)
[client](01-unicast/client/main.go)

- POSIXでは、直接`recvfrom()`や`sendto()`で通信してよい
  - server: `listen()`, `accept()`不要
  - client: `connect()`不要
