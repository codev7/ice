package ice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

var AuthenticateUser func(token string, conn Conn) Identity

type Identity interface {
	UserId() int64
	UserName() string
	UserRole() string
	CheckRole(role string) bool
}

type UserBase struct {
	Id       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Active   bool   `json:"active" db:"active"`
	Password string `json:"-" db:"password"`
	Token    string `json:"token" db:"token"`
	Role     string `json:"role" db:"role"`
}

func (u *UserBase) UserId() int64 {
	return u.Id
}

func (u *UserBase) UserName() string {
	return u.Name
}

func (u *UserBase) UserRole() string {
	return u.Role
}

func (p *UserBase) GenerateToken() error {
	b, err := bcrypt.GenerateFromPassword([]byte(p.Password+p.Name+p.Email), 13)
	if err != nil {
		return err
	}
	p.Token = string(b)
	return nil
}

func (p *UserBase) SetPassword(pwd string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), 13)
	if err != nil {
		return err
	}
	p.Password = string(password)
	return nil
}

func (p *UserBase) ComparePassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

func (p *UserBase) CheckRole(role string) bool {
	if p != nil {
		if p.Role == role {
			return true
		}
	}
	return false
}

type AnyAuthorizedUser struct{}

func (r *AnyAuthorizedUser) Authorize(conn Conn) bool {
	return conn.User() != nil
}

type OnlyUnauthorizedUser struct{}

func (r *OnlyUnauthorizedUser) Authorize(conn Conn) bool {
	return conn.User() == (*UserBase)(nil)
}

type AnyAdminUser struct{}

func (r *AnyAdminUser) Authorize(conn Conn) bool {
	return conn.User().CheckRole("admin")
}

func HmacToken(message interface{}, key []byte) ([]byte, error) {
	if len(key) == 0 {
		key = []byte(Config.Secret)
	}
	data, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	out := mac.Sum(data)
	//out = append(out,data ...)
	return out, nil
}
