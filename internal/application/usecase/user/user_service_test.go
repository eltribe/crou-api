package user

import (
	"crou-api/config/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserService Suite")

}

var _ = Describe("UserService 비즈니스 로직 테스트", func() {
	var (
		userService *UserService
		app         *fiber.App
		ctx         *fiber.Ctx
		fakedb      database.Persistent
	)
	BeforeEach(func() {
		fakedb = test.FakeDatabase()
		jwtService := common.NewJwtService(test.FakeConfig)
		fileService := common.NewFileService(test.FakeConfig, fakedb)
		myPromptService := NewMyPromptService(fakedb, jwtService, fileService, config.NewKopenAiGpt(test.FakeConfig))
		notificationService := notification.NewNotificationService(fakedb, jwtService)
		userService = NewUserService(fakedb, jwtService, myPromptService, notificationService)

		app = fiber.New()
		ctx = app.AcquireCtx(&fasthttp.RequestCtx{})
		claims := jwt.MapClaims{
			"email":    "test@email.com",
			"type":     "google",
			"sub":      "test",
			"nickname": "test",
		}
		ctx.Locals("claims", claims)

	})

	Context("GetUser ", func() {
		It("유저가 존재하지 않을때", func() {
			user, err := userService.GetUser(ctx)
			Expect(user).To(BeNil())
			Expect(err).ToNot(BeNil())
		})

		It("등록된 유저가 존재할 경우", func() {
			fakedb.DB().Save(&domains.User{
				Nickname:   "test",
				OauthType:  "google",
				OauthSub:   "test",
				OauthEmail: "test@email.com",
			})
			user, err := userService.GetUser(ctx)
			Expect(user).ToNot(BeNil())
			Expect(err).To(BeNil())
			Expect(user.Nickname).To(Equal("test"))
			Expect(user.OauthEmail).To(Equal("test@email.com"))
		})
	})

	Context("UpdateUser", func() {
		It("유저 닉네임 변경 완료", func() {
			fakedb.DB().Create(&domains.User{
				Nickname:   "test",
				OauthType:  "google",
				OauthSub:   "test",
				OauthEmail: "test@email.com",
				Taste:      "a,b,c",
			})
			user, err := userService.UpdateUser(ctx, &messages.UpdateUserRequest{
				Nickname: "닉네임변경",
			})
			Expect(user.Nickname).To(Equal("닉네임변경"))
			Expect(user.Taste).To(Equal("a,b,c"))
			Expect(err).To(BeNil())
		})
	})

	AfterEach(func() {
		app.ReleaseCtx(ctx)
		test.CloseFakeDatabase(fakedb)
	})
})
