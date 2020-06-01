package godloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/remeh/sizedwaitgroup"
)

//var wg sync.WaitGroup

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
	var b int
	var MinBlockSizeLimit int = 32768
	if api.BlockSizeLimit < MinBlockSizeLimit {
		b = MinBlockSizeLimit
	} else {

		b = api.BlockSizeLimit
	}
	//	var s int = length / b

	lenSub := length / b           // Bytes for each Go-routine
	diff := length % b             // Get the remaining for the last request
	body := make([]string, lenSub) // Make up a temporary array to hold the data to be written to the file
	var DownloadThreads int
	//	fmt.Printf("block size  : %d\n ", b)
	//	fmt.Printf("lensub size  : %d\n ", lenSub)
	//	if lenSub < api.ThreadLimit {

	//		DownloadThreads = lenSub

	//	} else {
	DownloadThreads = api.ThreadLimit

	//	}
	//fmt.Printf("downloads th : %d\n ", DownloadThreads)
	var wg = sizedwaitgroup.New(DownloadThreads)
	th := lenSub
	for i := 0; i < th; i++ {
		wg.Add()

		min := b * i       // Min range
		max := b * (i + 1) // Max range
		//	fmt.Printf("%d %d %d\n", i, min, max)
		if i == th-1 {
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
	return th, filename, length, true
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
	fs, err := api.Create(fn)
	if err != nil {
		pc, fn, line, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n, function %s line %d :error %s", details.Name(), fn, line, err)
		}
		return false
	}
	defer fs.Close()
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
