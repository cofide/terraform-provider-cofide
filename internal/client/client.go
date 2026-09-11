package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// jwtCredentials implements the grpc.PerRPCCredentials interface.
type jwtCredentials struct {
	token string
}

// GetRequestMetadata implements the grpc.PerRPCCredentials interface.
func (j *jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + j.token,
	}, nil
}

// RequireTransportSecurity implements the grpc.PerRPCCredentials interface.
func (j *jwtCredentials) RequireTransportSecurity() bool {
	return true
}

// NewTLSClient creates a new gPRC client with TLS credentials.
func NewTLSClient(connectAPIAddress, connectAPIServerName, jwtToken string, insecureSkipVerify bool, logger hclog.Logger, version string) (sdkclient.ClientSet, error) {
	tlsConfig, err := newTLSConfig(connectAPIServerName, insecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %v", err)
	}

	// We sometimes see auth errors such as 'Jwks remote fetch is failed', as
	// documented in https://github.com/cofide/cofide-connect/issues/223.
	// Use a retry policy to work around these.
	retryPolicy := `{
		"methodConfig": [{
		  "name": [{}],
		  "retryPolicy": {
			  "MaxAttempts": 10,
			  "InitialBackoff": ".1s",
			  "MaxBackoff": "1.0s",
			  "BackoffMultiplier": 2.0,
			  "RetryableStatusCodes": ["UNAUTHENTICATED"]
		  }
		}]}`

	opts := []grpc.DialOption{
		grpc.WithAuthority(connectAPIServerName),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithPerRPCCredentials(&jwtCredentials{token: jwtToken}),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithUserAgent(fmt.Sprintf("terraform-provider-cofide/%s", version)),
	}

	grpcConn, err := grpc.NewClient(connectAPIAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Connect gRPC server: %w", err)
	}

	logger.Info("Connecting to Connect gRPC server", "server", connectAPIAddress)

	return sdkclient.New(grpcConn), nil
}

// newTLSConfig creates a new TLS config based on the provided server name and insecure skip verify option.
func newTLSConfig(serverName string, insecureSkipVerify bool) (*tls.Config, error) {
	if insecureSkipVerify {
		return &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         serverName,
		}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("failed to load system root CA: %v", err)
	}

	return &tls.Config{
		RootCAs:    systemRoots,
		ServerName: serverName,
	}, nil
}
