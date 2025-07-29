package models

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey;autoIncrement:true;unique" json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email" gorm:"unique"`
	Password       string    `json:"-"`
	DeviceId       string    `json:"-"`
	ProfilePhoto   string    `json:"profilePhoto"`
	Age            string    `json:"age"`
	Country        string    `json:"country"`
	EducationLevel string    `json:"educationLevel"`
	TestedBefore   bool      `json:"testedBefore"`
	SaveResults    bool      `json:"SaveResults"`
	Gender         string    `json:"gender"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// FirstName ...
func (u *GoogleUserInfo) FirstName() string {
	return u.GivenName
}

// LastName ...
func (u *GoogleUserInfo) LastName() string {
	return u.FamilyName
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func ValidateGoogleToken(accessToken string) (*GoogleUserInfo, error) {
	url := "https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=" + accessToken

	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer " + accessToken,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
