package process

import (
	"github.com/google/uuid"
	"strings"
)

func GetNewUUID() string {
	res := uuid.New().String()
	res = strings.ReplaceAll(res, "-", "")

	return res
}
