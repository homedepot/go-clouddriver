package clouddriver

type Permissions struct {
	READ  []string `json:"READ"`
	WRITE []string `json:"WRITE"`
}

type ReadPermission struct {
	ID          string `json:"-" gorm:"primary_key"`
	AccountName string `json:"accountName" gorm:"index:account_name_idx"`
	ReadGroup   string `json:"readGroup"`
}

func (ReadPermission) TableName() string {
	return "provider_read_permissions"
}

type WritePermission struct {
	ID          string `json:"-" gorm:"primary_key"`
	AccountName string `json:"accountName" gorm:"index:account_name_idx"`
	WriteGroup  string `json:"writeGroup"`
}

func (WritePermission) TableName() string {
	return "provider_write_permissions"
}
