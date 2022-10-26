package indicator

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sort"
	"static-analyze/internal/app/proto"
	"sync"
	"time"
)

var (
	min        = int64(60)
	fiveMin    = int64(300)
	fifteenMin = int64(900)
	thirtyMin  = int64(1800)
	hour       = int64(3600)
)

type Get interface {
	Get(key string, group int64) ([]Indicator, time.Time, error)
}

type RowMapIndicator struct {
	sync.Mutex
	Indicators map[string]map[int64]ResponseIndicator
	client     proto.HistoryClient
	lastUpdate time.Time
}

// Get returns copy slice indicator
func (r *RowMapIndicator) Get(key string, group int64) ([]Indicator, time.Time, error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Indicators[key]; !ok {
		return nil, time.Time{}, fmt.Errorf("does not exist key-currency: %s", key)
	}

	if _, ok := r.Indicators[key][group]; !ok {
		return nil, time.Time{}, fmt.Errorf("does not exist key-group: %d", group)
	}

	temp := r.Indicators[key][group].Indic
	res := make([]Indicator, len(temp))

	copy(res, temp)
	return res, r.lastUpdate, nil
}

type ResponseIndicator struct {
	Indic []Indicator
}

func NewResponseIndicator(i []Indicator) ResponseIndicator {
	return ResponseIndicator{Indic: i}
}

func New(client proto.HistoryClient, pair []string) *RowMapIndicator {
	r := &RowMapIndicator{
		Indicators: map[string]map[int64]ResponseIndicator{},
		client:     client,
	}

	for _, val := range pair {
		r.Indicators[val] = make(map[int64]ResponseIndicator)
		r.Indicators[val][min] = ResponseIndicator{}
		r.Indicators[val][fiveMin] = ResponseIndicator{}
		r.Indicators[val][fifteenMin] = ResponseIndicator{}
		r.Indicators[val][thirtyMin] = ResponseIndicator{}
		r.Indicators[val][hour] = ResponseIndicator{}
	}

	return r
}

// Add to map new value in indicator
func (r *RowMapIndicator) Add(key string, i []Indicator, period int64) {
	r.Lock()
	defer r.Unlock()
	r.Indicators[key][period] = NewResponseIndicator(i)
}

// Updater method run update indicators from 00.00.00 to time of current day
func (r *RowMapIndicator) Updater() error {
	timeframe := []int64{min, fiveMin, fifteenMin, thirtyMin, hour}

	y, m, d := time.Now().Date()
	t1, err := time.Parse("2006 1 2", fmt.Sprintf("%d %d %d", y, m, d))
	if err != nil {
		return err
	}

	for key := range r.Indicators {
		resp, err := r.client.GetHistory(context.Background(), &proto.RequestMessage{
			Subject: key,
			From:    timestamppb.New(t1),
			To:      timestamppb.Now(),
		})
		if err != nil {
			log.Println(err)
			return err
		}

		for _, target := range timeframe {
			// goroutine for each timeframe
			go func(t int64, key string) {
				groupTime := make(map[int64][]*proto.Pair)
				for _, value := range resp.P {
					group := (value.GetTime().Seconds - t1.Unix()) / t
					groupTime[group] = append(groupTime[group], value)
				}
				var result []Indicator
				for group, value := range groupTime {
					var temp Indicator
					temp.Group = group
					temp.Open, temp.Close = value[0].Value, value[len(value)-1].Value
					temp.Low, temp.High = getMinMax(value)

					temp.Start, temp.End = periodToString(t1, temp.Group, t/60)
					result = append(result, temp)
				}
				sort.Slice(result, func(i, j int) bool {
					return result[i].Group < result[j].Group
				})

				r.Add(key, result, t)
			}(target, key)
		}
		r.lastUpdate = time.Now()
	}

	return err
}

func periodToString(t time.Time, group, c int64) (start, end string) {
	return time.Unix(t.Unix()+group*c, 0).String(), time.Unix(t.Unix()+(group+1)*c, 0).String()
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
