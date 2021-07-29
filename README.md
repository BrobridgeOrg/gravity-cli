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
  accessKey   Gravity Subscriber's Access Key Manager
  help        Help about any command
  setConfig   Set Gravity CLI configuration

Flags:
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
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
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
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

### Create Access Key for subscriber
```
./gravity-cli accessKey --help 
Gravity Subscriber's Access Key Manager

Usage:
  gravity-cli accessKey [flags]
  gravity-cli accessKey [command]

Available Commands:
  create      Create Gravity Subscriber's Access Key
  delete      Delete Gravity Subscriber's Access Key
  list        List Gravity Subscriber's Access Key
  update      Update Gravity Subscriber's Access Key

Flags:
  -h, --help   help for accessKey

Global Flags:
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")

Use "gravity-cli accessKey [command] --help" for more information about a command.
```

```bash
./gravity-cli accessKey create --help 
Create Gravity Subscriber's Access Key

Usage:
  gravity-cli accessKey create [flags]

Flags:
  -k, --accessKey string   Specify client's accessKey
  -i, --appID string       Specify client's appID
  -h, --help               help for create
  -n, --name string        Specify client's accessKey name

Global Flags:
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
```

#### Example
```bash
./gravity-cli accessKey create -i subID1 -n subscriber1 -k subscriber1AccessKey

./gravity-cli accessKey list
+--------+-------------+
| APPID  |   APPNAME   |
+--------+-------------+
| subID1 | subscriber1 |
+--------+-------------+
Total: 1

./gravity-cli accessKey list --all
+--------+-------------+----------------------+
| APPID  |   APPNAME   |      ACCESSKEY       |
+--------+-------------+----------------------+
| subID1 | subscriber1 | subscriber1AccessKey |
+--------+-------------+----------------------+
Total: 1

```

---

### Update Access Key
```bash
./gravity-cli accessKey update --help 
Update Gravity Subscriber's Access Key

Usage:
  gravity-cli accessKey update [AppID] [flags]

Flags:
  -k, --accessKey string   Specify new accessKey
  -h, --help               help for update
  -n, --name string        Specify new appName

Global Flags:
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")
```

#### Example
```bash
./gravity-cli accessKey update subID1 -n subscriberUpdateName -k subscriber1UpdateKey

./gravity-cli accessKey list --all
+--------+----------------------+----------------------+
| APPID  |       APPNAME        |      ACCESSKEY       |
+--------+----------------------+----------------------+
| subID1 | subscriberUpdateName | subscriber1UpdateKey |
+--------+----------------------+----------------------+
Total: 1

```

---

### Delete Access Key
```bash
./gravity-cli accessKey delete --help 
Delete Gravity Subscriber's Access Key

Usage:
  gravity-cli accessKey delete [AppID] [flags]

Flags:
  -h, --help   help for delete

Global Flags:
  -c, --config string   Specify Gravity CLI config file (default "$HOME/.gravity/config.toml")

```

#### Example
```bash
./gravity-cli accessKey list --all
+--------+----------------------+----------------------+
| APPID  |       APPNAME        |      ACCESSKEY       |
+--------+----------------------+----------------------+
| subID1 | subscriberUpdateName | subscriber1UpdateKey |
+--------+----------------------+----------------------+
| subID2 | subscriber2          | subscriber2AccessKey |
+--------+----------------------+----------------------+
| subID3 | subscriber3          | subscriber3AccessKey |
+--------+----------------------+----------------------+
| subID4 | subscriber4          | subscriber4AccessKey |
+--------+----------------------+----------------------+
| subID5 | subscriber5          | subscriber5AccessKey |
+--------+----------------------+----------------------+
Total: 5

./gravity-cli accessKey delete subID2 subID3 subID4 subID5

./gravity-cli accessKey list --all
+--------+----------------------+----------------------+
| APPID  |       APPNAME        |      ACCESSKEY       |
+--------+----------------------+----------------------+
| subID1 | subscriberUpdateName | subscriber1UpdateKey |
+--------+----------------------+----------------------+
Total: 1


```
---

## Author
Copyright(c) 2020 JheSue <jhe@brobridge.com>  
Copyright(c) 2020 Fred Chien <fred@brobridge.com> 
Copyright(c) 2020 Dagin Wu <daginwu@brobridge.com> 
