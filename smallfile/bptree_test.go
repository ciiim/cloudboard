package smallfile_test

import (
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/ciiim/cloudborad/smallfile"
)

func genKey(key string) []byte {
	sum := sha512.Sum512([]byte(key))
	return sum[:]
}

func TestInsert(t *testing.T) {
	tree := smallfile.NewBPTree(4)
	for i := 0; i < 100; i++ {
		err := tree.Insert(genKey(fmt.Sprintf("%d.txt", i)), 1)
		if err != nil {
			t.Error(err)
			return
		}
	}

}

func TestSearch(t *testing.T) {
	tree := smallfile.NewBPTree(4)
	value, found := tree.Search(genKey("1.txt"))
	if !found {
		t.Error("not found")
		return
	}
	if value != 1 {
		t.Error("value not match")
	}
	t.Log("value:", value)
}

func TestDump(t *testing.T) {
	tree := smallfile.NewBPTree(4)
	tree.Dump(1)
}
