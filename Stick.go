package godloader

import (
	"fmt"
	"log"
	"os"
)

//Stick приклеивает к файлу набор переданных байтов
func (api *StructAPI) Stick(pathFileLog string, insertData []byte) {
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
	_, err = fileLog.Write(insertData)
	if err != nil {
		log.Fatalln(err)
	}
	//	log.Println("log countN:", countN)

}
