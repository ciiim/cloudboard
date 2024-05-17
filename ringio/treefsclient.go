package ringio

import (
	"context"
	"errors"

	"github.com/ciiim/cloudborad/node"
	"github.com/ciiim/cloudborad/ringio/fspb"
	"github.com/ciiim/cloudborad/storage/tree"
	"github.com/ciiim/cloudborad/storage/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type rpcTreeClient struct {
}

func newRPCTreeClient() *rpcTreeClient {
	return &rpcTreeClient{}
}

func (r *rpcTreeClient) getSpaceStat(ctx context.Context, ni *node.Node, space string, key string) (*tree.SpaceStatElement, error) {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.GetSpaceStat(ctx, &fspb.GetSpaceStatRequest{
		Space: &fspb.SpaceRequest{Space: space},
		Key:   key,
	})
	if err != nil {
		return nil, err
	}
	return PbSpaceStatToSpaceStat(resp.GetStat()), nil
}

func (r *rpcTreeClient) setSpaceStat(ctx context.Context, ni *node.Node, space string, e *tree.SpaceStatElement) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	_, err = client.SetSpaceStat(ctx, &fspb.SetSpaceStatRequest{
		Space: &fspb.SpaceRequest{Space: space},
		Stat: &fspb.SpaceStat{
			Key:   e.Key(),
			Value: e.Value(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

var empty = &emptypb.Empty{}

func (r *rpcTreeClient) allSpaces(ctx context.Context, ni *node.Node) ([]tree.SpaceInfo, error) {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.AllSpaces(ctx, empty)
	if err != nil {
		return nil, err
	}
	return PbSpacesToSpaces(resp), nil
}

func (r *rpcTreeClient) getMetadata(ctx context.Context, ni *node.Node, space string, base string, name string) ([]byte, error) {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &fspb.TreeFileSystemBasicRequest{
		Space: space,
		Base:  base,
		Name:  name,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (r *rpcTreeClient) putMetadata(ctx context.Context, ni *node.Node, space string, base string, name string, data []byte) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.PutMetadata(ctx, &fspb.PutMetadataRequest{
		Src: &fspb.TreeFileSystemBasicRequest{
			Space: space,
			Base:  base,
			Name:  name,
		},
		Metadata: data,
	})
	respErr := errors.New(resp.Err)
	if respErr != nil {
		return err
	}

	return respErr
}

func (r *rpcTreeClient) deleteMetadata(ctx context.Context, ni *node.Node, space string, base string, name string) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.DeleteMetadata(ctx, &fspb.TreeFileSystemBasicRequest{
		Space: space,
		Base:  base,
		Name:  name,
	})
	if err != nil {
		return err
	}

	if resp.GetErr() != "" {
		return errors.New(resp.GetErr())
	}
	return nil
}

func (r *rpcTreeClient) makeDir(ctx context.Context, ni *node.Node, space string, base string, dir string) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	_, err = client.MakeDir(ctx, &fspb.TreeFileSystemBasicRequest{
		Space: space,
		Base:  base,
		Name:  dir,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *rpcTreeClient) renameDir(ctx context.Context, ni *node.Node, space string, base string, dir string, newName string) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	_, err = client.RenameDir(ctx, &fspb.RenameDirRequest{
		Src: &fspb.TreeFileSystemBasicRequest{
			Space: space,
			Base:  base,
			Name:  dir,
		},
		NewName: newName,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *rpcTreeClient) deleteDir(ctx context.Context, ni *node.Node, space string, base string, dir string) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	_, err = client.DeleteDir(ctx, &fspb.TreeFileSystemBasicRequest{
		Space: space,
		Base:  base,
		Name:  dir,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *rpcTreeClient) getDirSub(ctx context.Context, ni *node.Node, space string, base string, dir string) ([]*tree.SubInfo, error) {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.GetDirSub(ctx, &fspb.TreeFileSystemBasicRequest{
		Space: space,
		Base:  base,
		Name:  dir,
	})
	if err != nil {
		return nil, err
	}
	return PbSubsToSubs(resp.SubInfo), nil
}

func (r *rpcTreeClient) newSpace(ctx context.Context, ni *node.Node, space string, cap types.Byte) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	resp, err := client.NewSpace(ctx, &fspb.NewSpaceRequest{
		Space: space,
		Cap:   int64(cap),
	})
	if err != nil {
		return err
	}
	if resp.GetErr() != "" {
		return errors.New(resp.GetErr())
	}
	return nil
}

func (r *rpcTreeClient) deleteSpace(ctx context.Context, ni *node.Node, space string) error {
	conn, err := grpc.Dial(ni.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := fspb.NewTreeFileSystemServiceClient(conn)
	_, err = client.DeleteSpace(ctx, &fspb.SpaceRequest{
		Space: space,
	})
	if err != nil {
		return err
	}
	return nil
}
