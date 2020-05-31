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
)

var wg sync.WaitGroup
var defaultLim = 5

func (api *API) Getx(url string, filename string) (int, string) {
	//file := "http://localhost:8100/pinn-lite%281%29.zip"
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
	//fmt.Println(s)
	switch s {
	case 0:
		limit = 1
	//	fmt.Println("limit := 1  ")
	case 1:
		limit = 2
	//	fmt.Println("limit := 2  ")
	case 2:
		limit = 3
	//	fmt.Println("limit := 3  ")

	default:
		limit = defaultLim
		//	fmt.Println("limit := 10  ")
	}
	fmt.Println(limit)
	//limit := 10                // 10 Go-routines for the process so each downloads 18.7MB
	len_sub := length / limit  // Bytes for each Go-routine
	diff := length % limit     // Get the remaining for the last request
	body := make([]string, 11) // Make up a temporary array to hold the data to be written to the file
	for i := 0; i < limit; i++ {
		wg.Add(1)

		min := len_sub * i       // Min range
		max := len_sub * (i + 1) // Max range

		if i == limit-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int, url string) {

			for {
				var Res bool
				Res = true
				_, err := http.Head(url)
				if err != nil {
					Res = false
					_, fn, line, _ := runtime.Caller(1)
					log.Printf("[error] %s:%d %v", fn, line, err)

				}
				client := &http.Client{}
				req, _ := http.NewRequest("GET", url, nil)
				range_header := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1) // Add the data for the Range header of the form "bytes=0-100"
				req.Header.Add("Range", range_header)
				resp, err := client.Do(req)
				if err != nil {
					Res = false
					_, fn, line, _ := runtime.Caller(1)
					log.Printf("[error] %s:%d %v", fn, line, err)
				}
				defer resp.Body.Close()
				reader, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					Res = false
					_, fn, line, _ := runtime.Caller(1)
					log.Printf("[error] %s:%d %v", fn, line, err)
				}
				body[i] = string(reader)
				fn := filename + ".tmp." + strconv.Itoa(i)
				//fmt.Println("fn ", fn)
				api.Create(fn)
				//ioutil.WriteFile(fn, []byte(string(body[i])), 0x777) // Write to the file i as a byte array
				file, err := os.OpenFile(fn, os.O_WRONLY, 0666)
				if err != nil {
					fmt.Println("Unable to open file:", err)
					os.Exit(1)
				}
				defer file.Close()
				_, err = file.Write([]byte(string(body[i])))
				if err != nil {
					Res = false
					_, fn, line, _ := runtime.Caller(1)
					log.Printf("[error] %s:%d %v", fn, line, err)
				}
				if Res == true {
					break
				}
				//fmt.Println("res: ", Res)
			}
			wg.Done()

		}(min, max, i, url)
	}
	wg.Wait()
	return limit, filename
}
