# GO-Tile-Cache

Tile proxy caching server with gray-scale image converting ability.

## Build

```
go get github.com/astaxie/beego
go get github.com/beego/bee
go get github.com/harrydb/go/img/grayscale
bee pack
```

## Usage

Configure tiles servers in `conf/config.json` file. Use new tile server in your map.
```
http://SERVER-HOST:PORT/?x=310&y=158&z=9&server=SERVER_ALIAS&gs=1
```