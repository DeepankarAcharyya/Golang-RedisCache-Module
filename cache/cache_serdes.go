package redis_cache

import "github.com/vmihailenco/msgpack"

func (c *CacheDataUnit) Serialize(data CacheDataUnit) ([]byte, error) {
	return msgpack.Marshal(data)
}

func Desiarlize(data []byte, output_data interface{}) error {
	return msgpack.Unmarshal(data, output_data)
}
