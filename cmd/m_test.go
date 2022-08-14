package main

import (
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//y, m, d := time.Now().Date()
	t1, err := time.Parse("2006-1-2", time.Now().String())
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(timestamppb.Now().AsTime(), timestamppb.New(t1).AsTime())
}

func Test1(t *testing.T) {
	y, m, d := time.Now().Date()
	t1, err := time.Parse("2006-1-2", fmt.Sprintf("%d-%d-%d", y, m, d))
	if err != nil {
		log.Fatalln(err)
	}
	//end := timestamppb.Now()
	//start := timestamppb.New(t1)

	log.Println(t1.Minute(), time.Now().Minute())

}
