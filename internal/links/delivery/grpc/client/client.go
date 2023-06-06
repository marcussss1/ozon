package client

import (
	"context"
	"google.golang.org/grpc"
	"project/internal/generated"
	"project/internal/links"
	"project/internal/model"
)

type linksServiceGRPCClient struct {
	linksClient generated.LinksClient
}

func NewLinksServiceGRPSClient(con *grpc.ClientConn) links.Usecase {
	return &linksServiceGRPCClient{
		linksClient: generated.NewLinksClient(con),
	}
}

func (l linksServiceGRPCClient) GetOriginalLink(ctx context.Context, url string) (model.Link, error) {
	originalLink, err := l.linksClient.GetOriginalLink(ctx, &generated.Link{
		Url: url,
	})

	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url: originalLink.Url,
	}, nil
}

func (l linksServiceGRPCClient) SaveAbbreviatedLink(ctx context.Context, url string) (model.Link, error) {
	originalUrl, err := l.linksClient.SaveAbbreviatedLink(ctx, &generated.Link{
		Url: url,
	})

	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url: originalUrl.Url,
	}, nil
}
