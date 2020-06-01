package godloader

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

//Create создает файл вместе с указанным путем
func (api *StructAPI) Create(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("created ", p)
	return os.Create(p)
}

//ExtErr расширеннная версия error check с дебагом
func (api *StructAPI) ExtErr(err error) bool {
	pc, fn, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

	}
	return false
	//fmt.Println("fn", fn, " line ", line)
}

//StructAPI главная структура апи. можно указать максимальное количество потоков
type StructAPI struct {
	ThreadLimit    int
	BlockSizeLimit int
}

//API точка входа в lib
func API(ThreadLimit int, BlockSizeLimit int) *StructAPI {
	return &StructAPI{
		ThreadLimit:    ThreadLimit,
		BlockSizeLimit: BlockSizeLimit,
	}
}
