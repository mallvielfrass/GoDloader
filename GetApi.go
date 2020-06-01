package godloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup
var defaultLim = 5

//GetAPI получает готовую ссылку и имя файла от функций Get и GetAll, рассчитывает количество получившихся частей и загружает
func (api *StructAPI) GetAPI(url string, filename string) (int, string, int, bool) {
	//file := "http://localhost:8100/pinn-lite%281%29.zip"
	_, err := http.Head(url)
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf(" api.GetAPI:26 : called from %s, function %s line %d :error %s\n", details.Name(), fn, line, err)

		}
		return 0, "", 0, false
	}
	res, err := http.Head(url) // 187 MB file of random numbers per line
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", fn, line, err)
	}
	maps := res.Header
	length, _ := strconv.Atoi(maps["Content-Length"][0]) // Get the content length from the header request
	//fmt.Println(length)
	//	var a int = 646646464
	var b int = 32768
	var s int = length / b
	var limit int
	//проверка размера файла чтобы не резать на слишком мелкие куски
	switch s {
	case 0:
		limit = 1

	case 1:
		limit = 2

	case 2:
		limit = 3

	default:
		limit = defaultLim

	}
	//fmt.Println(limit)

	lenSub := length / limit   // Bytes for each Go-routine
	diff := length % limit     // Get the remaining for the last request
	body := make([]string, 11) // Make up a temporary array to hold the data to be written to the file
	for i := 0; i < limit; i++ {
		wg.Add(1)

		min := lenSub * i       // Min range
		max := lenSub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int, url string) {

			for {
				_, err := http.Head(url)
				if err != nil {
					pc, fn, line, ok := runtime.Caller(1)
					details := runtime.FuncForPC(pc)
					if ok && details != nil {
						fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

					}
					time.Sleep(5 * time.Second)
				} else {
					if api.Download(min, max, i, url, body, filename) == true {
						break
					}
				}
			}
			wg.Done()

		}(min, max, i, url)
	}

	wg.Wait()
	return limit, filename, length, true
}

//Download загружает часть
func (api *StructAPI) Download(min int, max int, i int, url string, body []string, filename string) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	rangeHeader := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
	req.Header.Add("Range", rangeHeader)
	resp, err := client.Do(req)
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

		}
		return false
	}
	defer resp.Body.Close()
	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

		}
		return false
	}
	body[i] = string(reader)
	fn := filename + ".tmp." + strconv.Itoa(i)
	api.Create(fn)
	file, err := os.OpenFile(fn, os.O_WRONLY, 0666)
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

		}
		return false
	}
	defer file.Close()
	_, err = file.Write([]byte(string(body[i])))
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)

		}
		return false
	}
	return true
}
