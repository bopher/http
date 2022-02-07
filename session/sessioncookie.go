package session

import (
	"encoding/json"
	"time"

	"github.com/bopher/cache"
	"github.com/bopher/caster"
	"github.com/bopher/utils"
	"github.com/gofiber/fiber/v2"
)

// cSession cookie bases session
type cSession struct {
	// cache driver
	cache cache.Cache
	// ctx request context
	ctx *fiber.Ctx
	// secure attribute for cookie
	secure bool
	// domain attribute for cookie
	domain string
	// sameSite attribute for cookie
	// Possible values: "Lax", "Strict", "None"
	// Optional. Default: "Lax"
	sameSite string
	// < 0 means when browser closes
	// > 0 is the time.Duration which the session cookies should expire.
	expiration time.Duration
	// Session id generator
	generator func() string
	// cookie name
	name string
	// cache key
	key  string
	data map[string]interface{}
}

func (this cSession) err(
	pattern string,
	params ...interface{},
) error {
	return utils.TaggedError([]string{"CookieSession"}, pattern, params...)
}

func (this *cSession) init(
	cache cache.Cache,
	ctx *fiber.Ctx,
	secure bool,
	domain string,
	sameSite string,
	exp time.Duration,
	generator func() string,
	name string,
) {
	this.cache = cache
	this.ctx = ctx
	this.secure = secure
	this.domain = domain
	this.sameSite = sameSite
	this.expiration = exp
	this.generator = generator
	this.name = name
	if this.name == "" {
		this.name = "session"
	}
	this.data = make(map[string]interface{})
}

func (this cSession) id() string {
	return "C_S_" + this.key
}

func (this cSession) ID() string {
	return this.key
}

func (this cSession) Context() *fiber.Ctx {
	return this.ctx
}

func (this *cSession) Parse() error {
	this.key = this.ctx.Cookies(this.name)
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

func (this *cSession) Regenerate() error {
	err := this.Destroy()
	if err != nil {
		return err
	}

	this.key = this.generator()
	cookie := fiber.Cookie{}
	cookie.Name = this.name
	cookie.Value = this.key
	cookie.Secure = this.secure
	cookie.Domain = this.domain
	cookie.SameSite = this.sameSite
	if this.expiration > 0 {
		cookie.Expires = time.Now().UTC().Add(this.expiration)
	}
	this.ctx.Cookie(&cookie)
	return nil
}

func (this *cSession) Set(key string, value interface{}) {
	this.data[key] = value
}

func (this cSession) Get(key string) interface{} {
	return this.data[key]
}

func (this *cSession) Delete(key string) {
	delete(this.data, key)
}

func (this cSession) Exists(key string) bool {
	_, ok := this.data[key]
	return ok
}

func (this cSession) Cast(key string) caster.Caster {
	return caster.NewCaster(this.data[key])
}

func (this *cSession) Destroy() error {
	err := this.cache.Forget(this.id())
	if err != nil {
		return this.err(err.Error())
	}
	this.key = ""
	this.data = make(map[string]interface{})
	return nil
}

func (this cSession) Save() error {
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
