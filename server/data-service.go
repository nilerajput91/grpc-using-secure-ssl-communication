package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"grpc/model"
)

type (
	myDataService struct {
	}
)

var (
	src = rand.NewSource(time.Now().Unix())
	r   = rand.New(src)
)

func (m *myDataService) Sum(in model.DataService_SumServer) error {
	if m == nil {
		return fmt.Errorf("sum called on the nill object")
	}

	var total int64
	v, err := in.Recv()

	for err == nil {
		total += v.Value
		v, err = in.Recv()

	}

	if err != io.EOF {
		return fmt.Errorf("Unexpected err:%v", err)
	}

	in.SendAndClose(&model.SumResponse{Total: total})
	return nil
}

func (m *myDataService) Random(in *model.RandomRequest, out model.DataService_RandomServer) error {
	if m == nil {
		return fmt.Errorf("Random call on Nil Object")
	}

	if in == nil {
		return fmt.Errorf("Random called with invalid parameter value nil")
	}

	count := int(in.Count1)
	var v int64

	if in.Bounded {
		for i := 0; i < count; i++ {
			v = r.Int63n(in.MaxValue-in.MinValue) + in.MinValue
			out.Send(&model.RandomResponse{Value: v})
		}
		return nil
	}

	for i := 0; i < count; i++ {
		v := r.Int63()
		out.Send(&model.RandomResponse{Value: v})
	}

	return nil
}
