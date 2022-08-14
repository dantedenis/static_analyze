package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"sort"
	"static-analyze/internal/app/proto"
	"time"
)

type Indicator struct {
	Open   float32 `json:"open"`
	High   float32 `json:"high"`
	Low    float32 `json:"low"`
	Close  float32 `json:"close"`
	Group  int64   `json:"-"`
	Period string  `json:"period"`
}

func main() {
	conn, err := grpc.Dial(os.Getenv("NETWORK_NAME")+os.Getenv("RPC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := proto.NewHistoryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	y, m, d := time.Now().Date()
	t1, err := time.Parse("2006 1 2", fmt.Sprintf("%d %d %d", y, m, d))

	r, err := c.GetHistory(ctx, &proto.RequestMessage{
		Subject: "USDRUB",
		From:    timestamppb.New(t1),
		To:      timestamppb.Now(),
	})

	if err != nil {
		log.Println(err)
		return
	}

	mapFiveMin := make(map[int64][]*proto.Pair)
	for _, value := range r.P {
		group := (value.GetTime().Seconds - t1.Unix()) / 300
		mapFiveMin[group] = append(mapFiveMin[group], value)
	}

	var result []Indicator
	for group, value := range mapFiveMin {
		var temp Indicator
		temp.Group = group
		temp.Open, temp.Close = value[0].Value, value[len(value)-1].Value
		temp.Low, temp.High = getMinMax(value)
		temp.Period = fmt.Sprintf("%s...%s", time.Unix((temp.Group-1)*5, 0).String(), time.Unix((temp.Group)*5, 0).String())
		result = append(result, temp)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Group < result[j].Group
	})

	indent, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(indent))
}

func getMinMax(values []*proto.Pair) (min, max float32) {
	min = 1.0
	for _, curr := range values {
		if curr.Value < min {
			min = curr.Value
		}
		if curr.Value > max {
			max = curr.Value
		}
	}
	return
}
