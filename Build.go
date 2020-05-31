package godloader

import (
	"fmt"
	"os"
	"strconv"
)

func (api *API) Build(name string, lim int) bool {
	for i := 0; i < lim; i++ {
		file, err := os.Open(name + ".tmp." + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
			return false
			//	os.Exit(1)
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			return false
		}

		// чтение файла
		bs := make([]byte, stat.Size())
		_, err = file.Read(bs)
		if err != nil {
			return false
		}

		//	str := string(bs)
		//fmt.Println(str)
		api.Stick(name, bs)

	}
	return true
}
