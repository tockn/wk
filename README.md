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
```
or
```
wk s
```

おしごとおわり〜〜

```
wk finish
```
or
```
wk f
```

### 時間指定もできるよ〜
```
wk s -t 8:00
```
```
wk f -t 23:30:48
```

## 仕組み

`~/.wk`ディレクトリにプロジェクトごとに打刻履歴がcsvで入ります
今はdefaultプロジェクトしか作れないけど、複数プロジェクト管理できるようにしたい