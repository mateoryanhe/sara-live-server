package entity

// Init CMS相关表迁移与 syndb 注册
func Init() {
	initCmsToken()
	InitCMSUser()
	InitCMSRole()
	InitPermission()
	initAccountCfg()
	initAppVersionCfg()
	initAppVersionUpdateDetail()
	initSimulatorCpuKeyword()
}
