package controllers

import (
	"github.com/astaxie/beego"
	"math/rand"
	"net/http"
	"time"
	"strconv"
	"go-tile-cache/models"
	"path/filepath"
	"fmt"
)

type TileController struct {
	beego.Controller
}

//Return blank image
func returnBlankTile(c *TileController){
	blankTileFile := filepath.Join("static", "img", "blank.png")
	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")
	http.ServeFile(c.Ctx.Output.Context.ResponseWriter, c.Ctx.Output.Context.Request, blankTileFile)
	c.StopRun()
}

func (c *TileController) Get() {

	var ttlSeconds int = 86400 //cache timeout in seconds

	var ttl time.Duration = time.Duration(time.Duration(ttlSeconds) * time.Second); //cache timeout in seconds
	var x string = c.GetString("x");
	var y string = c.GetString("y");
	var z string = c.GetString("z");
	var gs bool = false; //is grayScale
	gs, _ = strconv.ParseBool(c.GetString("gs"));

	var server = c.GetString("server"); //tile server

	var tileServers []string
	var url string

	switch server {
	case "mapnik":
		tileServers = append(tileServers, "a.tile.openstreetmap.org")
		tileServers = append(tileServers, "b.tile.openstreetmap.org")
		tileServers = append(tileServers, "c.tile.openstreetmap.org")

		randomServer := tileServers[rand.Intn(len(tileServers))]
		url = "http://" + randomServer + "/" + z + "/" + x + "/" + y + ".png";
		fmt.Println(url)

	case "yandex":
		tileServers = append(tileServers, "vec01.maps.yandex.net")
		tileServers = append(tileServers, "vec02.maps.yandex.net")
		tileServers = append(tileServers, "vec03.maps.yandex.net")
		tileServers = append(tileServers, "vec04.maps.yandex.net")

		randomServer := tileServers[rand.Intn(len(tileServers))]
		url = "http://" + randomServer + "/tiles?l=map&v=4.113.1&x=" + x + "&y=" + y + "&z=" + z + "&scale=1&lang=ru_RU";
	default:
		returnBlankTile(c)
	}

	//create tile manager
	var tileManager = models.TileManager{X: x, Y: y, Z: z, Server: server, GS: gs}

	//get tile file
	fileName, err := tileManager.Get(url, ttl);
	//If cannot get tile file
	if (err != nil){
		returnBlankTile(c)
	}

	//Setup headers
	c.Ctx.Output.Header("Expires", time.Now().Add(ttl).Format("Sat, 15 Oct 2016 19:31:22 GMT"))
	c.Ctx.Output.Header("Cache-Control", "public, max-age=" + strconv.Itoa(ttlSeconds * 60));
	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")

	//response file
	http.ServeFile(c.Ctx.Output.Context.ResponseWriter, c.Ctx.Output.Context.Request, fileName)
}
