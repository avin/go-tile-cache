package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
	"time"
	"strconv"
	"go-tile-cache/models"
	"path/filepath"
)

type TileController struct {
	beego.Controller
}

//Return blank image
func returnBlankTile(c *TileController) {
	blankTileFile := filepath.Join("static", "img", "blank.png")
	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")
	http.ServeFile(c.Ctx.Output.Context.ResponseWriter, c.Ctx.Output.Context.Request, blankTileFile)
	c.StopRun()
}

func (c *TileController) Get() {

	config := models.GetConfig()

	var ttlSeconds int = 86400 //cache timeout in seconds

	var x string = c.GetString("x");
	var y string = c.GetString("y");
	var z string = c.GetString("z");
	var gs bool = false; //is grayScale
	gs, _ = strconv.ParseBool(c.GetString("gs"));

	var server = c.GetString("server"); //tile server

	//create tile manager
	var tileManager = models.TileManager{X: x, Y: y, Z: z, Server: server, GS: gs}

	//get tile file
	fileName, err := tileManager.Get();
	//If cannot get tile file
	if (err != nil) {
		returnBlankTile(c)
	}

	//Setup headers
	c.Ctx.Output.Header("Expires", time.Now().Add(config.TtlDuration).Format("Sat, 15 Oct 2016 19:31:22 GMT"))
	c.Ctx.Output.Header("Cache-Control", "public, max-age=" + strconv.Itoa(ttlSeconds * 60));
	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")

	//response file
	http.ServeFile(c.Ctx.Output.Context.ResponseWriter, c.Ctx.Output.Context.Request, fileName)
}
