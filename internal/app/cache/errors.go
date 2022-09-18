package cache

import "errors"

var (
	ErrOldStats = errors.New("cannot add old stats to cache")
)
