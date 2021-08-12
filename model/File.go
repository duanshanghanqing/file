package model

type File struct {
	ID         int64  `gorm:"type:int;primaryKey;auto_increment"`
	FileId     string `gorm:"column:file_id;type:varchar(100);not null;default:''"`
	Bucket	   string `gorm:"column:bucket;type:varchar(100);not null;default:''"`
	Prefix	   string `gorm:"column:prefix;type:varchar(100);not null;default:''"`
	Name       string `gorm:"column:name;type:varchar(100);not null;default:''"`
	OriginName string `gorm:"column:origin_name;type:varchar(100);not null;default:''"`
	Suffix     string `gorm:"column:suffix;type:varchar(100);not null;default:''"`
	Size       int64  `gorm:"column:size;type:int;not null;default:0"`
	State      string `gorm:"column:state;type:varchar(10);not null;default:''"`
	CreateTime int64  `gorm:"column:create_time;type:bigint;not null;default:0"`
	UpdateTime int64  `gorm:"column:update_time;type:bigint;not null;default:0"`
}

// 文件删除状态
const (
	FileNormal         EnumType = iota
	FileLogicalDelete           // 逻辑删除
	FilePhysicalDelete          // 物理删除
)

var FileState = map[EnumType]string{
	FileNormal:         "正常",
	FileLogicalDelete:  "逻辑删除",
	FilePhysicalDelete: "物理删除",
}
