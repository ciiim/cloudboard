package ringio

import (
	"context"

	"github.com/ciiim/cloudborad/ringio/fspb"
	"github.com/ciiim/cloudborad/storage/types"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *rpcServer) GetSapceStat(ctx context.Context, req *fspb.GetSpaceStatRequest) (*fspb.GetSpaceStatResponse, error) {
	stat, err := r.tfs.Local().GetLocalSpace(req.GetSpace().GetSpace()).GetStatElement(req.GetKey())
	if err != nil {
		return nil, nil
	}
	return &fspb.GetSpaceStatResponse{
		Stat: &fspb.SpaceStat{
			Key:   stat.Key(),
			Value: stat.Value(),
		},
	}, nil
}

func (r *rpcServer) SetSapceStat(ctx context.Context, req *fspb.SetSpaceStatRequest) (*fspb.Error, error) {
	err := r.tfs.Local().GetLocalSpace(req.GetSpace().GetSpace()).SetStatElement(PbSpaceStatToSpaceStat(req.GetStat()))
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) AllSpaces(ctx context.Context, req *emptypb.Empty) (*fspb.SpaceInfos, error) {
	spaces := r.tfs.Local().AllSpaces()
	return SpacesToPbSpaces(spaces), nil
}

func (r *rpcServer) MakeDir(ctx context.Context, req *fspb.TreeFileSystemBasicRequest) (*fspb.Error, error) {
	err := r.tfs.MakeDir(req.Space, req.Base, req.Name)
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) RenameDir(ctx context.Context, req *fspb.RenameDirRequest) (*fspb.Error, error) {
	err := r.tfs.RenameDir(req.Src.Space, req.Src.Base, req.Src.Name, req.NewName)
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) DeleteDir(ctx context.Context, req *fspb.TreeFileSystemBasicRequest) (*fspb.Error, error) {
	err := r.tfs.DeleteDir(req.Space, req.Base, req.Name)
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) GetDirSub(ctx context.Context, req *fspb.TreeFileSystemBasicRequest) (*fspb.Subs, error) {
	subs, err := r.tfs.GetDirSub(req.Space, req.Base, req.Name)
	return &fspb.Subs{SubInfo: SubsToPbSubs(subs)}, err
}

func (r *rpcServer) NewSpace(ctx context.Context, space *fspb.NewSpaceRequest) (*fspb.Error, error) {
	err := r.tfs.Local().NewLocalSpace(space.Space, types.Byte(space.Cap))
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) DeleteSpace(ctx context.Context, space *fspb.SpaceRequest) (*fspb.Error, error) {
	err := r.tfs.Local().DeleteLocalSpace(space.Space)
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) GetMetadata(ctx context.Context, req *fspb.TreeFileSystemBasicRequest) (*fspb.BytesData, error) {
	data, err := r.tfs.GetMetadata(req.Space, req.Base, req.Name)
	if err != nil {
		return nil, err
	}
	return &fspb.BytesData{Data: data}, nil
}

func (r *rpcServer) PutMetadata(ctx context.Context, req *fspb.PutMetadataRequest) (*fspb.Error, error) {
	err := r.tfs.PutMetadata(req.Src.Space, req.Src.Base, req.Src.Name, req.Src.Hash, req.Metadata)
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}

func (r *rpcServer) DeleteMetadata(ctx context.Context, req *fspb.TreeFileSystemBasicRequest) (*fspb.Error, error) {
	err := r.tfs.DeleteMetadata(req.GetSpace(), req.GetBase(), req.GetName())
	if err != nil {
		return &fspb.Error{Err: err.Error()}, nil
	}
	return &fspb.Error{}, nil
}
