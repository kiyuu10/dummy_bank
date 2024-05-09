package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) ExtractMetadata(ctx context.Context) *Metadata {
	metaData := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			metaData.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			metaData.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		metaData.ClientIP = p.Addr.String()
	}

	return metaData
}
