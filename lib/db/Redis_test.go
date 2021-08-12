package db
/*
import (
	"omp/lib/util"
	"testing"
)

// go files -v -run TestRedis_SetToken Redis_test.go Redis.go
func TestRedis_SetToken(t *testing.T) {
	redisConn := GetRedisPool()
	uuId := util.NewUUID()
	_, err := redisConn.Do("Set", uuId, uuId, "EX", 3600*60)
	_ = redisConn.Close()
	if err != nil {
		t.Logf(`插入数据失败, %s`, err)
		return
	}
	t.Log(`插入数据成功`, uuId)
}

// go files -v -run TestRedis_GetToken Redis_test.go Redis.go
func TestRedis_GetToken(t *testing.T) {
	redisConn := GetRedisPool()
	rec1, err := redisConn.Do("Get", "fc89b827-9a48-4eab-86a4-f74bbc63e928")
	_ = redisConn.Close()
	if err != nil {
		t.Logf(`查询数据失败, %s`, err)
		return
	}
	t.Log(`查询数据成功`, string(rec1.([]byte)))
}
*/