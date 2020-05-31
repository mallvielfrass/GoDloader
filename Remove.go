package godloader

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

func (api *API) Remove(name string, lim int) {
	for i := 0; i < lim; i++ {

		err := os.Remove(name + ".tmp." + strconv.Itoa(i))

		if err != nil {
			for {
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
