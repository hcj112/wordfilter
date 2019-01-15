## installation
```
go get github.com/hcj112/wordfilter
```


## run

```
$ cd $GOPATH/github.com/hcj112/wordfilter/cmd
$ go run main.go -conf=example.toml
```

### Dependencies
[sego](https://github.com/huichen/sego)

[bolt](https://github.com/boltdb/bolt)




## filter keyword API

### error code
```
// ok
OK = 0

// request error
RequestErr = -400

// server error
ServerErr = -500
```

### add keyword
[GET] /filter/keyword/add

| Name            | Type     | Remork                 |
|:----------------|:--------:|:-----------------------|
| keyword | string    | keyword for response |

response:
```
{
    "code": 0,
    "message": ""
}
```

### remove keyword
[GET] /filter/keyword/remove

| Name            | Type     | Remork                 |
|:----------------|:--------:|:-----------------------|
| keyword | string    | keyword for response |

response:
```
{
    "code": 0,
    "message": ""
}
```

### replace keyword
[GET] /filter/keyword/replace

| Name            | Type     | Remork                 |
|:----------------|:--------:|:-----------------------|
| keyword | string    | keyword for response |

response:
```
{
    "code": 0,
    "message":",
    "data":"中华**共和国"
}
```
