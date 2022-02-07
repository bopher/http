package session

import (
	"encoding/json"
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/caster"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

type hSession struct {
	// cache driver
	cache cache.Cache
	// ctx request context
	ctx *fiber.Ctx
	// < 0 means 24 hours
	// > 0 is the time.Duration which the session should expire.
	expiration time.Duration
	// Session id generator
	generator func() string
	// header name
	name string
	// cache key
	key  string
	data map[string]interface{}
}

func (this hSession) err(
	pattern string,
	params ...interface{},
) error {
	return utils.TaggedError([]string{"HeaderSession"}, pattern, params...)
}

func (this *hSession) init(
	cache cache.Cache,
	ctx *fiber.Ctx,
	exp time.Duration,
	generator func() string,
	name string,
) {
	this.cache = cache
	this.ctx = ctx
	this.expiration = exp
	this.generator = generator
	this.name = name
	if this.name == "" {
		this.name = "X-SESSION-ID"
	}
	this.data = make(map[string]interface{})
}

func (this hSession) id() string {
	return "C_S_" + this.key
}

func (this hSession) ID() string {
	return this.key
}

func (this hSession) Context() *fiber.Ctx {
	return this.ctx
}

func (this *hSession) Parse() error {
	this.key = this.ctx.Get(this.name)
	exists := false
	var err error
	if this.key != "" {
		exists, err = this.cache.Exists(this.id())
		if err != nil {
			return this.err(err.Error())
		}
	}

	if !exists {
		return this.Regenerate()
	} else {
		res := make(map[string]interface{})
		caster, err := this.cache.Cast(this.id())
		if err != nil {
			return this.err(err.Error())
		}

		str, err := caster.String()
		if err != nil {
			return this.err(err.Error())
		}

		err = json.Unmarshal([]byte(str), &res)
		if err != nil {
			return this.err(err.Error())
		}

		this.data = res
		return nil
	}
}

func (this *hSession) Regenerate() error {
	err := this.Destroy()
	if err != nil {
		return err
	}

	this.key = this.generator()
	this.ctx.Set(this.name, this.key)
	return nil
}

func (s *hSession) Set(key string, value interface{}) {
	s.data[key] = value
}

func (this hSession) Get(key string) interface{} {
	return this.data[key]
}

func (this *hSession) Delete(key string) {
	delete(this.data, key)
}

func (this hSession) Exists(key string) bool {
	_, ok := this.data[key]
	return ok
}

func (this hSession) Cast(key string) caster.Caster {
	return caster.NewCaster(this.data[key])
}

func (this *hSession) Destroy() error {
	err := this.cache.Forget(this.id())
	if err != nil {
		return this.err(err.Error())
	}
	this.key = ""
	this.data = make(map[string]interface{})
	return nil
}

func (this hSession) Save() error {
	if this.key == "" {
		return nil
	}

	data, err := json.Marshal(this.data)
	if err != nil {
		return this.err(err.Error())
	}

	exists, err := this.cache.Set(this.id(), string(data))
	if err != nil {
		return this.err(err.Error())
	}

	if !exists {
		exp := this.expiration
		if exp <= 0 {
			exp = 24 * time.Hour
		}

		err = this.cache.Put(this.id(), string(data), exp)
		if err != nil {
			return this.err(err.Error())
		}
	}
	return nil
}
