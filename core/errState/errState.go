package errState

const (
	ErrBadReq     = 0001
	ErrDBConn     = 1001
	ErrDBUpdate   = 1002
	ErrDBDelete   = 1003
	ErrDBCreate   = 1004
	ErrDBSelect   = 1005
	ErrSaveFile   = 2001
	ErrCreatePath = 2002
)

var statusText = map[int]string{
	ErrBadReq     : "请求数据不规范",
	ErrDBConn     : "数据库连接出错",
	ErrDBUpdate   : "数据库更新数据出错",
	ErrDBDelete   : "数据库删除数据出错",
	ErrDBCreate   : "数据库创建数据出错",
	ErrDBSelect   : "数据库查询数据出错",
	ErrSaveFile   : "保存文件到本地磁盘失败",
	ErrCreatePath : "创建本地路径失败",
}

func StatusText(code int) string {
	return statusText[code]
}