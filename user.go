package ice

import (
	"golang.org/x/crypto/bcrypt"
)

var AuthenticateUser func(token string, conn Conn) *User

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
	Password string `json:"-"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

func (p *User) GenerateToken() error {
	b, err := bcrypt.GenerateFromPassword([]byte(p.Password+p.Name+p.Email), 10)
	if err != nil {
		return err
	}
	p.Token = string(b)
	return nil
}

func (p *User) SetPassword(pwd string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	if err != nil {
		return err
	}
	p.Password = string(password)
	return nil
}

func (p *User) ComparePassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

func (p *User) CheckRole(role string) bool {
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
	return conn.User() == nil
}

type AnyAdminUser struct{}

func (r *AnyAdminUser) Authorize(conn Conn) bool {
	return conn.User().CheckRole("admin")
}
