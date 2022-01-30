package btree

import (
	"crypto/sha1"
	"os"
	"testing"

	assert2 "github.com/stretchr/testify/assert"
)

func tempfile(t *testing.T) string {
	temp, err := os.CreateTemp(".", "temp.*.db")
	assert2.NoError(t, temp.Close())
	assert2.NoError(t, err)
	return temp.Name()
}

func TestCreate(t *testing.T) {
	f := tempfile(t)
	_, err := Create(f)
	assert2.NoError(t, err)
	defer os.Remove(f)
}

func TestBtree(t *testing.T) {
	f := tempfile(t)
	defer os.Remove(f)
	bt, err := Create(f)
	assert2.NoError(t, err)

	tests := []string{
		"hello world",
		"123",
		"We are met on a great battle-field of that war.",
		"Abraham Lincoln, November 19, 1863, Gettysburg, Pennsylvania",
	}
	sha := make([]*byte, len(tests))
	for i, tt := range tests {
		hash := sha1.New()
		hash.Write([]byte(tt))
		sha[i] = &hash.Sum(nil)[0]
		bt.Insert(sha[i], []byte(tt))
	}
	assert2.NoError(t, bt.Close())

	bt, err = Open(f)
	assert2.NoError(t, err)
	var ss []string
	bt.Foreach(func(key [16]byte, value []byte) {
		ss = append(ss, string(value))
	})
	assert2.ElementsMatch(t, tests, ss)

	for i, tt := range tests {
		assert2.Equal(t, []byte(tt), bt.Get(sha[i]))
	}

	for i := range tests {
		assert2.NoError(t, bt.Delete(sha[i]))
	}

	for i := range tests {
		assert2.Equal(t, []byte(nil), bt.Get(sha[i]))
	}

	assert2.NoError(t, bt.Close())
}
