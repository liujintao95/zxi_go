package conf

const SECRET_KEY = "XHSOI*Y9dfs9cshd9"
const COOKIE_MAXAGE = 60 * 60 * 24 * 7
const REDIS_MAXAGE = 60 * 60 * 12 * 1
const UPLOAD_DIRNAME = "/本地上传"

var MySqlConf = map[string]string{
	"user":     "root",
	"pwd":      "123456",
	"type":     "tcp",
	"address":  "127.0.0.1",
	"port":     "3306",
	"database": "zxi_net_disk",
}
