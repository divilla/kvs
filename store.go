package kvs

import (
	"bytes"
	"errors"
	"hash/fnv"
	"math"
	"math/rand"
	"strings"
	"time"

	"blainsmith.com/go/seahash"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"

	idxSize = uint64(math.MaxInt32 / math.MaxInt8)
)

var (
	rndKeySize = 1024
	src        = rand.NewSource(time.Now().UnixNano())
	hsh        = fnv.New64a()
)

type Store struct {
	kvp  []kvp
	idx  [][]int
	del  []int
	next int // next kvp index
	sb   strings.Builder
}

type kvp struct {
	key []byte
	val interface{}
}

func NewStore(kvpSize int) *Store {
	return &Store{
		kvp: make([]kvp, kvpSize),
		idx: make([][]int, idxSize),
	}
}

func (s *Store) Get(key string) (interface{}, error) {
	b := []byte(key)
	h, err := hash(b)
	if err != nil {
		return nil, err
	}
	return nil, nil

	for i := 0; i < len(s.idx[h]); i++ {
		idx := s.idx[h][i]
		if bytes.Equal(s.kvp[idx].key, b) {
			b = nil
			return s.kvp[i].val, nil
		}
	}

	return nil, errors.New("key not found")
}

func (s *Store) Set(key string, val interface{}) error {
	b := []byte(key)
	h, err := hash(b)
	if err != nil {
		return err
	}

	// check if value exists
	for i := 0; i < len(s.idx[h]); i++ {
		idx := s.idx[h][i]
		if bytes.Equal(s.kvp[idx].key, b) {
			s.kvp[idx].val = val
			return nil
		}
	}

	// reuse deleted kvp if possible
	if len(s.del) > 0 {
		i := s.del[len(s.del)-1]
		s.del = s.del[:len(s.del)-1]
		s.kvp[i].key = b
		s.kvp[i].val = val
		return nil
	}

	// use pre-allocated space if possible
	if s.next < len(s.kvp) {
		s.kvp[s.next].key = b
		s.kvp[s.next].val = val
	} else {
		s.kvp = append(s.kvp, kvp{key: b, val: val})
	}

	// add index
	s.idx[h] = append(s.idx[h], s.next)
	s.next++

	return nil
}

func (s *Store) SetRndKey(val interface{}) (string, error) {
	size := rand.Intn(rndKeySize)
	s.sb.Reset()
	s.sb.Grow(size)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := size-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			s.sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	str := s.sb.String()
	err := s.Set(str, val)
	return str, err
}

func (s *Store) Del(key string) error {
	b := []byte(key)
	h, err := hash(b)
	if err != nil {
		return err
	}

	for k, i := range s.idx[h] {
		if bytes.Equal(s.kvp[i].key, b) {
			// just remove index
			s.idx[h][k] = s.idx[h][len(s.idx[h])-1]
			s.idx[h] = s.idx[h][:len(s.idx[h])-1]
			// mark index as deleted
			s.del = append(s.del, i)
			return nil
		}
	}

	return errors.New("key not found")
}

func (s *Store) Len() int {
	return s.next
}

func hash(b []byte) (int, error) {
	return int(seahash.Sum64(b) % idxSize), nil
}
