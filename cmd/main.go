package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"static-analyze/internal/app/proto"
	"time"
)

func main() {
	conn, err := grpc.Dial(":8888", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := proto.NewHistoryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	t1, err := time.Parse("2006-02-01 15:04:05", "2000-01-01 12:49:34")
	t2, err := time.Parse("2006-02-01 15:04:05", "2006-06-12 12:49:34")

	r, err := c.GetHistory(ctx, &proto.RequestMessage{
		Subject: "USDRUB",
		From:    timestamppb.New(t1),
		To:      timestamppb.New(t2),
	})

	if err != nil {
		log.Println(err)
	}
	log.Println(r.P)
}
