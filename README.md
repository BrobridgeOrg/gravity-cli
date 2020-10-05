# gravity-cli

## How to Fast Run and Test

```
cd cmd/gravity-cli
go run main.go
```
## How to Build

```
cd cmd/gravity-cli
go build -o gravity-cli
```
## Command You can USE

Show all database manage by Gravity
```
gravity-cli store ls -a    
```

Search database manage by Gravity
```
gravity-cli store ls -d YOUR_DATABASE_NAME    
```

Recover all database manage by Gravity
```
gravity-cli store recover -a    
```

Search database and recover it that manage by Gravity
```
gravity-cli store recover -d YOUR_DATABASE_NAME  
``` 
## DEMO Screenshots
Main Page
![image](https://github.com/daginwu/demo-screenshot/blob/master/%5Bgravity-cli%5Dmainpage.png)
Show all database manage by Gravity
![image](https://github.com/daginwu/demo-screenshot/blob/master/%5Bgravity-cli%5Dls-all.png)
Search database manage by Gravity
![image](https://github.com/daginwu/demo-screenshot/blob/master/%5Bgravity-cli%5Dls-search.png)
Recover all database manage by Gravity
![image](https://github.com/daginwu/demo-screenshot/blob/master/%5Bgravity-cli%5Dls-all.png)
Search database and recover it that manage by Gravity
![image](https://github.com/daginwu/demo-screenshot/blob/master/%5Bgravity-cli%5Drecover-search-1.png)
## Author
Copyright(c) 2020 Dagin Wu <daginwu@brobridge.com>
