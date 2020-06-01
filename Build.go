package godloader

import (
	"fmt"
	"os"
	"strconv"
)

//Build собирает все части файла в один
func (api *StructAPI) Build(name string, lim int) bool {
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

		// читаем часть
		bs := make([]byte, stat.Size())
		_, err = file.Read(bs)
		if err != nil {
			return false
		}
		api.Stick(name, bs) //приклеиваем к целевому файлу
	}
	return true //если сборка прошла успешно, то возвращается true, иначе вернется false и можно будет обработать эту ошибку
}
