# kilit

kills an open file running on a port...
- clone into your gopath
- run ```go install```
- use in place of ```lsof -i :<port>``` and ```kill -9 <pid>``` combined
- change the port with the ```-p``` flag eg. ```kilit -p=8082``` (defaults to 8080)
