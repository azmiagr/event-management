package entity

type Role struct {
	RoleID   int    `json:"role_id" gorm:"type:int;primaryKey;autoIncrement"`
	RoleName string `json:"role_name" gorm:"type:varchar(25)"`

	Users []User `json:"users" gorm:"foreignKey:RoleID"`
}
