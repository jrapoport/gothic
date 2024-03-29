package root

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AdminClient is a grpc client for the admin service.
type AdminClient struct {
	admin.AdminClient
	conn *grpc.ClientConn
}

// Close closes the client connection.
func (client AdminClient) Close() error {
	return client.conn.Close()
}

// NewAdminClient returns a new grpc client for the admin service.
func NewAdminClient() (*AdminClient, error) {
	c := Config()
	if c.AdminAddress == "" {
		return nil, errors.New("admin rpc address required")
	}
	if c.RootPassword == "" {
		return nil, errors.New("admin root password required")
	}
	conn, err := newConnection(c.AdminAddress, c.RootPassword)
	if err != nil {
		return nil, err
	}
	client := admin.NewAdminClient(conn)
	return &AdminClient{client, conn}, nil
}

// newConnection returns an rpc client connection
func newConnection(address, pw string) (*grpc.ClientConn, error) {
	cred := insecure.NewCredentials()
	if adminCert != "" {
		// Create the client TLS credentials
		var err error
		cred, err = credentials.NewClientTLSFromFile(adminCert, "")
		if err != nil {
			return nil, fmt.Errorf("could not load admin tls cert: %s", err)
		}
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(cred),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx = metadata.AppendToOutgoingContext(ctx, rpc.RootPassword, pw)
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	}
	return grpc.Dial(address, opts...)
}
