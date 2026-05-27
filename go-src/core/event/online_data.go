package event

type OnlineData struct {
	RoleId uint64
}

func NewOnlineData(roleId uint64) *OnlineData {
	return &OnlineData{
		RoleId: roleId,
	}
}

type OfflineData struct {
	RoleId uint64
}

func NewOfflineData(roleId uint64) *OfflineData {
	return &OfflineData{
		RoleId: roleId,
	}
}
