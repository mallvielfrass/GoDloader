package godloader

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

//Remove удаляет все части файла
func (api *StructAPI) Remove(name string, lim int) {
	for i := 0; i < lim; i++ {

		err := os.Remove(name + ".tmp." + strconv.Itoa(i))
		if err != nil {
			for { //проверка. если файл занят, то ждать две секунды пока он освободится
				_, fn, line, _ := runtime.Caller(1)
				log.Printf("[error] %s:%d %v", fn, line, err)
				time.Sleep(2 * time.Second)
				err := os.Remove(name + ".tmp." + strconv.Itoa(i))
				if err != nil {
					time.Sleep(2 * time.Second)
				} else {
					break
				}
				//return
			}
		}

	}
}

//RemoveSimple удаляет указанный файл
func (api *StructAPI) RemoveSimple(name string) bool {
	if Exists(name) == true {
		err := os.Remove(name)
		if err != nil {
			for { //проверка. если файл занят, то ждать две секунды пока он освободится
				_, fn, line, _ := runtime.Caller(1)
				log.Printf("[error] %s:%d %v", fn, line, err)
				time.Sleep(2 * time.Second)
				err := os.Remove(name)
				if err != nil {
					time.Sleep(2 * time.Second)
				} else {
					break
				}
				//return
			}
		}
		return true
	}
	return false

}

//Exists проверяет файл на существование
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
