# Sardines
> The sardine are a decentralized system, which allows them to survive and evolve in the cruel nature.

A decentralized storage network based on Golang
<br>
This project will consider the function of Filecoin and ipfs, adopt one of its cores,a P2P go-library called 
*<u>libp2p</u>*.This wonderful library will be the fundamental part of our project,moreover, the file uploading and downloading will use the ipfs-go-API.
<br>
***
The project will be advanced by an tech called `Searchable Encryption` to build a more reliable decentralized storage network.


## Function List
1. p2p Chat 
2. p2p Ping 
3. p2p JoinApply 
4. p2p RouterDistribution 
5. p2p Inverted index table is done, it used `goleveldb` package, as we know, leveldb is a high-performance key-value storage library, however, this package is written in Go.
6. CLI(Command-line Interface) Some necessary tool-chains 
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


## Configuration 
Two configs are needed for running the server,they should be stored in directory `./values`. One is `config.json`, the other is `priv_key`. <br>
You can use cli tool-chain to configure:
```shell
$ ./sardines init

$ ./sardines -u root - p root -P 8082 -r 823 -b /ip4/x.x.x.x/tcp/8082/p2p/Qm... conf

$ ./sardines gen-key
```
Also, all of them can be made by API:
```go
c, err := config.New()
if err != nil {
	return
}    
c.Save()
```

## Run
Once configured compeletly:
```shell
$ ./sardines run
```

## Contributors

