package server

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"project/internal/generated"
	"project/internal/links"
)

type linksServiceGRPCServer struct {
	grpcServer   *grpc.Server
	linksUsecase links.Usecase
}

func NewLinksServiceGRPCServer(grpcServer *grpc.Server, linksUsecase links.Usecase) *linksServiceGRPCServer {
	return &linksServiceGRPCServer{
		grpcServer:   grpcServer,
		linksUsecase: linksUsecase,
	}
}

func (c *linksServiceGRPCServer) StartGRPCServer(listenURL string) error {
	lis, err := net.Listen("tcp", listenURL)
	if err != nil {
		return err
	}

	generated.RegisterLinksServer(c.grpcServer, c)

	return c.grpcServer.Serve(lis)
}

func (c *linksServiceGRPCServer) GetOriginalLink(ctx context.Context, link *generated.Link) (*generated.Link, error) {
	originalLink, err := c.linksUsecase.GetOriginalLink(ctx, link.Url)
	if err != nil {
		return nil, err
	}

	return &generated.Link{
		Url: originalLink.Url,
	}, nil
}

func (c *linksServiceGRPCServer) SaveAbbreviatedLink(ctx context.Context, link *generated.Link) (*generated.Link, error) {
	abbreviatedLink, err := c.linksUsecase.SaveAbbreviatedLink(ctx, link.Url)
	if err != nil {
		return nil, err
	}

	return &generated.Link{
		Url: abbreviatedLink.Url,
	}, nil
}
