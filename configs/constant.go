package configs

import (
	"os"

	"github.com/aryuuu/cepex-server/utils/converter"
)

type constant struct {
	Capacity int32
}

func initConstant() *constant {
	var capacity, _ = converter.ToInt(os.Getenv("CAPACITY"))

	result := &constant{
		Capacity: int32(capacity),
	}

	return result
}
