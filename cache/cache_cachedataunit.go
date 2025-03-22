package redis_cache

import (
	"golang.org/x/exp/constraints"
)

type CustomDataType interface {
	constraints.Ordered | []byte | []rune
}

type CacheDataUnit struct {
	Key                 string
	Data                CustomDataType
	LastUpdateTimestamp time.time
}
