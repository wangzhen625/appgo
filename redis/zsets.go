package redis

import (
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/oxfeeefeee/appgo"
)

type Zsets struct {
	namespace string
}

func NewZsets(namespace string) *Zsets {
	return &Zsets{namespace}
}

func (z *Zsets) Add(key interface{}, score float64, item interface{}) error {
	if _, err := Do("ZADD", z.keystr(key), score, item); err != nil {
		return err
	}
	return nil
}

func (z *Zsets) Rem(key, item interface{}) error {
	if _, err := Do("ZREM", z.keystr(key), item); err != nil {
		return err
	}
	return nil
}

func (z *Zsets) Min(key interface{}) (interface{}, error) {
	return oneFromSlice(z.Range(key, 0, 0))
}

func (z *Zsets) Max(key interface{}) (interface{}, error) {
	return oneFromSlice(z.Range(key, -1, -1))
}

func (z *Zsets) Range(key interface{}, b, e int) ([]interface{}, error) {
	return redigo.Values(Do("ZRANGE", z.keystr(key), b, e))
}

func (z *Zsets) MinStr(key interface{}) (string, error) {
	return redigo.String(z.Min(key))
}

func (z *Zsets) MaxStr(key interface{}) (string, error) {
	return redigo.String(z.Max(key))
}

func (z *Zsets) RangeStr(key interface{}, b, e int) ([]string, error) {
	return redigo.Strings(z.Range(key, b, e))
}

func (z *Zsets) keystr(k interface{}) string {
	return fmt.Sprintf("%s:%v", z.namespace, k)
}

func oneFromSlice(vals []interface{}, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}
	if len(vals) == 0 {
		return nil, appgo.NotFoundErr
	}
	return vals[0], nil
}
