package cachemulti

import (
	"fmt"
	"testing"

	coretesting "cosmossdk.io/core/testing"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/store/cachekv"
	"cosmossdk.io/store/dbadapter"
	"cosmossdk.io/store/types"
)

func TestStoreGetKVStore(t *testing.T) {
	require := require.New(t)

	s := Store{stores: map[types.StoreKey]types.CacheWrap{}}
	key := types.NewKVStoreKey("abc")
	errMsg := fmt.Sprintf("kv store with key %v has not been registered in stores", key)

	require.PanicsWithValue(errMsg,
		func() { s.GetStore(key) })

	require.PanicsWithValue(errMsg,
		func() { s.GetKVStore(key) })
}

func bz(s string) []byte  { return []byte(s) }
func keyFmt(i int) []byte { return bz(fmt.Sprintf("key%0.8d", i)) }
func valFmt(i int) []byte { return bz(fmt.Sprintf("value%0.8d", i)) }

func TestClonedCacheMultiStore(t *testing.T) {
	mem := dbadapter.Store{DB: coretesting.NewMemDB()}
	cacheKV := cachekv.NewStore(mem)
	key := types.NewKVStoreKey("test")
	stores := map[types.StoreKey]types.CacheWrapper{
		key: cacheKV,
	}

	store := NewFromKVStore(mem, stores, map[string]types.StoreKey{}, nil, nil)

	st := store.GetKVStore(key)

	// do some initial setup to the store
	require.Empty(t, st.Get(keyFmt(1)), "Expected `key1` to be empty")

	// // put something in mem and in cache
	mem.Set(keyFmt(1), valFmt(1))
	st.Set(keyFmt(1), valFmt(1))
	require.Equal(t, valFmt(1), st.Get(keyFmt(1)))

	store.Write()
	require.Equal(t, valFmt(1), st.Get(keyFmt(1)))

	cloned := store.Clone()
	require.Equal(t, valFmt(1), cloned.GetKVStore(key).Get(keyFmt(1)))

	st.Set(keyFmt(1), valFmt(2))
	require.Equal(t, valFmt(1), cloned.GetKVStore(key).Get(keyFmt(1)))

	store.Write()
	require.Equal(t, valFmt(2), st.Get(keyFmt(1)))
	require.Equal(t, valFmt(1), cloned.GetKVStore(key).Get(keyFmt(1)))

	cloned.GetKVStore(key).Set(keyFmt(1), valFmt(3))
	store.Write()
	require.Equal(t, valFmt(2), st.Get(keyFmt(1)))
	require.Equal(t, valFmt(3), cloned.GetKVStore(key).Get(keyFmt(1)))

}
