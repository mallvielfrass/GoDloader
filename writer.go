package godloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (api *API) Create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("created ", p)
	return os.Create(p)
}
func (api *API) Stick(pathFileLog string, insertData []byte) {
	var fileLog *os.File
	var err error
	_, err = os.Stat(pathFileLog)
	if os.IsNotExist(err) {
		api.Create(pathFileLog)
	}
	fileLog, err = os.OpenFile(pathFileLog, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = fileLog.Close()
		if err != nil {
			log.Fatalln(err)
			fmt.Println("Unable to create file:")
		}
	}()
	//	insertData := []byte("test")
	countN, err := fileLog.Write(insertData)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("log countN:", countN)

}

func (api *API) Get(url string) (int, string) {
	//file := "http://localhost:8100/pinn-lite%281%29.zip"
	//stick("ff/ffis", []byte("lollll"))
	//	original :="http://localhost:8100/pinn-lite%281%29.zip"
	s := strings.Split(url, "/")
	// fmt.Println(len(s))
	num := len(s)
	//fmt.Println(s[num-1])
	limit, filename := api.Getx(url, "./"+s[num-1])
	return limit, filename
}
func (api *API) GetAll(url string) (int, string) {
	//file := "http://localhost:8100/pinn-lite%281%29.zip"
	//stick("ff/ffis", []byte("lollll"))
	s := strings.Split(url, "/")
	// fmt.Println(len(s))
	num := len(s)
	str := ""
	fmt.Println(num)
	for i := 3; i < num; i++ {
		str = str + "/" + s[i]
	}
	fmt.Println(s[num-1])
	fmt.Println(str)
	limit, filename := api.Getx(url, "./"+str)
	return limit, filename
}

type API struct {
	ThreadLimit int
}

func Api(ThreadLimit int) *API {
	return &API{
		ThreadLimit: ThreadLimit,
	}
}
