package main

import "github.com/DaHuangQwQ/webook/ioc"

func main() {
	server := ioc.InitApp()
	err := server.Web.Start()
	if err != nil {
		return
	}
}
