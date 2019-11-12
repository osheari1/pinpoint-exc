# Search engine
#### Steps to run
* Build app 
```bash
./scripts/setup.sh
```
* Start server
```
go run main.go
```
* Send requests 
```bash
>>> ./scripts/start_cli.sh
micro> call search Indexer.Index {"url": "http://<ENTER URL>"}
micro> call search Indexer.Search {"word": "<ENTER WORD>"}
```



