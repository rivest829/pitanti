package utils

import "testing"

func TestGenRoleId(t *testing.T) {
	serverIds := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for _, serverId := range serverIds {
		for range 3 {
			roleId := GenRoleId(serverId)
			t.Logf("serverId: %d, roleId: %s", serverId, roleId)
			if GetServerId(roleId) != serverId {
				t.Errorf("serverId: %d, roleId: %s, parsed serverId: %d", serverId, roleId, GetServerId(roleId))
			}
		}
	}
}

func TestGetServerId(t *testing.T) {
	// 测试特定的 roleId 解析
	testCases := []struct {
		serverId uint32
	}{
		{1}, {2}, {3}, {10}, {100}, {1000}, {8191}, // 8191 是 0x1fff 的最大值
	}

	for _, tc := range testCases {
		roleId := GenRoleId(tc.serverId)
		parsedServerId := GetServerId(roleId)
		if parsedServerId != tc.serverId {
			t.Errorf("Expected serverId %d, got %d for roleId %s", tc.serverId, parsedServerId, roleId)
		}
	}
}
