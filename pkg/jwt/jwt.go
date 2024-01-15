package jwt

import (
	"errors"
	"github.com/eliofery/go-chix/pkg/chix"
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/log"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	CookieTokenName    = "jwt"
	ExpiresTimeDefault = 3600
)

// TokenManager jwt токен
type TokenManager interface {
	GenerateToken(userId int) (token string, err error)
	VerifyToken(token string) (issuer string, err error)
	GetExpiresTime() time.Duration
	SetCookieToken(ctx *chix.Ctx, token string)
	RemoveCookieToken(ctx *chix.Ctx)
}

type tokenManager struct {
	conf config.Config
}

func NewTokenManager(conf config.Config) TokenManager {
	log.Debug("Инициализация менеджера токенов")

	return &tokenManager{conf: conf}
}

// GenerateToken создание токена
func (t *tokenManager) GenerateToken(userId int) (token string, err error) {
	if err = t.isSecretEmpty(); err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": strconv.Itoa(userId),
		"exp": time.Now().Add(time.Second * t.GetExpiresTime()).Unix(),
	})

	token, err = claims.SignedString([]byte(t.conf.Get("jwt.secret")))
	if err != nil {
		log.Error("Не удалось создать токен", slog.String("err", err.Error()))
		return "", err
	}

	return token, nil
}

// VerifyToken валидация токена
func (t *tokenManager) VerifyToken(token string) (issuer string, err error) {
	if err = t.isSecretEmpty(); err != nil {
		return "", err
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Неожиданный метод подписи токена", slog.String("method", token.Method.Alg()))
			return nil, errors.New("неожиданный метод подписи токена")
		}

		return []byte(t.conf.Get("jwt.secret")), nil
	})
	if err != nil || !parsedToken.Valid {
		log.Error("Не верный токен", slog.String("err", err.Error()))
		return "", err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Не верный тип токена", slog.String("err", err.Error()))
		return "", err
	}

	issuer, err = claims.GetIssuer()
	if err != nil {
		log.Error("Issuer не получен", slog.String("err", err.Error()))
		return "", err
	}

	return issuer, nil
}

// GetExpiresTime получить время истечения токена
func (t *tokenManager) GetExpiresTime() time.Duration {
	expiresTime, ok := t.conf.GetAny("jwt.expires").(int)
	if !ok {
		log.Warn("Не удалось получить время истечения токена", slog.Any("expiresTime", expiresTime))
		expiresTime = ExpiresTimeDefault
	}

	return time.Duration(expiresTime)
}

// SetCookieToken сохранить токен в куки
func (t *tokenManager) SetCookieToken(ctx *chix.Ctx, token string) {
	cookie := http.Cookie{
		Name:     CookieTokenName,
		Value:    token,
		Expires:  time.Now().Add(time.Second * t.GetExpiresTime()),
		HttpOnly: true,
	}

	ctx.Cookie(&cookie)
}

// RemoveCookieToken удалить токен из куки
func (t *tokenManager) RemoveCookieToken(ctx *chix.Ctx) {
	ctx.ClearCookie(CookieTokenName)
}

// isSecretEmpty проверка секретного ключа
func (t *tokenManager) isSecretEmpty() error {
	secret := t.conf.Get("jwt.secret")

	if secret == "" {
		log.Error("Секретный ключ токена не может быть пустым", slog.String("secret", secret))
		return errors.New("секретный ключ токена пуст")
	}

	return nil
}
