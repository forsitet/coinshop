package service 

import "errors"

var ErrItemNotFound = errors.New("item not found")
var ErrCoinNotEnough = errors.New("not enough coins")
