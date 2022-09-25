package err

import "errors"

var (
	ErrDirExists error = errors.New("dir already exists")
	ErrFileNotFound error = errors.New("file not found")
	ErrNothingToStore error = errors.New("there is nothing to be stored")
	ErrIllegalSeed error = errors.New("the RandomSeed is not a positive number")
	ErrIllFormedIP error = errors.New("the format of ipAddr is wrong")
)