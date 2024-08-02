package v1

import (
	"context"
	"google.golang.org/grpc"
	"plugin-demo/common"
	"plugin-demo/proto"
)

// BackupItemActionGRPCClient implements the backup/ItemAction interface and uses a
// gRPC client to make calls to the plugin server.
type PluginClient struct {
	Context context.Context
	*common.ClientBase
	grpcClient proto.KVClient
}

func newBackupItemActionGRPCClient(base *common.ClientBase, clientConn *grpc.ClientConn) interface{} {
	return &PluginClient{
		Context:    context.Background(),
		ClientBase: base,
		grpcClient: proto.NewKVClient(clientConn),
	}
}

func (c *PluginClient) Get(key string) ([]byte, error) {
	resp, err := c.grpcClient.Get(c.Context, &proto.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *PluginClient) Put(key string, value []byte) error {
	_, err := c.grpcClient.Put(c.Context, &proto.PutRequest{Key: key, Value: value})
	return err
}
