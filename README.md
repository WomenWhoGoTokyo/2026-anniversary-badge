# 2026-anniversary-badge

このリポジトリは Women Who Go Tokyo 10 周年記念バッジ (通称 `wwgt2026badge`) のファームウェア / 回路図 / 組み立てガイドをまとめたものです。
RP2040 Zero と TinyGo を組み合わせた、はんだ付けから始められるワークショップ向けキットです。

> ## Special Thanks
>
> 本ワークショップは [**TinyGo Keeb**](https://tinygo-keeb.org/) コミュニティの多大なご協力のもと制作しました。
> ベースとなった [tinygo-keeb/workshop-conf2025badge](https://github.com/tinygo-keeb/workshop-conf2025badge) を含む、
> サンプルコード・回路設計・ワークショップ運営等の惜しみない共有に、心より御礼申し上げます。

不明点は Issue で気軽に質問してください。

## このリポジトリの内容

- `src/` : TinyGo で書かれたサンプルファームウェア
- `wwgt2026badge/` : KiCad プロジェクト (回路図・PCB データ)
- `lib/` : 自作シンボル / フットプリントライブラリ
- `tools/imgconv/` : 液晶用 RGB565 画像の生成ツール
- `build/BUILD.md` : ハードウェアの組み立て手順

## 回路図

KiCanvas でブラウザから直接閲覧できます。

- [kiCanvas で開く](https://kicanvas.org/?repo=https%3A%2F%2Fgithub.com%2Fsago35%2F2026-anniversary-badge%2Ftree%2Fmain%2Fwwgt2026badge)

## ハードウェアの組み立て

はんだ付けの手順、組み立て順序、必要な工具などはビルドガイドにまとめています。

- [ビルドガイド](./build/BUILD.md)

# 環境設定

[TinyGo Keeb の環境設定](https://github.com/tinygo-keeb/workshop-conf2025badge#%E7%92%B0%E5%A2%83%E8%A8%AD%E5%AE%9A) をご確認ください。


# 開発対象

このバッジのマイコンは RP2040 (Cortex M0+) で、マイコンボードとして [Waveshare RP2040-Zero](https://www.waveshare.com/wiki/RP2040-Zero) を使用しています。TinyGo のターゲット名は `waveshare-rp2040-zero` です。

主な機能は次の通りです。

- Waveshare RP2040-Zero
- ST7789 液晶
- LED 付きロータリーエンコーダー
- アナログジョイスティック
- RGB LED 付きキースイッチ × 2
- ブザー
- Grove 端子

ピン配置の概要は次の通りです。詳細は回路図を参照してください。

| 機能 | ピン |
| ---- | ---- |
| 上ボタン | GPIO0 |
| 下ボタン | GPIO1 |
| ロータリーエンコーダー A / B | GPIO2 / GPIO3 |
| SK6812MINI-E (WS2812B 互換) | GPIO4 |
| エンコーダースイッチ | D5 (GPIO5) |
| エンコーダー LED | D7 (GPIO7) |
| ブザー | D8 (GPIO8) PWM4 |
| ST7789 RST / DC / CS / BL | GPIO9 / GPIO12 / GPIO13 / GPIO14 |
| ジョイスティック X / Y | GPIO29 / GPIO28 (ADC) |

# サンプルファームウェアの実行

最初にこのリポジトリを clone してください。以降のコマンドはリポジトリのルートから実行します。

```
$ git clone https://github.com/WomenWhoGoTokyo/2026-anniversary-badge
$ cd 2026-anniversary-badge
```

サンプルは `src/` 以下に 1 ファイル 1 機能のフラットな構成で置いてあります。複数のファイルがそれぞれ `package main` を持っているため、ビルドするときはディレクトリではなく **ファイルを直接指定** します。

```
$ tinygo build -o out.uf2 --target waveshare-rp2040-zero --size short ./src/blink.go
```

## ビルド + 書き込み (ブートローダー経由)

RP2040 Zero の `BOOT` ボタンを押しながら `RESET` ボタンを押す (または USB を抜き差しする) と、ブートローダーが起動して PC から外付けドライブとして見えます。あとは `*.uf2` ファイルを D&D でコピーすれば書き込み完了です。

この方式は TinyGo 以外で作られた uf2 にも使えるので、最初に動作確認したいときに便利です。

## ビルド + 書き込み (tinygo flash)

`tinygo flash` を使えば、ビルドと書き込みを一度に実行できます。エラーメッセージが出なければ成功です。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/blink.go
```

シリアル出力の確認は `tinygo monitor` を使います。ポートを自動検出できない場合は `--port` で明示指定してください。

```
$ tinygo ports
Port                 ID        Boards
COM7                 2E8A:000A waveshare-rp2040-zero

$ tinygo monitor --port COM7
```

`tinygo flash --monitor` でまとめて実行することもできますが、環境によってはポートを誤判定するので、うまく動かないときは別々に実行してください。

# サンプル一覧

`src/` には次のサンプルが入っています。それぞれ `tinygo flash --target waveshare-rp2040-zero ./src/<ファイル名>` で書き込めます。

| ファイル | 内容                                 |
| -------- |------------------------------------|
| `blink.go` | ロータリーエンコーダー LED の点滅                |
| `sk6812.go` | キースイッチ下の RGB LED を順番に光らせる          |
| `key_input.go` | 上下キーの押下検出                          |
| `rotary.go` | ロータリーエンコーダーの回転値取得                  |
| `encorder_sw.go` | ロータリーエンコーダー押下の検出                   |
| `joystick.go` | アナログジョイスティックの XY 値取得               |
| `st7789_txt.go` | ST7789 液晶へのテキスト表示                  |
| `buzzer.go` | ブザーでドレミファソラシドを鳴らす                  |
| `all.go` | 全機能を統合したデモ (ロゴ表示 + 入力連動 LED + ブザー) |

## L チカ (blink.go)

ロータリーエンコーダー LED (D7) を 1 秒間隔で点滅させます。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/blink.go
```

LED が光ったら成功です。`time.Sleep` の値を変えて点滅速度を調整してみましょう。

## RGB LED (sk6812.go)

基板上の SK6812MINI-E は 2 個搭載されています。`WriteRaw([]uint32{...})` で複数の LED を一度に制御できます。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/sk6812.go
```

色は uint32 で指定し、最上位から 8bit ずつ Green / Red / Blue / (SK6812 RGBW の場合 White) を表します。

```go
colors := []uint32{
    0xFFFFFFFF, // white
    0xFF0000FF, // green
    0x00FF00FF, // red
    0x0000FFFF, // blue
}
```

## キー入力 (key_input.go)

上下 2 キーをマイコンの GPIO 入力で読み取ります。プルアップにしているので、押下時に Low になります。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/key_input.go
$ tinygo monitor
上側のボタンが押されました
下側のボタンが押されました
```

## ロータリーエンコーダー (rotary.go)

`tinygo.org/x/drivers/encoders` の Quadrature ドライバを使います。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/rotary.go
$ tinygo monitor
value:  -1
value:  0
value:  1
```

## エンコーダースイッチ (encorder_sw.go)

ロータリーエンコーダーは押し込むとボタンとしても動作します。GPIO5 を Low 検出するだけです。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/encorder_sw.go
$ tinygo monitor
encorder sw is pressed!!
```

## アナログジョイスティック (joystick.go)

XY の二軸を ADC で読み取り、16 進数で表示します。何もしていない時は `0x8000` 付近の値が表示されます。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/joystick.go
$ tinygo monitor
8000 8000
6E10 7E10
```

## 液晶 (st7789_txt.go)

`tinygo.org/x/drivers/st7789` と `tinygo.org/x/tinyfont` を組み合わせて、の液晶に文字を表示します。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/st7789_txt.go
```

書き込みに成功すると "Hello / Gophers!" の文字が表示されます。

画像を表示したい場合は、`tools/imgconv` で PNG を RGB565 形式に変換し、`go:embed` で埋め込みます。

```
$ go run ./tools/imgconv -in src/images/badge240.png -out src/images/badge.rgb565
```

## ブザー (buzzer.go)

GPIO トグルでブザーを鳴らし、ドレミファソラシドを再生します。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/buzzer.go
```

`all.go` では PWM を使った実装になっています。長時間鳴らし続ける用途では PWM 方式の方が CPU 負荷を抑えられます。

## 全機能デモ (all.go)

液晶のロゴ表示、入力連動の RGB LED、ジョイスティック方向に応じた虹色表示、エンコーダー押下でブザー ON / OFF、エンコーダー回転で液晶を回転、といった全部入りのデモです。

```
$ tinygo flash --target waveshare-rp2040-zero --size short ./src/all.go
```

# 参考リンク

- [tinygo-keeb/workshop-conf2025badge](https://github.com/tinygo-keeb/workshop-conf2025badge) — 本リポジトリのベースとなったワークショップ
- [sago35/tinygo-keyboard を用いて自作キーボードを作ろう](https://qiita.com/sago35/items/b008ed03cd403742e7aa)
- [koebiten — TinyGo 向け 2D ゲームエンジン](https://github.com/sago35/koebiten)
- [基礎から学ぶ TinyGo の組込み開発](https://sago35.hatenablog.com/entry/2022/11/04/230919)
