package modules

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/hfdend/cxz/cli"
	"github.com/hfdend/cxz/conf"
	"github.com/hfdend/cxz/errors"
	"github.com/hfdend/cxz/models"
	"github.com/hfdend/cxz/utils"
	"github.com/hfdend/cxz/wx/miniprogram"
)

type passport int

var Passport passport

const (
	KEY_VerificationCode = "verification_code_"
	KEY_Register         = "register"
)

func (p passport) Register(phone, code, password string) (*models.User, error) {
	if err := p.CheckVerificationCode(phone, code, KEY_Register); err != nil {
		return nil, err
	}
	user := new(models.User)
	user.Phone = phone
	user.Password = utils.EncodePassword(password)
	if n, err := user.Insert(); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, errors.New("手机号已注册")
	}
	return user, nil
}

func (p passport) SendRegisterCode(phone string) (code string, err error) {
	var u *models.User
	if u, err = models.UserDefault.GetByPhone(phone); err != nil {
		return
	} else if u.ID != 0 {
		err = errors.New("此号码以及被注册")
		return
	}
	code = "1234"
	if conf.Config.Aliyun.SMS.Test {
		code = "1234"
	} else {
		code = fmt.Sprintf("%0.4d", utils.RandInterval(0, 10000))
		if err = SMS.Send(phone, code); err != nil {
			return
		}
	}

	if err = p.SaveVerificationCode(phone, code, KEY_Register, 10*time.Minute); err != nil {
		return
	}
	return
}

func (passport) Login(phone, password string) (token *models.Token, err error) {
	var user *models.User
	if user, err = models.UserDefault.GetByPhone(phone); err != nil {
		return
	} else if user == nil || user.ID == 0 {
		err = errors.New("该手机号未注册")
		return
	}
	if utils.EncodePassword(password) != user.Password {
		err = errors.New("密码错误")
		return
	}
	token = models.NewToken(time.Now().Add(7 * 24 * time.Hour).Local())
	if err = token.Clean(user.ID); err != nil {
		return
	}
	if err = token.SaveUser(user.ID); err != nil {
		return
	}
	return
}

func (passport) LoginByJsCode(code string) (token *models.Token, err error) {
	c := conf.Config.MiniProgram
	var session miniprogram.Session
	if session, err = miniprogram.GetSession(c.AppID, c.Secret, code); err != nil {
		return
	}
	var user *models.User
	if user, err = models.UserDefault.GetByOpenID(session.OpenID); err != nil {
		return
	} else if user == nil || user.ID == 0 {
		user.OpenID = session.OpenID
		// 生成一个唯一临时手机占位号
		user.Phone = fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
		var n int64
		if n, err = user.Insert(); err != nil {
			return
		} else if n == 0 {
			err = errors.New("unionid repetition")
			return
		}
	}
	token = models.NewToken(time.Now().Add(7 * 24 * time.Hour).Local())
	if err = token.Clean(user.ID); err != nil {
		return
	}
	if err = token.SaveUser(user.ID); err != nil {
		return
	}
	return
}

func (passport) BindPhone(userID int, phone, code string) error {
	// TODO CODE
	user, err := models.UserDefault.GetByPhone(phone)
	if err != nil {
		return err
	} else if user.ID != 0 {
		return errors.New("手机号已注册")
	}
	return models.UserDefault.UpdatePhone(userID, phone)
}

func (passport) SaveVerificationCode(phone, code, typ string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s_%s", KEY_VerificationCode, typ, phone)
	return cli.Redis.Set(key, code, expiration).Err()
}

func (p passport) CheckVerificationCode(phone, code, typ string) error {
	sCode, err := p.GetVerificationCode(phone, typ)
	if err != nil {
		return err
	}
	if sCode != code {
		return errors.New("验证码错误")
	}
	return nil
}

func (passport) GetVerificationCode(phone, typ string) (code string, err error) {
	key := fmt.Sprintf("%s%s_%s", KEY_VerificationCode, typ, phone)
	if code, err = cli.Redis.Get(key).Result(); err != nil && err == redis.Nil {
		err = errors.New("验证码已过期")
	}
	return
}
