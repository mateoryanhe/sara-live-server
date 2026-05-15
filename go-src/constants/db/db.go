package db

type TbName string
type TbCol string

const (
	// IdName 默认表主键名为id
	IdName        TbCol = "id"
	CreatedAtName TbCol = "created_at"
	UpdatedAtName TbCol = "updated_at"
	DeletedAtName TbCol = "deleted_at"
	IsDeletedName TbCol = "is_deleted"
)
