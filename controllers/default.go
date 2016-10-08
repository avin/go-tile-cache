package controllers

import (
	"github.com/harrydb/go/img/grayscale"
	"image/png"
	"github.com/astaxie/beego"
	"math/rand"
	"net/http"
	"path/filepath"
	"io"
	"os"
	"time"
	"strconv"
)

type TileController struct {
	beego.Controller
}

func GrayScale(filename string) string {

	//read file
	infile, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer infile.Close()

	//decode PNG
	src, err := png.Decode(infile)
	if err != nil {
		panic(err.Error())
	}

	//convert
	gray := grayscale.Convert(src, grayscale.ToGrayLuminance)

	//get name for converted file
	outFileName := filename[:len(filename) - 4] + ".gs.png"

	if _, err := os.Stat(outFileName); err == nil {
		//if file exist - return it now
		return outFileName
	}

	//save converted file
	outfile, err := os.Create(outFileName)
	if err != nil {
		panic(err.Error())
	}
	defer outfile.Close()
	png.Encode(outfile, gray)

	return outFileName
}

//download file
func DownloadFile(fileName string, url string) {
	//create output file
	out, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()

	//get url file
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	//write content to destination file
	io.Copy(out, resp.Body)
}

//get file from cache or download (and convert to GS)
func GetTile(url string, x string, y string, z string, server string, gs bool, ttl time.Duration) (string, error) {

	//make path
	path := filepath.Join("cache", server, x, y)
	os.MkdirAll(path, os.ModePerm)
	fileName := filepath.Join(path, z + ".png")

	if file, err := os.Stat(fileName); err != nil {
		//if file not exist
		DownloadFile(fileName, url)
	} else {
		//if file older then ttl duration
		if (time.Now().Add(-ttl).After(file.ModTime())) {
			DownloadFile(fileName, url)
		}
	}

	//if we gonna get grayScale tile - convert it
	if (gs) {
		gsFileName := GrayScale(fileName)
		return gsFileName, nil
	}

	return fileName, nil

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
		url = "http://" + randomServer + "/" + x + "/" + y + "/" + z + ".png";
	case "yandex":
		tileServers = append(tileServers, "vec01.maps.yandex.net")
		tileServers = append(tileServers, "vec02.maps.yandex.net")
		tileServers = append(tileServers, "vec03.maps.yandex.net")
		tileServers = append(tileServers, "vec04.maps.yandex.net")

		randomServer := tileServers[rand.Intn(len(tileServers))]
		url = "http://" + randomServer + "/tiles?l=map&v=4.113.1&x=" + x + "&y=" + y + "&z=" + z + "&scale=1&lang=ru_RU";
	default:
		c.Ctx.Output.Body([]byte(""))
		c.StopRun()
	}

	//get tile file
	fileName, _ := GetTile(url, x, y, z, server, gs, ttl);

	//Setup headers
	c.Ctx.Output.Header("Expires", time.Now().Add(ttl).Format("Sat, 15 Oct 2016 19:31:22 GMT"))
	c.Ctx.Output.Header("Cache-Control", "public, max-age=" + strconv.Itoa(ttlSeconds * 60));
	c.Ctx.Output.Header("Content-Type", "image/png")
	c.Ctx.Output.Header("Content-Transfer-Encoding", "binary")

	//response file
	http.ServeFile(c.Ctx.Output.Context.ResponseWriter, c.Ctx.Output.Context.Request, fileName)
}
