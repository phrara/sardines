# decentralodge
A decentralized storage network based on Golang
<br>
This project will consider the function of IPFS and adopt one of its cores,a P2P go-library called 
*<u>libp2p</u>*.This wonderful library will be the fundamental part of our project,however,writer's
capability is quite poor,so the usage of the lib seem stupid and unfamiliar.
<br>
***
The progress of the project would be slow,writer sincerely hope it could
make great progress one day.


### Function List
p2p Chat is done
<br>
p2p Ping is done
<br>
p2p JoinApply is done 
<br>
p2p RouterDistribution is done
<br>
p2p Inverted index table is done, it used `goleveldb` package, as we know, leveldb is a high-performance key-value storage library, however, this package is written in Go.
```go
hn, _ := core.GenerateNode()
{
	hn.JoinNetwork()
	hn.RouterDistributeOn(false, 10)
	// hn.RouterDistributeOn(true, 0)
	hn.RouterDistributeOff()
}
select{}
```


### Configuration 
Two configs are needed for running the server,they should be
stored in directory /root/values. One is `config.json`, the
other is `priv_key`, both of them can be made by API:
```go
c, err := config.New()
if err != nil {
	return
}    
c.Save()
```

> author phrara