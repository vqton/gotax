package gl

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ─── RSA Key Pair ───────────────────────────────────────────────

var (
	jwtPrivateKey *rsa.PrivateKey
	jwtPublicKey  *rsa.PublicKey
)

func GenerateRSAKeyPair(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	jwtPrivateKey = privateKey
	jwtPublicKey = &privateKey.PublicKey
	return nil
}

func SetJWTSecret(secret string) {
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	if err := GenerateRSAKeyPair(2048); err != nil {
		log.Fatalf("failed to generate RSA key pair: %v", err)
	}
}

// ─── Claims ─────────────────────────────────────────────────────

type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
	TokenID  string   `json:"jti,omitempty"`
	jwt.RegisteredClaims
}

// ─── Password Hashing ───────────────────────────────────────────

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ─── Token Generation ───────────────────────────────────────────

func GenerateToken(user *User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("fallback-secret"))
}

func GenerateAccessToken(user *User) (string, error) {
	if jwtPrivateKey == nil {
		return "", errors.New("JWT key pair not initialized")
	}
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(jwtPrivateKey)
}

func GenerateRefreshTokenRaw(userID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashRefreshToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

func GeneratePasswordResetTokenRaw() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ─── Token Parsing ──────────────────────────────────────────────

func ParseAndValidateAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtPublicKey, nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	return claims, nil
}

// ─── Rate Limiter ───────────────────────────────────────────────

type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*attemptInfo
	limit    int
	window   time.Duration
}

type attemptInfo struct {
	count   int
	windowStart time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		attempts: make(map[string]*attemptInfo),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	info, ok := rl.attempts[key]
	if !ok || now.Sub(info.windowStart) > rl.window {
		rl.attempts[key] = &attemptInfo{count: 1, windowStart: now}
		return true
	}
	if info.count >= rl.limit {
		return false
	}
	info.count++
	return true
}

// ─── Middlewares ─────────────────────────────────────────────────

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		claims, err := ParseAndValidateAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

func RoleMiddleware(roles ...UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleStr, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		role := UserRole(roleStr.(string))
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

func GetUserID(c *gin.Context) string {
	uid, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	s, ok := uid.(string)
	if !ok {
		return ""
	}
	return s
}

func GetUserRole(c *gin.Context) UserRole {
	roleStr, exists := c.Get("role")
	if !exists {
		return ""
	}
	return UserRole(roleStr.(string))
}

func GetUsername(c *gin.Context) string {
	u, exists := c.Get("username")
	if !exists {
		return ""
	}
	return u.(string)
}

// ─── TOTP ───────────────────────────────────────────────────────

func GenerateTOTPSecret() string {
	b := make([]byte, 20)
	rand.Read(b)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func GenerateTOTPURL(username, secret string) string {
	return fmt.Sprintf("otpauth://totp/GoTax:%s?secret=%s&issuer=GoTax&algorithm=SHA1&digits=6&period=30", username, secret)
}

func VerifyTOTP(secret, code string) bool {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}
	counter := time.Now().Unix() / 30
	for i := int64(-1); i <= 1; i++ {
		if generateTOTP(key, counter+i) == code {
			return true
		}
	}
	return false
}

func generateTOTP(key []byte, counter int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	hash := mac.Sum(nil)
	offset := hash[len(hash)-1] & 0x0f
	code := (int32(hash[offset]&0x7f) << 24) |
		(int32(hash[offset+1]&0xff) << 16) |
		(int32(hash[offset+2]&0xff) << 8) |
		(int32(hash[offset+3] & 0xff))
	code = code % 1000000
	return fmt.Sprintf("%06d", code)
}

func CheckBackupCode(user *User, code string) bool {
	for i, bc := range user.BackupCodes {
		if bc == code {
			user.BackupCodes = append(user.BackupCodes[:i], user.BackupCodes[i+1:]...)
			return true
		}
	}
	return false
}
