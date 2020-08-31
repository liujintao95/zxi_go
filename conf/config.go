package conf

const SECRET_KEY = "XHSOI*Y9dfs9cshd9"
const BLOCK_SIZE = 1024*1024*5
const SAVE_PATH = "./upload_file"

var MySqlConf = map[string]string{
	"user":     "root",
	"pwd":      "123456",
	"type":     "tcp",
	"address":  "127.0.0.1",
	"port":     "3306",
	"database": "zxi_net_disk",
}
