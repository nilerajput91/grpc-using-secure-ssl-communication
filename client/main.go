package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"grpc/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	server   = "expert-Inspiron-3541:8080"
	count    = 10
	bounded  bool
	minValue int64
	maxValue int64

	src = rand.NewSource(time.Now().Unix())
	r   = rand.New(src)
)

func main() {
	flag.StringVar(&server, "a", server, "gRPC server address host:port")
	flag.IntVar(&count, "c", count, "number of random request")
	flag.BoolVar(&bounded, "b", bounded, "whether the random number within a range ,if b is true ,then min and max must be specifed")
	flag.Int64Var(&minValue, "min", minValue, "minium random number in sequence")
	flag.Int64Var(&maxValue, "max", maxValue, "maximum random number in sequence")
	flag.Parse()

	if count < 1 {
		log.Fatal("count must be greater than 0")
	}

	if bounded {
		if minValue < 1 || maxValue < 1 {
			log.Fatal("Bounded 'b' is set to true,min and max must>0 ")
		}

		if minValue >= maxValue {
			log.Fatal("Bounded 'b' is set to true but min greater than or equals to max ")
		}
	}

	var opts []grpc.DialOption
	creds, _ := credentials.NewClientTLSFromFile("../cert.pem", "")

	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(server, opts...)

	if err != nil {
		log.Fatal(fmt.Errorf("unable to connect to the grpc Service :%v ", err))
	}

	defer conn.Close()

	testMathService(conn)
	testDataServiceRandom(conn)
	testDataServiceSum(conn)
}

func testDataServiceSum(conn *grpc.ClientConn) {
	client := model.NewDataServiceClient(conn)

	ctx := context.Background()

	stream, err := client.Sum(ctx)

	if err != nil {
		log.Fatalf("DataService.Sum()RPC failed:%v ", err)
	}

	for i := 0; i < count; i++ {
		v := r.Int63n(50)
		fmt.Printf("%2v) Sending %v\n", i, v)
		in := &model.SumRequest{Value: v}
		stream.Send(in)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("DataService.Sum() RPC failed: %v", err)
	}

	fmt.Printf("Total: %v\n", res.Total)
}

func testDataServiceRandom(conn *grpc.ClientConn) {
	client := model.NewDataServiceClient(conn)

	ctx := context.Background()
	in := &model.RandomRequest{

		Count1:   int32(count),
		Bounded:  bounded,
		MinValue: minValue,
		MaxValue: maxValue,
	}

	stream, err := client.Random(ctx, in)

	if err != nil {
		log.Fatalf("Dataservice.Random() RPC failed :%v", err)
	}

	v, err := stream.Recv()
	i := 1

	for err == nil {
		fmt.Printf("%2v.%v\n", i, v.Value)
		i++
		v, err = stream.Recv()

	}

	if err != io.EOF {
		fmt.Printf("Error to reading the random stream:%v\n", err)

	}

}

func testMathService(conn *grpc.ClientConn) {

	client := model.NewMyMathServiceClient(conn)
	ctx := context.Background()
	in := &model.MathRequest{Operand1: 12, Operand2: 6}

	//call Add on the client stud
	result, err := client.Add(ctx, in)

	if err != nil {
		log.Fatal(fmt.Errorf("Add rpc failed:%v", err))
	}

	fmt.Printf("Add(%v)=>%v\n", in, result)

	//call sub on the  client stub

	result, err = client.Sub(ctx, in)

	if err != nil {
		log.Fatal(fmt.Errorf("Sub rpc failed:%v", err))
	}

	fmt.Printf("Sub(%v)=>%v\n", in, result)

	//call Mul on the client stub

	result, err = client.Mul(ctx, in)

	if err != nil {
		log.Fatal(fmt.Errorf("Mul rpc failed:%v", err))
	}

	fmt.Printf("Mul(%v)=>%v\n", in, result)

	//Call Div on the client stub

	result, err = client.Div(ctx, in)

	if err != nil {
		log.Fatal(fmt.Errorf("Div rpc failed:%v", err))
	}

	fmt.Printf("Div(%v)=>%v\n", in, result)

	//call Mod on the client stub

	result, err = client.Mod(ctx, in)

	if err != nil {
		log.Fatal(fmt.Errorf("Mod rpc failed:%v", err))
	}

	fmt.Printf("Mod(%v)=>%v\n", in, result)

}
