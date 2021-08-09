# Gravity CLI

## Build

``` bash
go build ./cmd/gravity-cli/gravity-cli.go
```

---

## Usage
```
./gravity-cli  --help 
Gravity CLI tool

If $GRACONFIG environment variable is set, then that config file is loaded.

Usage:
  gravity-cli [command]

Available Commands:
  accessKey   Gravity Access Key Manager
  adapter     Gravity Adapter Manager
  help        Help about any command
  setConfig   Set Gravity CLI configuration
  subscriber  Gravity Subscriber Manager

Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
  -h, --help            help for gravity-cli

Use "gravity-cli [command] --help" for more information about a command.

```

---

### Set connect info
``` bash
./gravity-cli setConfig --help 
Set Gravity CLI configuration and write config file to Gravity CLI config file

Usage:
  gravity-cli setConfig [flags]

Flags:
  -k, --accessKey string   Specify Gravity AccessKey.
  -i, --appID string       Specify Application ID
  -d, --domain string      Specify Gravity Domain (default "gravity")
  -h, --help               help for setConfig
  -H, --host string        Specify Gravity host:port (example "127.0.0.1:4222")

Global Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
```

#### Example
``` bash
./gravity-cli setConfig -H 172.17.0.1:4222 -d gravity -i gravity -k GRAVITYaccessKEY

cat ~/.gravity/config.toml

[gravity]
  accesskey = "GRAVITYaccessKEY"
  appid = "gravity"
  domain = "gravity"
  host = "172.17.0.1:4222"
```

---

### Create Access Key
```
./gravity-cli accessKey --help 
Gravity Access Key Manager

Usage:
  gravity-cli accessKey [flags]
  gravity-cli accessKey [command]

Available Commands:
  create      Create Gravity Access Key
  delete      Delete Gravity Access Key
  list        List Gravity Access Key
  update      Update Gravity Access Key

Flags:
  -h, --help   help for accessKey

Global Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")

Use "gravity-cli accessKey [command] --help" for more information about a command.
```

```bash
./gravity-cli accessKey create --help 
Create Gravity Access Key

Usage:
  gravity-cli accessKey create [flags]

Flags:
  -k, --accessKey string     Specify client's accessKey
  -i, --appID string         Specify client's appID
  -c, --collections string   Assign accessKey can subscribe's collections, This flag can using "," to  specified multiple collections. (Default can subscribe any collections.)
  -h, --help                 help for create
  -n, --name string          Specify client's accessKey name
  -r, --roles string         Specify accessKey's roles [ SYSTEM | ADAPTER | SUBSCRBIER ], This flag can using "," to  specified multiple roles.

Global Flags:
      --config string   Specify Gravity CLI config file (default "/home/jhe/.gravity/config.toml")

```

#### Example
```bash
./gravity-cli accessKey create -i subID1 -n subscriber1 -k subscriber1AccessKey -r SUBSCRIBER -c accountData

./gravity-cli accessKey list
APPID 	APPNAME     
subID1	subscriber1	

./gravity-cli accessKey list --all
APPID 	APPNAME    	ACCESSKEY           	ROLES     	COLLECTIONS 
subID1	subscriber1	subscriber1AccessKey	SUBSCRIBER	accountData

```

---

### Update Access Key
```bash
./gravity-cli accessKey update --help 
Update Gravity Access Key

Usage:
  gravity-cli accessKey update [AppID] [flags]

Flags:
  -k, --accessKey string     Specify new accessKey
  -c, --collections string   Assign accessKey can subscribe's collections, This flag can using "," to  specified multiple collections. (Default can subscribe any collections.)
  -h, --help                 help for update
  -n, --name string          Specify new appName
  -r, --roles string         Specify accessKey's roles [ SYSTEM | ADAPTER | SUBSCRBIER ], This flag can using "," to  specified multiple roles.

Global Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
```

#### Example
```bash
./gravity-cli accessKey update subID1 -n subscriberUpdateName -k subscriber1UpdateKey -r ADAPTER,SUBSCRIBER

./gravity-cli accessKey list --all
APPID 	APPNAME             	ACCESSKEY           	ROLES             	COLLECTIONS 
subID1	subscriberUpdateName	subscriber1UpdateKey	ADAPTER,SUBSCRIBER	accountData
```

---

### Delete Access Key
```bash
./gravity-cli accessKey delete --help 
Delete Gravity Access Key

Usage:
  gravity-cli accessKey delete [AppID] [flags]

Flags:
  -h, --help   help for delete

Global Flags:
      --config string   Specify Gravity CLI config file (default "/home/jhe/.gravity/config.toml")

```

#### Example
```bash
./gravity-cli accessKey list --all
APPID 	APPNAME             	ACCESSKEY           	ROLES             	COLLECTIONS 
subID1	subscriberUpdateName	subscriber1UpdateKey	ADAPTER,SUBSCRIBER	accountData	
subID2	subscriber2         	subscriber2AccessKey	SUBSCRIBER        	userData   	
subID3	subscriber3         	subscriber3AccessKey	ADAPTER           	           	


./gravity-cli accessKey delete subID1 subID2

./gravity-cli accessKey list --all
APPID 	APPNAME    	ACCESSKEY           	ROLES  	COLLECTIONS 
subID3	subscriber3	subscriber3AccessKey	ADAPTER	           

```

---

### List Adapter
```bash
./gravity-cli adapter --help 
Gravity Adapter Manager

Usage:
  gravity-cli adapter [flags]
  gravity-cli adapter [command]

Available Commands:
  list        List Gravity Adapters

Flags:
  -h, --help   help for adapter

Global Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")

Use "gravity-cli adapter [command] --help" for more information about a command.
```

#### Example
```bash
./gravity-cli adapter list
ID              	NAME            	COMPONENT 
postgres_adapter	Postgres Adapter	postgres
```
---

### List Subscriber
```bash
./gravity-cli subscriber --help 
Gravity Subscriber Manager

Usage:
  gravity-cli subscriber [flags]
  gravity-cli subscriber [command]

Available Commands:
  list        List Gravity Subscribers

Flags:
  -h, --help   help for subscriber

Global Flags:
      --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")

Use "gravity-cli subscriber [command] --help" for more information about a command.
```

#### Example
```bash
./gravity-cli subscriber list
ID                  	NAME                	COMPONENT	TYPE        
postgres_transmitter	Postgres Transmitter	postgres 	TRANSMITTER
```

---

## Author
Copyright(c) 2020 JheSue <jhe@brobridge.com>  
Copyright(c) 2020 Fred Chien <fred@brobridge.com>  
Copyright(c) 2020 Dagin Wu <daginwu@brobridge.com> 
