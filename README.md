# gopac [![Gobuild Download](http://beta.gobuild.io/badge/github.com/futurespace/gopac/download.png)](http://beta.gobuild.io/github.com/futurespace/gopac)

A PAC(Proxy auto-config) file generator with fetched China IP range,   
which helps walk around GFW.   
Forked from [Flora_Pac][] and ported [Flora_Pac][] to Golang.


## Usage:

```
$ gopac help
```

```sh
// First, generate PAC file.
$ gopac build
// Second, host the PAC file.
$ gopac serve -p 8970
```


[Flora_Pac]: https://github.com/Leask/Flora_Pac
