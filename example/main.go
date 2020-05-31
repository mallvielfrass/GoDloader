package main

import (
	"fmt"

	godloader "github.com/mallvielfrass/GoDloader"
)

func main() {
	ThreadLimit := 5
	api := godloader.Api(ThreadLimit)
	//api.GetAll("http://localhost:8100/static/testolde/parse/pinn-lite%281%29.zip")
	lim, name := api.Get("https://sourceforge.net/projects/reactos/files/ReactOS/0.4.13/ReactOS-0.4.13-iso.zip")
	fmt.Printf("%d %s \n", lim, name)
	if api.Build(name, lim) == true {
		api.Remove(name, lim)
	}
}

//2d7a5278b84a4f5c1e75c86dd88405bf4e34df3e *pinn-lite(1).zip
