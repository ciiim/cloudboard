package smallfile_test

import (
	"os"
	"testing"

	"github.com/ciiim/cloudborad/smallfile"
	"github.com/go-playground/assert/v2"
)

// // 测试创建文件的性能
// func BenchmarkCreateFile(b *testing.B) {
// 	_ = os.Mkdir("tmp", 0755)

// 	// files := make([]*os.File, 0, b.N)

// 	// defer func() {
// 	// 	for _, file := range files {
// 	// 		file.Close()
// 	// 	}
// 	// 	os.RemoveAll("./tmp")
// 	// }()

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		file, err := os.CreateTemp("./tmp", "test*.txt")
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		// files = append(files, file)

// 		file.WriteString("data")

// 		file.Close()
// 		os.Remove(file.Name())
// 	}
// 	b.StopTimer()

// }

const (
	TestFileName = "db.dat"
)

func BenchmarkAllocCachePages1(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	defer a.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = a.Alloc(smallfile.PageSize * 1)
		id, err := a.Alloc(smallfile.PageSize * 1)
		if err != nil {
			b.Fatal(err)
		}
		span, err := a.Get(id)
		if err != nil {
			b.Fatal(err)
		}
		// spans = append(spans, id)
		copy(span.FixedBytes(), []byte("datadata"))

		a.ForceSync(id)
		// a.MarkDirty(id)

		a.Free(id)
	}
	b.StopTimer()
}

func BenchmarkAllocPages1(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	spans := make([]smallfile.GlobalPageID, 0, b.N)

	defer func(spans *[]smallfile.GlobalPageID) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i])
		}
	}(&spans)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.Alloc(smallfile.PageSize * 1)
		assert.Equal(b, err, nil)
		spans = append(spans, span)
	}
	b.StopTimer()

}

func BenchmarkAllocPages16(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	spans := make([]smallfile.GlobalPageID, 0, b.N)

	defer func(spans *[]smallfile.GlobalPageID) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i])
		}
	}(&spans)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.Alloc(smallfile.PageSize * 16)
		assert.Equal(b, err, nil)
		spans = append(spans, span)
	}
	b.StopTimer()
}

func BenchmarkAllocPages32(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	spans := make([]smallfile.GlobalPageID, 0, b.N)

	defer func(spans *[]smallfile.GlobalPageID) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i])
		}
	}(&spans)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.Alloc(smallfile.PageSize * 32)
		assert.Equal(b, err, nil)
		spans = append(spans, span)
	}
	b.StopTimer()
}

func BenchmarkAllocPages64(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	spans := make([]smallfile.GlobalPageID, 0, b.N)

	defer func(spans *[]smallfile.GlobalPageID) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i])
		}
	}(&spans)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.Alloc(smallfile.PageSize * 64)
		assert.Equal(b, err, nil)
		spans = append(spans, span)
	}
	b.StopTimer()
}

func BenchmarkAllocPages128(b *testing.B) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	spans := make([]smallfile.GlobalPageID, 0, b.N)

	defer func(spans *[]smallfile.GlobalPageID) {
		for i := 0; i < b.N; i++ {
			a.Free((*spans)[i])
		}
	}(&spans)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		span, err := a.Alloc(smallfile.PageSize * 128)
		assert.Equal(b, err, nil)
		spans = append(spans, span)
	}
	b.StopTimer()
}

func allocNPagesNoCache(t *testing.T, a *smallfile.Allocator, n, pages int, free bool) {
	spans := make([]*smallfile.SpanInCache, n)

	t.Logf("---test alloc and free %d page(s)", n)

	for i := 0; i < n; i++ {
		span, err := a.AllocNoCache(pages * smallfile.PageSize)
		assert.Equal(t, err, nil)
		if err != nil {
			a.Dump(0, 1024, 12)
		}
		spans[i] = span
	}
	if free {
		for i := 0; i < n; i++ {
			a.Free(spans[i].Id())
		}
	}
}

func allocNPagesWithCache(t *testing.T, a *smallfile.Allocator, n, pages int, free bool) {
	spans := make([]smallfile.GlobalPageID, n)

	t.Logf("---test alloc and free %d page(s)", n)

	for i := 0; i < n; i++ {
		id, err := a.Alloc(pages * smallfile.PageSize)
		assert.Equal(t, err, nil)

		spans[i] = id

		span, err := a.Get(id)
		assert.Equal(t, err, nil)

		copy(span.FixedBytes(), []byte("test"))

		a.ForceSync(id)
	}
	if free {
		for i := 0; i < n; i++ {
			a.Free(spans[i])
		}
	}
}

func TestAllocWithCacheAndFree(t *testing.T) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)

	allocNPagesNoCache(t, a, 16, 1, true)
	allocNPagesNoCache(t, a, 16, 4, true)
	allocNPagesNoCache(t, a, 16, 8, true)

	// a.Dump(0, 0, 12)

}

func TestAllocNoCacheAndFree(t *testing.T) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)

	allocNPagesWithCache(t, a, 16, 1, true)
	allocNPagesWithCache(t, a, 32, 1, true)

	allocNPagesWithCache(t, a, 16, 4, true)
	allocNPagesWithCache(t, a, 32, 4, true)

	allocNPagesWithCache(t, a, 16, 8, true)
	allocNPagesWithCache(t, a, 32, 8, true)

}

func TestNewChunk(t *testing.T) {
	a := smallfile.NewAllocator(TestFileName)
	defer os.Remove(TestFileName)
	defer a.Close()

	//分配chunk最大page数的span
	id, err := a.Alloc(smallfile.MaxUserPagePerChunk * smallfile.PageSize)
	assert.Equal(t, err, nil)

	pages, _ := a.GetSpanPages(id)
	t.Log(id, pages)

	// 第二个chunk
	id2, err := a.Alloc(1024)
	assert.Equal(t, err, nil)

	pages, _ = a.GetSpanPages(id2)
	t.Log(id2, pages)
}