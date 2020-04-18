package client

import (
	"github.com/akula410/web/server"
	"context"
	"flag"
	"google.golang.org/grpc"
	"fmt"
)

func Listen(){
	flag.Parse()
	if flag.NArg() > 0 {
		var conn *grpc.ClientConn
		var connProto server.ApiClient
		var r *server.Response
		var err error

		conn, err = grpc.Dial(":8081", grpc.WithInsecure())

		if err != nil {
			panic(err)
		}

		connProto = proto.NewApiClient(conn)

		r, err = connProto.Add(context.Background(), &server.Request{Command: flag.Arg(0)})

		if err != nil {
			panic(err)
		}

		if r.Message != nil {
			for _, m := range r.Message{
				fmt.Println(m)
			}
		}else{
			fmt.Println(r.Result)
		}
	}else{
		fmt.Println("Not enough arguments")
	}
}
