Qmeta
==================

Qmeta is written in Go. It is an queue meta backup tool used for RabbitMQ. It depends on RabbitMQ management plugin's API, and backup metadata in the form of JSON format. You can recover the queue meta through management plugin when needed. 


##Usage

###Build
```
$ git clone git@github.com:fengchj/ftool.git
$ cd ftool/qmeta
$ go build 
```


###Use

```
$ qmeta -h
Usage of qmeta:
  -file="config": config file contains RabbitMQ node addrs.
$ cat config
127.0.0.1:15672 guest guest
10.20.30.40:15672 guest guest
$ qmeta -file config
Host 127.0.0.1 backup done!
Host 10.20.30.40 backup done!
$ cd qmeta_{time-in-second}
$ ls -l
127.0.0.1.json
10.20.30.40.json
```

