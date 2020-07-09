package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type GRPCCommand struct {
	addr          string
	service       string
	tls           bool
	tlsNoVerify   bool
	tlsCACert     string
	tlsClientCert string
	tlsClientKey  string
	tlsServerName string
}

func (c *GRPCCommand) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("grpc", flag.ExitOnError)
	fs.StringVar(&c.addr, "addr", c.addr, "address (required)")
	fs.StringVar(&c.addr, "service", c.addr, "service name to check")
	fs.BoolVar(&c.tls, "tls", c.tls, "use TLS")
	fs.BoolVar(&c.tlsNoVerify, "tls-no-verify", c.tlsNoVerify, "do not verify the certificate")
	fs.StringVar(&c.tlsCACert, "tls-ca-cert", c.tlsCACert, "trusted certificates for verifying server")
	fs.StringVar(&c.tlsClientCert, "tls-client-cert", c.tlsClientCert, "")
	fs.StringVar(&c.tlsClientKey, "tls-client-key", c.tlsClientKey, "")
	fs.StringVar(&c.tlsServerName, "tls-server-name", c.tlsServerName, "")
	return fs
}

func (c *GRPCCommand) Run(args []string) error {
	if c.addr == "" {
		return flag.ErrHelp
	}
	var opts []grpc.DialOption
	if c.tls {
		var tlsConfig tls.Config
		if c.tlsClientCert != "" && c.tlsClientKey != "" {
			keyPair, err := tls.LoadX509KeyPair(c.tlsClientCert, c.tlsClientKey)
			if err != nil {
				return fmt.Errorf("load client key pair: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{keyPair}
		}
		if c.tlsNoVerify {
			tlsConfig.InsecureSkipVerify = true
		} else if c.tlsCACert != "" {
			certPool := x509.NewCertPool()
			pem, err := ioutil.ReadFile(c.tlsCACert)
			if err != nil {
				return fmt.Errorf("read ca cert: %w", err)
			}
			if !certPool.AppendCertsFromPEM(pem) {
				return fmt.Errorf("not found valid certificates: %s", c.tlsCACert)
			}
			tlsConfig.RootCAs = certPool
		}
		if c.tlsServerName != "" {
			tlsConfig.ServerName = c.tlsServerName
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	ctx := CommandContext()
	for range Continue(ctx, interval) {
		permanent, err := c.Request(ctx, opts)
		if permanent {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return ctx.Err()
}

func (c *GRPCCommand) Request(ctx context.Context, opts []grpc.DialOption) (bool, error) {
	conn, err := grpc.DialContext(ctx, c.addr, opts...)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	resp, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: c.service})
	if err != nil {
		if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
			return true, fmt.Errorf("unimplemented the grpc health check protocol")
		}
		return false, err
	}
	if resp.GetStatus() == grpc_health_v1.HealthCheckResponse_SERVING {
		return true, nil
	}
	return false, fmt.Errorf("unhealthy: %s", resp.GetStatus())
}
