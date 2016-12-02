package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
)

var cc cache.Cache

func InitCache() {
	cacheConfig := beego.AppConfig.String("cache")
	cc = nil

	if "redis" == cacheConfig {
		initRedis()
	} else {
		initMemcache()
	}
}

func initRedis() {
	var err error

	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()

	cc, err = cache.NewCache("redis", `{"conn":"`+beego.AppConfig.String("redis_host")+`"}`)
	if err != nil {
		beego.Info(err)
	}
}

func initMemcache() {
	var err error
	cc, err = cache.NewCache("memcache", `{"conn":"`+beego.AppConfig.String("memcache_host")+`"}`)

	if err != nil {
		beego.Info(err)
	}
}

func SetCache(key string, value interface{}, timeout int) error {
	data, err := Encode(value)
	if err != nil {
		return err
	}
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()
	timeouts := time.Duration(timeout) * time.Second
	err = cc.Put(key, data, timeouts)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func GetCache(key string, to interface{}) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()

	data := cc.Get(key)
	if data == nil {
		return errors.New("cache is not exists.")
	}
	err := Decode(data.([]byte), to)
	if err != nil {
		beego.Info("get cache fail", key, err)
	} else {
		beego.Info("get cache success.", key)
	}

	return err
}

func DelCache(key string) error {
	if cc == nil {
		return errors.New("cc is nil")
	}

	defer func() {
		if r := recover(); r != nil {
			cc = nil
		}
	}()

	err := cc.Delete(key)
	if err != nil {
		return errors.New("cache delete fail.")
	} else {
		return nil
	}
}

//Encode
//use gob encode data
//
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//Decode
//use gob decode data
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}
