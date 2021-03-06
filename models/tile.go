package models

import (
	"os"
	"net/http"
	"io"
	"time"
	"path/filepath"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"github.com/astaxie/beego"
	"log"
	"gopkg.in/mgo.v2"
)

type TileManager struct {
	X      string
	Y      string
	Z      string
	Server string
	GS     bool
}

//download file
func downloadFile(fileName string, url string) (error) {

	//get url file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if (resp.StatusCode != 200) {
		return errors.New("error status code")
	}

	//create output file
	out, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()

	//write content to output file
	io.Copy(out, resp.Body)

	return nil
}

// NewTaskManager returns an empty TaskManager.
func NewTileManager() *TileManager {
	return &TileManager{}
}

//get file from cache or download (and convert to GS)
func (tm *TileManager) getTileUrl() (string, error) {

	var url string

	config := GetConfig()

	for _, serverItem := range config.Servers {
		if (serverItem.Alias == tm.Server) {

			url = serverItem.Url

			if (serverItem.Subdomains != nil) {
				subdomain := serverItem.Subdomains[rand.Intn(len(serverItem.Subdomains))]
				url = strings.Replace(url, "{s}", subdomain, -1)
			}

			url = strings.Replace(url, "{x}", tm.X, -1)
			url = strings.Replace(url, "{y}", tm.Y, -1)
			url = strings.Replace(url, "{z}", tm.Z, -1)

			break;
		}
	}

	if (len(url) == 0) {
		return "", errors.New("Server [" + tm.Server + "] not configured")
	}

	if (beego.AppConfig.String("runmode") == "dev") {
		fmt.Println(url)
	}

	return url, nil
}

//get file from cache or download (and convert to GS)
func (tm *TileManager) Get() (string, error) {

	type Person struct {
		Name  string
		Phone string
	}

	config := GetConfig()

	session, err := mgo.Dial(config.Mongodb.Host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	file, err := session.DB(config.Mongodb.DbName).GridFS("fs").Create("myfile.txt")
	_, err = file.Write([]byte("Hello world!"))
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return "ok.png", nil
}

func init() {

	config := GetConfig()

	if (config.ClearOldCache) {
		//Start clear old cache tile files ticker
		ticker := time.NewTicker(time.Minute * 30)
		go func() {
			for t := range ticker.C {
				fmt.Println("Tick at", t)

				visit := func(path string, f os.FileInfo, err error) error {

					if file, err := os.Stat(path); err == nil {
						//if file older then ttl duration
						if (time.Now().Add(0 - config.TtlDuration).After(file.ModTime())) {
							//delete file/path
							os.RemoveAll(path)
						}
					}

					return nil
				}

				dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
				if err != nil {
					log.Fatal(err)
				}
				path := filepath.Join(dir, "cache")
				err = filepath.Walk(path, visit)
				if (err != nil) {
					fmt.Printf("filepath.Walk() returned %v\n", err)
				}
			}
		}()
	}

}