package utils

import (
	"math/rand/v2"
	"strconv"
	"strings"
)

func GetServerId(roleId string) uint32 {
	return uint32(ParseBase36(roleId) & 0x1fff)
}

func GenRoleId(serverId uint32) string {
	randNum := uint64(rand.Uint32N(0x1fff))

	roleIdNum := randNum<<13 | (uint64(serverId) & 0x1fff) // 将 serverId 设置到低 13 位
	return ToBase36(roleIdNum)
}

func ToBase36(num uint64) string {
	return strings.ToLower(strconv.FormatUint(num, 36))
}

func ParseBase36(id string) int64 {
	num, err := strconv.ParseInt(strings.ToLower(id), 36, 64)
	Must(err)
	return num
}
