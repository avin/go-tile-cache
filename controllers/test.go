package controllers

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
	"log"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) Get() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	serversConfigFile := filepath.Join(dir, "conf", "config.json")

	file, err := ioutil.ReadFile(serversConfigFile)
	if err != nil {
		fmt.Println("Cannot open servers configuration file:", err)
		os.Exit(1)
	}

	c.Data["serversConfig"] = string(file)

	c.TplName = "test.tpl"
}