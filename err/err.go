package err

import "errors"

var (
	ErrDirExists      = errors.New("dir already exists")
	ErrFileNotFound   = errors.New("file not found")
	ErrNothingToStore = errors.New("there is nothing to be stored")
	ErrIllegalSeed    = errors.New("the RandomSeed is not a positive number")
	ErrIllFormedIP    = errors.New("the format of ipAddr is wrong")
	ErrConf           = errors.New("configure failed")
	ErrJoinNetwork    = errors.New("join network failed")
)
