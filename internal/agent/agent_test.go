package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"
	"time"

	api "github.com/ShinyaIshitobi/distributed-go/api/v1"
	"github.com/ShinyaIshitobi/distributed-go/internal/config"
	"github.com/ShinyaIshitobi/distributed-go/internal/loadbalance"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func TestAgent(t *testing.T) {
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	var agents []*Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]

		dataDir, err := os.MkdirTemp("", "agent-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		agent, err := New(Config{
			NodeName:        fmt.Sprintf("%d", i),
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			DataDir:         dataDir,
			StartJoinAddrs:  startJoinAddrs,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			Bootstrap:       i == 0,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir))
		}
	}()
	time.Sleep(3 * time.Second)

	leaderClient := client(t, agents[0], peerTLSConfig)
	produceResponse, err := leaderClient.Produce(
		context.Background(),
		&api.ProduceRequest{
			Record: &api.Record{Value: []byte("foo")},
		},
	)

	require.NoError(t, err)

	// wait for replication
	time.Sleep(3 * time.Second)

	consumeResponse, err := leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)

	require.NoError(t, err)
	require.Equal(t, "foo", string(consumeResponse.Record.Value))

	followerClient := client(t, agents[1], peerTLSConfig)
	consumeResponse, err = followerClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "foo", string(consumeResponse.Record.Value))

	consumeResponse, err = leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset + 1,
		},
	)
	require.Nil(t, consumeResponse)
	require.Error(t, err)
	got := status.Code(err)
	want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	require.Equal(t, got, want)
}

func client(t *testing.T, agent *Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf(
		"%s///%s",
		loadbalance.Name,
		rpcAddr,
	), opts...)
	require.NoError(t, err)

	return api.NewLogClient(conn)
}
