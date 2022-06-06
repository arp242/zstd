package zsync_test

import (
	"reflect"
	"sync"
	"testing"

	"zgo.at/zstd/zsync"
)

func TestSynced(t *testing.T) {
	{
		s := zsync.NewSynced[string]("")
		if s.Get() != "" {
			t.Error()
		}

		s.Set("asd")
		if s.Get() != "asd" {
			t.Error()
		}
	}

	{
		s := zsync.NewSynced[[]string](nil)
		if s.Get() != nil {
			t.Error()
		}

		s.Set([]string{"A", "B"})
		if !reflect.DeepEqual(s.Get(), []string{"A", "B"}) {
			t.Error()
		}
	}

	{
		s := zsync.NewSynced[int](0)

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			func(i int) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = s.Get()
					s.Set(i)
				}()
			}(i)
		}
		wg.Wait()
	}
}
