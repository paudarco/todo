package service

import (
	"crypto/sha256"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paudarco/todo/internal/config"
	"github.com/paudarco/todo/internal/errors"
	"github.com/paudarco/todo/internal/models"
)

// Структура для сервиса аутентификации.
type AuthService struct {
	cfg            *config.Config
	tokenExpiry    time.Duration
	secret         []byte
	hashedPassword [32]byte
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		cfg:         cfg,
		tokenExpiry: time.Duration(cfg.TTL) * time.Hour,
		secret:      []byte(cfg.Secret),
	}
}

// SignIn проверяет введенный пароль и возвращает токен, если пароль верен.
func (s *AuthService) SignIn(user models.User) (models.AuthResponse, error) {
	// Если пароль неверный - возвращаем ошибку.
	if user.Password != s.cfg.Password {
		return models.AuthResponse{}, errors.ErrInvalidPassword
	}

	// Хешируем пароль.
	hashedPassword := sha256.Sum256([]byte(user.Password))

	s.hashedPassword = hashedPassword

	// Генерируем токен JWT.
	token, err := s.GenerateToken(hashedPassword)
	if err != nil {
		return models.AuthResponse{}, err
	}

	response := models.AuthResponse{
		Token: token,
	}

	return response, nil
}

// GenerateToker возвращает сгенерированый JWT токен.
func (s *AuthService) GenerateToken(hashedPassword [32]byte) (string, error) {
	// Создаем JWT токен с данными payload'а.
	// Раз уж пароль передается в JSON незашифрованным,
	// для создания уникального токена используем его в payload.
	// Увы, тесты требуют постоянный токен.
	payload := jwt.MapClaims{
		"sub": s.hashedPassword,
		// "exp": time.Now().Add(time.Duration(s.tokenExpiry) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString(s.secret)
}

// ValidateToken проверяет валидность JWT токена и возвращает nil, если токен валиден.
func (s *AuthService) ValidateToken(accessToken string) error {
	// Парсим токен и проверяем метод подписи.
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrInvalidSigningMethod
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return err
	}

	// Проверяем валидность payload'а JWT токена.
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.ErrInvalidClaims
	}

	// Проверяем, что срок хранения из payload приводится к типу string.
	// expStr, ok := claims["exp"].(string)
	// if !ok {
	// 	return err
	// }

	// Приводим строку со сроком хранения к числовому типу для получения UNIX формата.
	// exp, err := strconv.ParseInt(expStr, 10, 64)
	// if err != nil {
	// 	return err
	// }

	// ПРоверяем, не истек ли срок действия в токене.
	// if time.Now().Unix() > exp {
	// 	return errors.ErrExpiredToken
	// }

	return nil
}
