package grpc

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata  = status.Errorf(codes.InvalidArgument, "missing metadata")
	ErrMissingHmac      = status.Errorf(codes.InvalidArgument, "missing x-hmac-signature metadata")
	ErrMissingHmacKeyID = status.Errorf(codes.InvalidArgument, "missing x-hmac-key-id metadata")
	ErrUnauthorized     = status.Errorf(codes.Unauthenticated, "unauthorized")
)

type clientAuth struct {
	hmacKeyId, hmacSecret string
}

type serverAuth struct {
	secrets map[string]string
}

// NewClientAuthInterceptor returns a new ClientAuthInterceptor to be used as a grpc.DialOption
// user StreamInterceptOpt or UnaryInterceptOpt to get the grpc.ServerOption.
func NewClientAuthInterceptor(hmacKeyId, hmacSecret string) grpc.UnaryClientInterceptor {
	c := &clientAuth{hmacKeyId, hmacSecret}
	return c.Interceptor
}

// NewServerAuthInterceptor returns a new ServerAuthInterceptor to be used as a grpc.ServerOption
// user StreamInterceptOpt or UnaryInterceptOpt to get the grpc.ServerOption.
func NewServerAuthInterceptor(secrets map[string]string) grpc.UnaryServerInterceptor {
	s := &serverAuth{secrets}
	return s.Interceptor
}

func (c *clientAuth) Interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	tflog.Info(ctx, "authenticating request")
	ctx = metadata.AppendToOutgoingContext(ctx, "x-hmac-key-id", c.hmacKeyId)
	plaintext, err := plainText(req, method)
	if err != nil {
		return err
	}
	log.Debug().Str("plaintext", plaintext).Msg("client plaintext")
	ctx = metadata.AppendToOutgoingContext(ctx, "x-hmac-signature", signature(c.hmacSecret, plaintext))
	tflog.Info(ctx, "sending request")
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (s *serverAuth) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Debug().Msg("authenticating request")
	var md metadata.MD
	var ok bool
	var hmacSign []string
	var hmacKeyID []string
	var secretKey string
	var err error
	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		log.Debug().Msg("missing metadata")
		return nil, ErrMissingMetadata
	}
	if hmacSign = md.Get("x-hmac-signature"); len(hmacSign) != 1 {
		log.Debug().Msg("missing x-hmac-signature metadata")
		return nil, ErrMissingHmac
	}
	if hmacKeyID = md.Get("x-hmac-key-id"); len(hmacKeyID) != 1 {
		log.Debug().Msg("missing x-hmac-key-id metadata")
		return nil, ErrMissingHmacKeyID
	}
	if secretKey, err = s.getHMACSecretKey(hmacKeyID[0]); err != nil {
		log.Debug().Err(err).Msg("failed to get HMAC secret key")
		return nil, err
	}
	plaintext, err := plainText(req, info.FullMethod)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get plaintext")
		return nil, err
	}
	log.Debug().Str("plaintext", plaintext).Msg("server plaintext")
	// Compare HMAC signatures.
	if !hmac.Equal([]byte(hmacSign[0]), signatureBytes(secretKey, plaintext)) {
		log.Debug().Msg("invalid HMAC signature")
		return nil, ErrUnauthorized
	}
	return handler(ctx, req)
}

func (s *serverAuth) getHMACSecretKey(key string) (string, error) {
	log.Debug().Str("key", key).Msg("getting HMAC secret key")
	if secretKey, ok := s.secrets[key]; ok {
		return secretKey, nil
	}
	return "", ErrUnauthorized
}

func plainText(req interface{}, method string) (string, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(req); err != nil {
		if strings.Contains(err.Error(), "has no exported fields") {
			goto ADDMETHOD
		}
		return "", fmt.Errorf("failed to encode request: %w", err)
	}
ADDMETHOD:
	buf.WriteString("method=" + method)
	return buf.String(), nil
}

// signatureBytes generates a HMAC signature and returns it as a base64 encoded []byte.
func signatureBytes(secretKey string, message string) []byte {
	return []byte(signature(secretKey, message))
}

// signature generates a HMAC signature and returns it as a base64 encoded string.
func signature(secretKey string, message string) string {
	mac := hmac.New(sha512.New512_256, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
