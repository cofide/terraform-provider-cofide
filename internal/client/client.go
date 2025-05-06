package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	sdkclient "github.com/cofide/cofide-api-sdk/pkg/connect/client"
	"github.com/cofide/terraform-provider-cofide/internal/consts"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type jwtCredentials struct {
	token string
}

func (j *jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + j.token,
	}, nil
}

func (j *jwtCredentials) RequireTransportSecurity() bool {
	return true
}

func NewTLSClient(baseAddr string, jwtToken string, devMode bool, logger hclog.Logger) (sdkclient.ClientSet, error) {
	serverName, err := getServerName(baseAddr)
	if err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if !devMode {
		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to load system root CA: %v", err)
		}
		tlsConfig = &tls.Config{
			RootCAs:    systemRoots,
			ServerName: serverName,
		}
	} else {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         serverName,
		}
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
		grpc.WithAuthority(serverName),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithPerRPCCredentials(&jwtCredentials{token: jwtToken}),
		grpc.WithDefaultServiceConfig(retryPolicy),
	}

	connectUri := fmt.Sprintf("dns:///%s.%s", consts.ServerAuthoritySubdomain, baseAddr)

	grpcConn, err := grpc.NewClient(connectUri, opts...)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Connect gRPC server: %w", err)
	}

	logger.Info("Connecting to Connect instance", "server", connectUri)

	return sdkclient.New(grpcConn), nil
}

func getServerName(baseAddr string) (string, error) {
	serverHost, _, err := net.SplitHostPort(baseAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", consts.ServerAuthoritySubdomain, serverHost), nil
}
