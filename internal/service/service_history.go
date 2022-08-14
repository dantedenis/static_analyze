package service

type HistoryClient struct {
}

func New() *HistoryClient {
	return &HistoryClient{}
}

/*
func (h HistoryClient) GetHistory(ctx context.Context, in *proto.RequestMessage, opts ...grpc.CallOption) (*proto.ResponseMessage, error) {
	t1, err := time.Parse(time.RFC3339, "")
	t2, err := time.Parse(time.RFC3339, "")

	pbFrom := timestamppb.New(t1)
	pbTO := timestamppb.New(t2)

}
*/
