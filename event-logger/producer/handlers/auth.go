package handlers

import (
	"errors"
	"event-logger/internal/data"
	"event-logger/internal/lib"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/segmentio/kafka-go"
)

var (
	ErrMissingSignUpAttr = errors.New("Username, email, and password are required")
	ErrMissingLoginAttr  = errors.New("Username, and password are required")
	ErrInsecurePassword  = errors.New("Password must be at least 8 characters")
	ErrUserTaken         = errors.New("Username and email must be unique")
	ErrUnknownUser       = errors.New("Password or username are incorrect")
	ErrInvalidSession    = errors.New("Session id is invalid or expired")
)

func AddUser(user *data.UserParams) error {

	if user.Email == "" || user.Username == "" || user.Password == "" {
		return ErrMissingSignUpAttr
	}

	if len(user.Password) < 8 {
		return ErrInsecurePassword
	}

	for _, u := range data.Users {
		if u.Username == user.Username || u.Email == user.Email {
			return ErrUserTaken
		}
	}

	data.Users[data.Seq] = user

	return nil
}

func FindUser(signin *data.SignInProps) (error, *data.UserParams) {
	if signin.Username == "" || signin.Password == "" {
		return ErrMissingLoginAttr, nil
	}

	for _, u := range data.Users {
		if u.Username == signin.Username {
			if u.Password != signin.Password {
				return ErrUnknownUser, nil
			}

			return nil, u
		}
	}

	return ErrUnknownUser, nil
}

func AddSession(u *data.UserParams) string {
	id := lib.RandStringRunes(32)

	// Purposefully does not store password in session
	data.Sessions[id] = &data.UserType{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}

	return id
}

func FindUserBySession(id string) *data.UserType {
	for i, s := range data.Sessions {
		if i == id {
			return s
		}
	}

	return nil
}

func (h *Handler) Register(c *echo.Context) error {

	lock.Lock()
	defer lock.Unlock()

	var user data.UserParams

	hasWriter := h.Writer != nil

	if err := c.Bind(&user); err != nil {
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to register user %v", err.Error()),
			})
		}
		return c.String(http.StatusBadRequest, "")
	}

	if err := AddUser(&user); err != nil {
		m := err.Error()
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to register user %v", m),
			})
		}
		return c.String(http.StatusUnprocessableEntity, m)
	}

	if hasWriter {
		h.Writer.WriteMessages(kafka.Message{
			Value: fmt.Appendf(nil, "registered %s", user.Email),
		})
	}

	return c.String(http.StatusCreated, "")
}

func (h *Handler) Login(c *echo.Context) error {

	lock.Lock()
	defer lock.Unlock()

	var signin data.SignInProps
	hasWriter := h.Writer != nil

	if err := c.Bind(&signin); err != nil {
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to sign in user %v", err.Error()),
			})
		} else {
			fmt.Printf("Failed to sign in user %v", err.Error())
		}
		return c.String(http.StatusBadRequest, "")
	}

	err, u := FindUser(&signin)
	if err != nil {
		if hasWriter {
			h.Writer.WriteMessages(kafka.Message{
				Value: fmt.Appendf(nil, "Failed to sign in user %v", err.Error()),
			})
		} else {
			fmt.Printf("Failed to sign in user %v", err.Error())
		}
		return c.String(http.StatusBadRequest, "")
	}
	sessionId := AddSession(u)

	cookie := new(http.Cookie)
	cookie.Name = "authorization"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(24 * 30 * time.Hour)

	if hasWriter {
		h.Writer.WriteMessages(kafka.Message{
			Value: fmt.Appendf(nil, "New session for %s", u.Username),
		})
	}

	c.SetCookie(cookie)
	return c.String(http.StatusOK, "")

}

func (h *Handler) Me(c *echo.Context) error {

	lock.Lock()
	defer lock.Unlock()

	auth, err := c.Cookie("authorization")

	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid session")
	}

	user := FindUserBySession(auth.Value)
	if user == nil {
		return ErrInvalidSession
	}
	return c.JSON(http.StatusOK, user)
}
