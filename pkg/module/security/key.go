package security

import (
	"context"
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/crypt"
	"github.com/davycun/eta/pkg/common/crypt/crypt_asym"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/locker"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

func GetPublicKey(algo string) string {
	key, err := GetAsymmetryKey(algo)
	if err != nil || !key.Valid() {
		logger.Errorf("get public key err %s", err)
		return ""
	}
	return key.PublicKey
}
func GetPrivateKey(algo string) string {
	key, err := GetAsymmetryKey(algo)
	if err != nil || !key.Valid() {
		logger.Errorf("get private key err %s", err)
		return ""
	}
	return key.PrivateKey
}

func GetAsymmetryKey(algo string) (crypt_asym.KeyPair, error) {

	var (
		prf        = os.Getenv("ETA_PRIVATE_KEY_FILE") //可以是密钥的内容，也可以是存储密钥的文件的地址
		pbf        = os.Getenv("ETA_PUBLIC_KEY_FILE")  //可以是公钥的内容，也可以是存储公钥的文件的地址
		expiration = time.Hour * 24
		err        error
		key        = crypt_asym.KeyPair{}
	)

	//1. 从缓存中获取
	_, err = cache.Get(getRedisKey(algo), &key)
	if err != nil || key.Valid() {
		return key, err
	}

	//2. 从环境变量中指定的文件中获取，如果环境变量直接存储的是公钥内容，那么直接返回
	//如果是通过环境变量直接配置内容的情况，那么直接返回（除非是配置的文件，避免每次都读文件才设置redis）
	key.PublicKey = pbf
	key.PrivateKey = prf
	if key.Valid() {
		return key, nil
	}

	//3. 环境变量可能存储的是公钥文件的地址，从文件中读取
	if pbf != "" && prf != "" {
		pubBytes, _ := os.ReadFile(pbf)
		priBytes, _ := os.ReadFile(prf)
		key.PublicKey = string(pubBytes)
		key.PrivateKey = string(priBytes)
	}

	if key.Valid() {
		_ = cache.SetEx(getRedisKey(algo), key, expiration)
		return key, nil
	}
	//4. 自动生成密钥需要通过分布式锁去设置
	redLock := locker.NewRedLock([]*redis.Client{global.GetRedis()})
	defer redLock.Release()
	lkOk := redLock.Acquire()
	if !lkOk {
		//如果没有获取到，可能已经被其他节点获取到锁，所以直接等待其他节点设置完，直接从缓存中获取即可
		return getKeyFromCache(algo, redLock.TTL)
	}
	//自动生成公钥密钥
	key, err = crypt.GenKeypair(algo, 2048)
	if err != nil || !key.Valid() {
		return key, err
	}
	//TODO 是否需要考虑设置超时SetEx，过一段时间后重新生成密钥
	err = cache.Set(getRedisKey(algo), key)
	return key, err
}

// 从缓存中获取密钥，如果获取到就返回值，如果直到超时都没有获取到就直接返回
func getKeyFromCache(algo string, ttl time.Duration) (crypt_asym.KeyPair, error) {
	var (
		key = crypt_asym.KeyPair{}
	)
	ok, err := cache.Get(getRedisKey(algo), &key)
	if ok {
		return key, err
	}
	ct, cancel := context.WithTimeout(context.Background(), ttl)
	defer cancel()
	tk := time.NewTicker(time.Second * 1)
	defer tk.Stop()

	//两秒一次去获取，如果获取到就返回值，如果直到超时都没有获取到就直接返回
	for {
		select {
		case <-ct.Done():
			_, err = cache.Get(getRedisKey(algo), &key)
			return key, err
		case <-tk.C:
			ok, err = cache.Get(getRedisKey(algo), &key)
			if ok {
				return key, err
			}
		}
	}
}

func getRedisKey(algo string) string {
	return fmt.Sprintf("eta:crypt:%s", algo)
}
