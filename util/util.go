package util

import "fmt"

func Panic(err error, action string){
	if err != nil{
		fmt.Printf("Error while %s: panicing!!!!", action)
		panic(err)
	}
}

func Error(err error, action string){
	if err != nil {
		fmt.Printf("Error while %s: %s", action, err)
	}
}