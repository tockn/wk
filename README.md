# wk

超素朴打刻ツール

## INSTALLATION

```
go get -u github.com/tockn/wk
```

## USAGE

おしごと開始〜〜
```
wk start
wk s
```

おしごとおわり〜〜
```
wk finish
wk f
```

今日の休憩時間（min）
```
wk rest 100
wk r 100
```

合計おしごと時間
```
wk total
wk t
```

2020年9月の合計おしごと時間
```
wk t 2020-9
```


### 時間指定もできるよ〜
```
wk s -t 8:00
```
```
wk f -t 23:30:48
```

## 仕組みとか

AM 6:00 ~ AM 5:59までが1日です。

`~/.wk`にプロジェクトごとにディレクトリが作成され、打刻履歴が月毎にcsvで入ります
