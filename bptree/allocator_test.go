package bptree_test

import (
	"os"
	"testing"

	"github.com/ciiim/cloudborad/bptree"
	"github.com/go-playground/assert/v2"
)

func BenchmarkAllocNPages(b *testing.B) {
	a := bptree.NewAllocator()
	defer os.Remove("db.dat")
	spans := make([]*bptree.SpanInCache, 0, b.N)

	defer func(spans *[]*bptree.SpanInCache) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i].Id())
		}
	}(&spans)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.AllocNoCache(bptree.PageSize)
		if err != nil {
			b.Fatal(err)
		}
		spans = append(spans, span)
	}
	b.StopTimer()
}

func allocNPages(t *testing.T, a *bptree.Allocator, n int, free bool) {
	spans := make([]*bptree.SpanInCache, n)

	t.Logf("---test alloc and free %d page(s)", n)

	for i := 0; i < n; i++ {
		span, err := a.AllocNoCache(n * bptree.PageSize)
		if err != nil {
			assert.Equal(t, err, nil)
		}
		spans[i] = span
	}
	if free {
		for i := 0; i < n; i++ {
			a.Free(spans[i].Id())
		}
	}
}

func TestAllocNoCacheAndFree(t *testing.T) {
	a := bptree.NewAllocator()
	defer os.Remove("db.dat")

	allocNPages(t, a, 1, true)
	allocNPages(t, a, 4, true)
	allocNPages(t, a, 8, true)

}
