package entity

type UserAuth struct {
	IdUser              *int    `json:"id"`
	Npp                 *string `json:"npp"`
	Branchalias         *int    `json:"branchalias"`
	Role                *string `json:"role"`
	Title               *string `json:"title"`
	Full_name           *string `json:"full_name"`
	Email               *string `json:"email"`
	Password            *string `json:"password"`
	Access_token        *string `json:"access_token"`
	Access_token_expire *string `json:"access_token_expire"`
	Ldap_expire         *string `json:"ldap_expire"`
	Login_failed_count  *int    `json:"login_failed_count"`
	Login_failed_expire *string `json:"login_failed_expire"`
	First_login_at      *string `json:"first_login_at"`
	Last_login_at       *string `json:"last_login_at"`
	Status              *int    `json:"status"`
}

type UserMe struct {
	IdUser              *int    `json:"id"`
	Npp                 *string `json:"npp"`
	Branchalias         *int    `json:"branchalias"`
	Role                *string `json:"role"`
	RoleId              *string `json:"role_id"`
	Title               *string `json:"title"`
	Full_name           *string `json:"full_name"`
	Email               *string `json:"email"`
	Access_token        *string `json:"access_token"`
	Access_token_expire *string `json:"access_token_expire"`
	Ldap_expire         *string `json:"ldap_expire"`
	Login_failed_count  *int    `json:"login_failed_count"`
	Login_failed_expire *string `json:"login_failed_expire"`
	First_login_at      *string `json:"first_login_at"`
	Last_login_at       *string `json:"last_login_at"`
	Status              *int    `json:"status"`
}

func (UserAuth) TableName() string {
	return "users"
}
