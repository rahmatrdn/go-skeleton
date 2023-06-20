package middleware

import (
	"context"
	"log"
	"strings"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/config"
	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/entity"

	pb "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/proto"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apperr "gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/error"
)

func Authenticator(c *fiber.Ctx) error {
	log.Println("HIT MASUK KE GRPC")

	cfg := config.NewConfig()

	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		panic(apperr.UnauthorizedError{
			Message: "You are not logged in",
		})
	}
	conn, err := grpc.Dial(cfg.MiddlewareAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("failed connect on " + cfg.MiddlewareAddress + " " + err.Error())
	}
	defer conn.Close()

	cl := pb.NewValidateTokenServiceClient(conn)
	reqs := &pb.ValidateTokenRequest{
		Token: tokenString,
		Path:  c.Route().Path,
	}
	resp, err := cl.Validate(context.Background(), reqs)
	if err != nil {
		panic(apperr.InternalError{
			Message: "Cannot connect to middleware",
		})
	}
	if !resp.Val {
		panic(apperr.UnauthorizedError{
			Message: resp.Message,
		})
	}
	uid := int64(resp.User.Id)
	iuid := int(uid)

	ba := int64(resp.User.Branchalias)
	uba := int(ba)

	lfc := int64(resp.User.LoginFailedCount)
	ulfc := int(lfc)

	stat := int64(resp.User.Status)
	ustat := int(stat)
	var datauser entity.UserMe
	datauser.Access_token = &resp.User.AccessToken
	datauser.IdUser = &iuid
	datauser.Npp = &resp.User.Npp
	datauser.Branchalias = &uba
	datauser.Role = &resp.User.Role
	datauser.RoleId = &resp.User.RoleId
	datauser.Title = &resp.User.Title
	datauser.Full_name = &resp.User.FullName
	datauser.Email = &resp.User.Email
	datauser.Access_token_expire = &resp.User.AccessTokenExpire
	datauser.Login_failed_expire = &resp.User.LdapExpire
	datauser.Login_failed_count = &ulfc
	datauser.Login_failed_expire = &resp.User.LoginFailedExpire
	datauser.First_login_at = &resp.User.FirstLoginAt
	datauser.Last_login_at = &resp.User.LastLoginAt
	datauser.Status = &ustat
	c.Locals("user", &datauser)
	return c.Next()
}

// func Authenticator(c *fiber.Ctx) error {
// 	var tokenString string
// 	cfg := config.New()

// 	authorization := c.Get("Authorization")

// 	if strings.HasPrefix(authorization, "Bearer ") {
// 		tokenString = strings.TrimPrefix(authorization, "Bearer ")
// 	} else if c.Cookies("token") != "" {
// 		tokenString = c.Cookies("token")
// 	}

// 	if tokenString == "" {
// 		panic(apperr.UnauthorizedError{
// 			Message: "You are not logged in",
// 		})
// 	}

// 	user := GetToken(cfg, tokenString)

// 	if user.Access_token_expire != nil && *helper.String(helper.GetDateNow()) >= *user.Access_token_expire {
// 		panic(apperr.UnauthorizedError{
// 			Message: "Your token is expired.",
// 		})
// 	}

// 	if user.Ldap_expire != nil && *helper.String(helper.GetDateNow()) >= *user.Ldap_expire {
// 		LDAPauthentication(cfg, *user.IdUser)
// 	}

// 	datauser := UpdateUsers(cfg, tokenString)

// 	c.Locals("user", &datauser)

// 	return c.Next()
// }

// func LDAPauthentication(cfg config.Config, id_user int) error {

// 	database := config.MainDatabase(cfg)

// 	log.Println("================ HIT LDAP ===============")

// 	sqlStatement := `
// 	UPDATE users
// 	SET ldap_expire = ?, updated_at = now()
// 	WHERE id = ?;`

// 	_, err := database.Exec(sqlStatement, helper.AddMinutes(120), id_user)
// 	if err != nil {
// 		panic(apperr.InternalError{
// 			Message: err.Error(),
// 		})
// 	}

// 	return nil
// }

// func GetToken(cfg config.Config, tokenString string) entity.UserMe {

// 	database := config.MainDatabase(cfg)
// 	var userResult entity.UserMe
// 	err := database.QueryRow("select id, npp, branchalias, role, title, full_name, email, access_token, access_token_expire, ldap_expire, login_failed_count, "+
// 		"login_failed_expire, first_login_at, last_login_at, status FROM users WHERE access_token = ? AND last_login_at IS NULL", tokenString).
// 		Scan(
// 			&userResult.IdUser,
// 			&userResult.Npp,
// 			&userResult.Branchalias,
// 			&userResult.Role,
// 			&userResult.Title,
// 			&userResult.Full_name,
// 			&userResult.Email,
// 			&userResult.Access_token,
// 			&userResult.Access_token_expire,
// 			&userResult.Ldap_expire,
// 			&userResult.Login_failed_count,
// 			&userResult.Login_failed_expire,
// 			&userResult.First_login_at,
// 			&userResult.Last_login_at,
// 			&userResult.Status,
// 		)

// 	if err != nil {
// 		panic(apperr.InternalError{
// 			Message: "Invalid Token",
// 		})
// 	}

// 	return userResult
// }

// func UpdateUsers(cfg config.Config, tokenString string) entity.UserMe {

// 	database := config.MainDatabase(cfg)
// 	sqlStatement := `
// 		UPDATE users
// 		SET access_token_expire = ?
// 		WHERE access_token = ?;`

// 	_, err := database.Exec(sqlStatement, helper.AddMinutes(60), tokenString)
// 	if err != nil {
// 		panic(apperr.InternalError{
// 			Message: err.Error(),
// 		})
// 	}

// 	user := GetToken(cfg, tokenString)

// 	return user
// }
