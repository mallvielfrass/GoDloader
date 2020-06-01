package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	godloader "github.com/mallvielfrass/GoDloader"
)

func main() {

	ThreadLimit := 8
	BlockSizeLimit := 5242880 //5MB one part
	api := godloader.API(ThreadLimit, BlockSizeLimit)
	//api.GetAll("http://localhost:8100/static/testolde/parse/pinn-lite%281%29.zip")
	url := "http://localhost:8100/static/ReactOS/ReactOS-0.3.17-REL-iso.zip" //run ../SeverForTest/server.go on port 8100
	s := strings.Split(url, "/")
	fmt.Printf("File ReactOS-0.4.13-iso.zip  %t was deleted\n", api.RemoveSimple("./"+s[len(s)-1]))
	//md5File := "3cc4988d6536e53b48eda5736147a457"
	md5File := "" //если имеется md5, используйте его. если он пустой, то будет проверяться длинна файла
	lim, name, length, status := api.Get(url)
	if status != false {
		fmt.Printf("%d %s \n", lim, name)
		api.Build(name, lim)
		if len(md5File) != 0 {
			hash, err := api.HashFileMd5(name)
			if err == nil {
				if hash == md5File {
					fmt.Printf("file %s has been downloaded and verified successfully\n", name)
					api.Remove(name, lim)
				}
			}
		} else {

			file, err := os.Open(name)
			if err != nil {
				log.Fatal(err)
			}
			fi, err := file.Stat()
			if err != nil {
				log.Fatal(err)
			}
			if fi.Size() == int64(length) {

				fmt.Printf("file %s has been downloaded and length checked successfully\n", name)
				api.Remove(name, lim)
			} else {
				fmt.Println("fi.size", fi.Size())
				fmt.Println("file size", int64(length))
				fmt.Printf("file %s hasn`t downloaded. lenght not checked\n", name)
			}
		}
	} else {

		fmt.Printf("file %s not downloaded. check internet", url)
	}
	//if api.Build(name, lim) == true {

	//}
}
