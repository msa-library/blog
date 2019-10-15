package app

import (
	"log"
	"os"
	"fmt"
)

func AccesLog(msg string){
	filename:="/app/logs/access.log"
	writeLog(filename,msg)
}

func ErrLog(msg string){
	filename:="/app/logs/error.log"
	writeLog(filename,msg)
}

func writeLog(filename, msg string){
	f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[ERR] error opening file: %v \n",err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(msg)
}