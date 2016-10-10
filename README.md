vimlex
======
A vim log lexer

example
=======
Quick pattern stats:

```shell
  cat /path/to/vim/logs* | vimlex | sort | uniq -c | sort -nr | head -n7
  2835 k
  2827 j
  1307 l
   821 h
   152 w
   136 :
   114 d
   ...
```

setup
=====
Track vim keystrokes by creating an alias that uses the scriptout flag:

```shell
  alias vim="vim -w ~/vimlogs/\$(date '+%Y%m%d').$RANDOM.log"
```

This will write keystrokes to the supplied file when vim is closed.

install
=======
```shell
$ go get github.com/dstokes/vimlex
```

Make sure your `PATH` includes your `GOPATH` bin directory:

```shell
export PATH=$PATH:$GOPATH/bin
```
