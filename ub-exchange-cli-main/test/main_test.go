package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/test/data/seed"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	sarulabsDI "github.com/sarulabs/di"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//gormLogger "gorm.io/gorm/logger"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const DbDSN = "exchange_user:exchange_pass@tcp(db:3306)/exchange_go_test?multiStatements=true&parseTime=true"
const RedisAddr = "redis:6379"

var container sarulabsDI.Container
var mainDb *gorm.DB
var redisClient *redis.Client
var testUser *userActor

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

type userActor struct {
	ID    int
	Token string
	Email string
}

var Seeders = map[string]func(db *gorm.DB){
	"currencySeeder":       seed.CurrencySeed,
	"externalExchangeSeed": seed.ExternalExchangeSeed,
	"userLevelSeed":        seed.UserLevelSeed,
	"userPermissionSeed":   seed.UserPermissionSeed,
	"roleSeed":             seed.RoleSeed,
}

func getContainer() sarulabsDI.Container {
	if container != nil {
		return container
	}
	c := di.NewContainer()

	container = c
	return container
}

func getDb() *gorm.DB {
	if mainDb != nil {
		return mainDb
	}
	db, err := gorm.Open(mysql.Open(DbDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		//Logger:                 gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		panic("can not establish connection to database")
	}
	mainDb = db
	return mainDb
}

func getRedis() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	rc := redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: "",
		DB:       0,
	})

	//s := miniredis.NewMiniRedis()
	//defer s.Close()
	//_ = s.Start()
	//rc := redis.NewClient(&redis.Options{Addr: s.Addr()})

	redisClient = rc
	return redisClient
}

func TestMain(m *testing.M) {
	setupDb()
	runSeeders()
	exitCode := m.Run()
	os.Exit(exitCode)

}

func setupDb() {
	db := getDb()
	content, err := os.ReadFile("./data/db.sql")

	err = db.Exec(string(content)).Error
	if err != nil {
		fmt.Println(err.Error())
	}
}

func runSeeders() {
	db := getDb()
	for name, seedFunc := range Seeders {
		fmt.Println("running seed ", name)
		seedFunc(db)
	}
}

func getUserActor() *userActor {
	if testUser != nil {
		return testUser
	}
	db := getDb()
	passwordEncoder := platform.NewPasswordEncoder()
	encodedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	u := user.User{
		Email:               "test@test.com",
		Password:            string(encodedPassword),
		Kyc:                 user.KycLevelMinimum,
		Status:              "VERIFIED",
		AccountStatus:       "UNBLOCKED",
		ExchangeNumber:      12,
		Google2faSecretCode: sql.NullString{String: "HWOAQZBGXCKJZQVH", Valid: true},
		IsTwoFaEnabled:      true,
		UbID:                "",
		VerificationCode:    "",
		PrivateChannelName:  "userActorPrivateChannel",
		UserLevelID:         1,
	}
	err := db.Create(&u).Error
	if err != nil {
		fmt.Println("err", err)
	}

	up := &user.Profile{
		UserID:     u.ID,
		TrustLevel: 0,
	}

	err = db.Create(&up).Error
	if err != nil {
		fmt.Println("err", err)
	}

	twoFaCode, _ := totp.GenerateCode("HWOAQZBGXCKJZQVH", time.Now())

	res := httptest.NewRecorder()
	body := `{"username":"test@test.com","password":"123456789","2fa_code":"` + twoFaCode + `","recaptcha":"recaptcha"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(body)))
	c := getContainer()
	httpServer := c.Get(di.HTTPServer).(api.HTTPServer)
	engine := httpServer.GetEngine()
	engine.ServeHTTP(res, req)

	response := struct {
		Token string
	}{}
	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		fmt.Println("err", err)
	}

	ua := userActor{
		ID:    u.ID,
		Token: response.Token,
		Email: "test@test.com",
	}

	testUser = &ua
	return testUser
}

func getAdminUserActor() *userActor {
	db := getDb()
	userActor := getNewUserActor()
	userRole := user.UserRole{
		UserID: int64(userActor.ID),
		RoleID: 1, //super admin see the roleSeeder
	}
	db.Create(userRole)
	return userActor
}

func getNewUserActor() *userActor {
	db := getDb()
	passwordEncoder := platform.NewPasswordEncoder()
	encodedPassword, _ := passwordEncoder.GenerateFromPassword("123456789")
	email := uuid.NewString() + "@test.com"
	ubID := randSeq(10)
	verificationCode := randSeq(12)
	PrivateChannelName := randSeq(14)
	u := user.User{
		Email:              email,
		Password:           string(encodedPassword),
		Kyc:                user.KycLevelMinimum,
		Status:             "VERIFIED",
		AccountStatus:      "UNBLOCKED",
		ExchangeNumber:     12,
		IsTwoFaEnabled:     false,
		UbID:               ubID,
		VerificationCode:   verificationCode,
		PrivateChannelName: PrivateChannelName,
		UserLevelID:        1,
	}
	err := db.Create(&u).Error
	if err != nil {
		fmt.Println("err", err)

	}

	res := httptest.NewRecorder()
	body := "{" +
		"\"username\":\"" + email + "\"," +
		"\"password\": \"123456789\"," +
		"\"recaptcha\": \"recaptcha\"" +
		"}"

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(body)))

	c := getContainer()
	httpServer := c.Get(di.HTTPServer).(api.HTTPServer)
	engine := httpServer.GetEngine()
	engine.ServeHTTP(res, req)

	response := struct {
		Token string
	}{}

	err = json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		fmt.Println("err", err)
	}

	ua := userActor{
		ID:    u.ID,
		Token: response.Token,
		Email: email,
	}

	return &ua
}
