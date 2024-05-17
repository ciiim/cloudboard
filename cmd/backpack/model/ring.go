package model

import (
	"crypto/md5"
	"log/slog"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"

	"github.com/ciiim/cloudborad/chunkpool"
	"github.com/ciiim/cloudborad/ringio"
	"github.com/urfave/cli/v2"
)

const (
	secretKey = "Ew3124*9djfercS@"
)

var Ring *RingModel

type RingModel struct {
	*ringio.Ring
	chunkPool *chunkpool.ChunkPool
}

func EncryptPassword(password string) string {
	return crypto.
		FromString(password).
		SetKey(secretKey).
		SetIv(secretKey).
		Aes().
		CBC().
		PKCS7Padding().
		Encrypt().
		ToBase64String()
}

func DecryptToken(token string) string {
	return crypto.
		FromBase64String(token).
		SetKey(secretKey).
		SetIv(secretKey).
		Aes().
		CBC().
		PKCS7Padding().
		Decrypt().
		ToString()
}

func md5Hash(data []byte) []byte {
	sum := md5.Sum(data)
	return sum[:]
}

func InitRingModel(flags *cli.Context) {
	config := &ringio.RingConfig{
		Name:          flags.String("hostname"),
		Port:          flags.Int("port"),
		Replica:       flags.Int("replica"),
		ChunkMaxSize:  ringio.DefaultChunkSize,
		HashFn:        md5Hash,
		RootPath:      flags.String("root"),
		LogLevel:      slog.LevelInfo,
		EnableReplica: true, //FIXME: test
	}

	Ring = NewRingModel(ringio.NewRing(config))

}

func NewRingModel(ring *ringio.Ring) *RingModel {
	return &RingModel{
		Ring:      ring,
		chunkPool: chunkpool.NewChunkPool(ring.StorageSystem.Config().HCSConfig.ChunkMaxSize),
	}
}
