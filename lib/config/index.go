package config

// 端口
const AppPort int64 = 8881

// 路由配置
const BaseUrl string = "/file"
const ApiBaseUrl string = "/file-api"
const BaseAssets string = "/file-assets"

// mysql 配置
const DB_W_host string = ""
const DB_W_port int64 = 3306
const DB_W_username string = ""
const DB_W_password string = ""
const DB_W_database string = ""

// redis 配置
const RedisHost string = ""
const RedisPort int64 = 6379
const RedisPassword string = ""

// OSS 配置
const OSS_protocol string = "https://"
const OSS_domain string = ""
const OSS_endpoint string = OSS_protocol+ OSS_domain
const OSS_accessKeyID string  = ""
const OSS_accessKeySecret string  = ""