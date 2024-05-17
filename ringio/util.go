package ringio

import (
	"github.com/ciiim/cloudborad/replica"
	"github.com/ciiim/cloudborad/ringio/fspb"
	"github.com/ciiim/cloudborad/storage/hashchunk"
	"github.com/ciiim/cloudborad/storage/tree"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func HashChunkInfoToPBChunkInfo(info *hashchunk.HashChunkInfo) *fspb.HashChunkInfo {
	if info == nil {
		return nil
	}
	return &fspb.HashChunkInfo{
		ChunkCount: info.ChunkCount,
		ChunkName:  info.ChunkName,
		ChunkHash:  info.ChunkHash,
		BasePath:   info.ChunkPath,
		Size:       info.ChunkSize,
		ModTime:    timestamppb.New(info.ChunkModTime),
		CreateTime: timestamppb.New(info.ChunkCreateTime),
	}
}

func PBChunkInfoToHashChunkInfo(pb *fspb.HashChunkInfo) *hashchunk.HashChunkInfo {
	if pb == nil {
		return nil
	}
	return &hashchunk.HashChunkInfo{
		ChunkCount:   pb.ChunkCount,
		ChunkName:    pb.ChunkName,
		ChunkHash:    pb.ChunkHash,
		ChunkPath:    pb.BasePath,
		ChunkSize:    pb.Size,
		ChunkModTime: pb.ModTime.AsTime(),
	}
}

func PbSpaceStatToSpaceStat(pb *fspb.SpaceStat) *tree.SpaceStatElement {
	if pb == nil {
		return nil
	}
	return tree.NewSpaceStatElement(pb.Key, pb.Value)
}

func SpacesToPbSpaces(spaces []tree.SpaceInfo) *fspb.SpaceInfos {
	pbSpaces := new(fspb.SpaceInfos)
	pbSpaces.SpaceName = make([]string, len(spaces))
	for i, v := range spaces {
		pbSpaces.SpaceName[i] = v.SpaceName
	}
	return pbSpaces
}

func PbSpacesToSpaces(pb *fspb.SpaceInfos) []tree.SpaceInfo {
	spaces := make([]tree.SpaceInfo, len(pb.SpaceName))
	for i, v := range pb.SpaceName {
		spaces[i] = tree.SpaceInfo{SpaceName: v}
	}
	return spaces
}

func PbSubsToSubs(pb []*fspb.SubInfo) []*tree.SubInfo {
	subs := make([]*tree.SubInfo, len(pb))
	for i, v := range pb {
		subs[i] = &tree.SubInfo{
			Name:    v.Name,
			IsDir:   v.IsDir,
			ModTime: v.ModTime.AsTime(),
		}
	}
	return subs
}

func SubsToPbSubs(subs []*tree.SubInfo) []*fspb.SubInfo {
	pbSubs := make([]*fspb.SubInfo, len(subs))
	for i, v := range subs {
		pbSubs[i] = &fspb.SubInfo{
			Name:    v.Name,
			IsDir:   v.IsDir,
			ModTime: timestamppb.New(v.ModTime),
		}
	}
	return pbSubs
}

func WithPort(addr string, port string) string {
	return addr + ":" + port
}

func ReplicaInfoToPBReplicaInfo(info *replica.ReplicaObjectInfoG[*hashchunk.HashChunkInfo]) *fspb.ReplicaChunkInfo {
	if info == nil {
		return nil
	}
	rci := &fspb.ReplicaChunkInfo{
		Key:          info.Key,
		ChunkInfo:    HashChunkInfoToPBChunkInfo(info.Custom),
		ReplicaCount: int64(info.ExpectedReplicaCount),
		Checksum:     info.Checksum,
		NodeIds:      info.All,
	}
	return rci
}

func PBReplicaInfoToReplicaInfo(pb *fspb.ReplicaChunkInfo) *replica.ReplicaObjectInfoG[*hashchunk.HashChunkInfo] {
	if pb == nil {
		return nil
	}
	chunkInfo := PBChunkInfoToHashChunkInfo(pb.ChunkInfo)
	return &replica.ReplicaObjectInfoG[*hashchunk.HashChunkInfo]{
		Key:                  pb.Key,
		ExpectedReplicaCount: int(pb.ReplicaCount),
		Checksum:             pb.Checksum,
		All:                  pb.NodeIds,
		Custom:               chunkInfo,
	}
}
