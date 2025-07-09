package suite

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	ssov1 "github.com/m1al04949/contracts/contracts/gen/go/sso"
	"github.com/m1al04949/sso-gRPC/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("error loading .env file")
	}
	path := os.Getenv("CONFIG_TESTS_PATH")

	cfg := config.MustLoadByPath(path)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.DialContext(context.Background(), grpcAdress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAdress(cfg *config.Config) string {

	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
