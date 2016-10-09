package controllers

import (
	"github.com/astaxie/beego"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"os"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) Get() {

	serversConfigFile := filepath.Join("conf", "config.json")

	file, err := ioutil.ReadFile(serversConfigFile)
	if err != nil {
		fmt.Println("Cannot open servers configuration file:", err)
		os.Exit(1)
	}

	c.Data["serversConfig"] = string(file)

	c.TplName = "test.tpl"
}