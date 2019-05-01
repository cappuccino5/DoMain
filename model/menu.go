package model

/***************这是后台权限管理管理数据结构*******
1.Admin.RoleName关联EmployeeRole.RoleName
2.EmployeeRole.Permission关联Menu.Id
3.Menu表中填写具体的域名，主菜单，子菜单
4.Menu.SubMenu为子菜单，为拼接字符串，有多个url
*********************************************/
// 菜单，包含域名，1级菜单，1级下子菜单
type Menu struct {
	Id       uint64 `json:"id"`
	DoMain   string `json:"do_main"`
	MainMenu string `json:"main_menu"`
	SubMenu  string `json:"sub_menu"`
}

//职务角色增删查改
type EmployeeRole struct {
	RoleName   string `json:"role_name"` // 职位名是唯一
	Permission string `json:"permission"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time,omitempty"`
	UpdateTime uint64 `json:"update_time"`
}

// 管理员账号
type Admin struct {
	Id         uint64 `json:"id"`
	UserName   string `json:"user_name"`
	RoleName   string `json:"role_name"` // 职位名
	Status     int    `json:"status"`
	CreateTime uint64 `json:"create_time"`
	UpdateTime uint64 `json:"update_time"`
	Password   string `json:"password"` // 密码加盐
	Salt       string `json:"salt"`     // 盐
}

