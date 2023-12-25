package errcode

// DatabaseErrorx 扩展数据库错误的错误信息类
type DatabaseErrorx struct {
	*ErrCode
}

func NewDatabaseErrorx() *DatabaseErrorx {
	return &DatabaseErrorx{DatabaseError}
}

func (e *DatabaseErrorx) GetError(err error) *ErrCode {
	e.SetMsg("获取数据时发生错误").SetDetails(err.Error())
	return e.ErrCode
}
func (e *DatabaseErrorx) UpdateError(err error) *ErrCode {
	e.SetMsg("更新数据时发生错误").SetDetails(err.Error())
	return e.ErrCode
}

func (e *DatabaseErrorx) DeleteError(err error) *ErrCode {
	e.SetMsg("删除数据时发生错误").SetDetails(err.Error())
	return e.ErrCode
}

func (e *DatabaseErrorx) CreateError(err error) *ErrCode {
	e.SetMsg("创建数据时发生错误").SetDetails(err.Error())
	return e.ErrCode
}
