package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"static-analyze/internal/app/api"
	"static-analyze/internal/app/proto"
	indic "static-analyze/pkg/indicator"
	"strconv"
	"strings"
	"time"
)

var Result *indic.RowMapIndicator
var Pair []string

func main() {
	Pair = strings.Split(os.Getenv("CUR_PAIR"), ",")
	timer, err := strconv.Atoi(os.Getenv("PERIOD_SEC"))
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := grpc.Dial(os.Getenv("NETWORK_NAME")+os.Getenv("RPC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	client := proto.NewHistoryClient(conn)
	Result = indic.New(client, Pair)

	go func() {
		for {
			log.Println("update data")
			err = Result.Updater()
			if err != nil {
				log.Fatalln(err)
			}
			<-time.After(time.Duration(timer) * time.Second)
		}
	}()

	serv := api.NewServer(Result)
	err = serv.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
