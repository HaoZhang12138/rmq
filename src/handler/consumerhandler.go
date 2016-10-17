package handler

import (
	"dao"
	"gopkg.in/mgo.v2/bson"
	"log"
	"fmt"
	"time"
)

func Handleafterconsume(body []byte) (err error){
	res := new(dao.RechargeInfo)
	err = bson.Unmarshal(body, res)
	if err != nil {
		log.Println("failed to unmarshal to RechargeInfo")
	}
	fmt.Printf("%v:%v\n", time.Now(), res)
	//fmt.Println(res)
	return
}
