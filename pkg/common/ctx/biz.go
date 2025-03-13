package ctx

import "gorm.io/gorm"

func (c *Context) GetAppGorm() *gorm.DB {
	value, exists := c.Get(GormAppKey)
	if exists {
		return value.(*gorm.DB)
	}
	return nil
}
func (c *Context) GetAppDoris() *gorm.DB {
	value, exists := c.Get(DorisAppKey)
	if exists {
		return value.(*gorm.DB)
	}
	return nil
}
func (c *Context) GetContextGorm() *gorm.DB {
	value, exists := c.Get(GormContextKey)
	if exists {
		return value.(*gorm.DB)
	}
	return nil
}

func (c *Context) SetAppGorm(db *gorm.DB) {
	if db == nil {
		return
	}
	c.Set(GormAppKey, db)
}
func (c *Context) SetAppDoris(db *gorm.DB) {
	if db == nil {
		return
	}
	c.Set(DorisAppKey, db)
}
func (c *Context) SetContextGorm(db *gorm.DB) {
	if db == nil {
		return
	}
	c.Set(GormContextKey, db)
}

func (c *Context) GetContextUserId() string {
	return c.GetString(UserIdContextKey)
}
func (c *Context) GetContextUserKey() string {
	return c.GetString(UserKeyContextKey)
}
func (c *Context) GetContextUserName() string {
	return c.GetString(UserNameContextKey)
}
func (c *Context) GetContextAppId() string {
	return c.GetString(AppIdContextKey)
}

func (c *Context) GetContextCurrentDeptId() string {
	return c.GetString(CurrentDeptIdContextKey)
}
func (c *Context) GetContextToken() string {
	return c.GetString(TokenContextKey)
}
func (c *Context) GetContextIsManager() bool {
	return c.GetBool(ManagerContextKey)
}

func (c *Context) SetContextUserId(id string) {
	c.Set(UserIdContextKey, id)
}
func (c *Context) SetContextUserKey(id string) {
	c.Set(UserKeyContextKey, id)
}

func (c *Context) SetContextUserName(name string) {
	c.Set(UserNameContextKey, name)
}
func (c *Context) SetContextAppId(id string) {
	c.Set(AppIdContextKey, id)
}
func (c *Context) SetContextCurrentDeptId(id string) {
	c.Set(CurrentDeptIdContextKey, id)
}
func (c *Context) SetContextToken(id string) {
	c.Set(TokenContextKey, id)
}
func (c *Context) SetContextIsManager(isManager bool) {
	c.Set(ManagerContextKey, isManager)
}
