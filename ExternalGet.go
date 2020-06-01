package godloader

import "strings"

//Get скачивает файл
func (api *StructAPI) Get(url string) (int, string, int, bool) {
	s := strings.Split(url, "/")
	num := len(s)
	limit, filename, length, status := api.GetAPI(url, "./"+s[num-1])
	return limit, filename, length, status
}

//GetCreateFolderTree скачивает файл, создавая структуру директорий как на сервере
func (api *StructAPI) GetCreateFolderTree(url string) (int, string, int, bool) {
	s := strings.Split(url, "/")
	str := ""
	for i := 3; i < len(s); i++ {
		str = str + "/" + s[i]
	}
	limit, filename, length, status := api.GetAPI(url, "./"+str)
	return limit, filename, length, status
}
