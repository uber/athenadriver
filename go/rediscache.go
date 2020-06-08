package athenadriver

import (
	"errors"
	"fmt"
	redis "github.com/go-redis/redis/v7"
	"strconv"
	"strings"
	"time"
)

// Redis Cache for limit rater and result set cache
type CacheInRedis struct {
	client *redis.Client
}

var gRedis *CacheInRedis = newCacheInRedis()

func GetCacheInRedis() *CacheInRedis {
	return gRedis
}

func newCacheInRedis() *CacheInRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     RedisServerHost+":"+strconv.Itoa(RedisServerPort),
		Password: RedisServerPassword, // no password set
		DB:       RedisServerDefaultDB,  // use default DB
	})
	return &CacheInRedis{
		client: client,
	}
}

func (c *CacheInRedis) SetQID(query string, dataSize int64, QID string, t int64) {
	c.client.Set("⌘" + QID, query, 0)
	c.client.Set("I" + QID, dataSize, 0)
	c.client.Set("Q" + query, QID, 0)
	c.client.Set("T" + QID, t, 0)
}

func (c *CacheInRedis) SetSaving(username string, QID string) {
	dataSize := c.GetDataSize(QID)
	c.client.LPush("q" + username, QID + "\t" + time.Now().String())
	c.client.IncrBy("S" + username, dataSize)
	c.client.IncrByFloat("M" + username, getCost(dataSize))
	c.client.Incr("C" + username)
}

func (c *CacheInRedis) SetTableLastModified(table string, lastModified int64) {
	c.client.Set("⚚" + table, lastModified, 0)
}

// https://docs.google.com/document/d/1hf6IzerIIEY0Xd9e7tUuT7Sa0ROq0L4YODlbGBh7S9A/edit#bookmark=id.6oxz4x3t955h
func (c *CacheInRedis) SetTableS3Location(table string, s3loction string) {
	c.client.Set("♠" + s3loction, table, 0)
	c.client.Set("⚇" + table, s3loction, 0)
}

// https://docs.google.com/document/d/1hf6IzerIIEY0Xd9e7tUuT7Sa0ROq0L4YODlbGBh7S9A/edit#bookmark=id.6oxz4x3t955h
func (c *CacheInRedis) GetTableFromS3Location(s3loction string) string {
	if val, err := c.client.Get("♠" + s3loction).Result(); err == nil {
		return val
	}
	return ""
}

func (c *CacheInRedis) SetCost(username string, dataSize int64, QID string) {
	c.client.LPush("q" + username, QID + "\t" + time.Now().String())
	c.client.Set("s" + username, dataSize, 0)
	c.setMiss(username)
	c.setMoney(username, getCost(dataSize))
}

func (c *CacheInRedis) setMiss(username string) {
	c.client.Incr("c" + username)
}

func (c *CacheInRedis) setMoney(username string, cost float64) {
	c.client.IncrByFloat("m" + username, cost)
}

func (c *CacheInRedis) RemoveQID(QID string) error {
	query := c.GetQuery(QID)
	if query == "" {
		return errors.New("query id " + QID + " not found")
	}
	if err := c.RemoveQuery(query); err != nil {
		return err
	}
	err := c.client.Del("⌘" + QID).Err()
	if err != nil {
		return err
	}
	err = c.client.Del("I" + QID).Err()
	if err != nil {
		return err
	}
	err = c.client.Del("T" + QID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheInRedis) RemoveQuery(s string) error {
	err := c.client.Del("Q" + s).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheInRedis) RemoveQueryAndQID(query string) error {
	QID := c.GetQID(query)
	if err := c.RemoveQID(QID); err != nil {
		return err
	}
	err := c.client.Del("Q" + query).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheInRedis) GetValidCachedQID(query string, table string) string {
	account := "henry"
	cachedQID := c.GetQID(query)
	if cachedQID != "" {
		QIDSetTime, err := c.client.Get("T" + cachedQID).Result()
		if err != nil {
			return ""
		}
		tableLastModified, err := c.client.Get("⚚" + account +"."+ table).Result()
		if tableLastModified != "" {
			if QIDSetTime < tableLastModified {
				return ""
			}
		}
	}
	return cachedQID
}

func (c *CacheInRedis) GetQID(query string) string {
	val, err := c.client.Get("Q" + query).Result()
	if err != nil {
		return ""
	}
	return val
}

func (c *CacheInRedis) GetQuery(QID string) string {
	val, err := c.client.Get( "⌘" + QID).Result()
	if err != nil {
		return ""
	}
	return val
}

func (c *CacheInRedis) GetDataSize(QID string) int64 {
	val, err := c.client.Get("I" + QID).Result()
	if err != nil {
		return 0
	}
	n, err := strconv.ParseInt(val, 10, 64)
	if err == nil {
		return n
	}
	return 0
}

func (c *CacheInRedis) GetQueryLogOfUser(user string) [][]string {
	qids, err := c.client.Do("LRANGE" ,"q" + user, 0 , 9).Result()
	if err != nil {
		return [][]string{}
	}
	res :=[][]string{}
	if values, ok := qids.([]interface{}); ok {
		for _, qidTime := range values {
			res = append(res, strings.Split(qidTime.(string), "\t"))
		}
	}
	return res
}

func ppp(l []string){
	for i, v := range l{
		fmt.Println(i, v)
	}
}

func (c *CacheInRedis) GetAllUsers() []string {
	stringSliceCmd := c.client.Keys("M*")
	if stringSliceCmd.Err() != nil{
		return []string{}
	}
	users :=[]string{}
	for _, u := range stringSliceCmd.Val() {
		users = append(users, u[1:])
	}
	return users
}

func (c *CacheInRedis) PrintStatsForUser(user string) {
	var saved, spent, hit, miss, dataSaving, dataPaid string
	saved, err := c.client.Get("M"+user).Result()
	if err != nil {
		return
	}
	spent, err = c.client.Get("m"+user).Result()
	if err != nil {
		return
	}
	miss, err = c.client.Get("c"+user).Result()
	if err != nil {
		return
	}
	hit, err = c.client.Get("C"+user).Result()
	if err != nil {
		return
	}
	dataSaving, err = c.client.Get("S"+user).Result()
	if err != nil {
		return
	}
	dataPaid, err = c.client.Get("s"+user).Result()
	if err != nil {
		dataPaid = "0"
	}

	spentUSD, _ := strconv.ParseFloat(strings.TrimSpace(spent), 64)
	savedUSD, _ := strconv.ParseFloat(strings.TrimSpace(saved), 64)
	rate := int(100 * savedUSD/(savedUSD + spentUSD))
	print("user: " + user + ", ")
	fmt.Printf("saving_rate: %d%%, ", rate)
	print("money_spent: " + spent + "USD, ")
	print("money_saved: " + saved + "USD, ")
	print("data_saved: " + dataSaving + ", ")
	print("data_paid: " + dataPaid + ", ")
	print("hit: " + hit + ", ")
	println("miss: " + miss)
}

