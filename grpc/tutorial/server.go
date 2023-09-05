package tutorial

import (
	context "context"
	"fmt"
)

type Server struct {
	UnimplementedTutorialServer
}

func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	fmt.Println("hi")
	return &HelloReply{Message: "Hello, " + in.GetName()}, nil
}
