package Repositories

import (
	Cfg "gold-store/Config"
	M "gold-store/Models"
	"log"
)

func ReadUserWithID(user M.User) (M.User, bool) {
	db := Cfg.DBConnect()
	notFound := db.First(&user).RecordNotFound()

	if notFound {
		defer db.Close()
		return user, false
	} else {

		db.First(&user.Role, "id=?", user.RoleID)

		defer db.Close()
		return user, true
	}
}

func ReadUserWithEmail(user M.User, where ...interface{}) (M.User, bool) {
	db := Cfg.DBConnect()
	var notFound bool
	if where != nil {
		notFound = db.First(&user, "email=? and "+where[0].(string), user.Email).RecordNotFound()
	} else {
		notFound = db.First(&user, "email=?", user.Email).RecordNotFound()
	}

	if notFound {
		defer db.Close()
		return user, false
	} else {
		db.First(&user.Wallet, "user_id=?", user.ID)
		db.First(&user.DefaultAddress, "user_id=? and status=?", user.ID, 1)
		defer db.Close()
		return user, true
	}

}

func AddSocialUser(user M.User) (M.User, bool) {
	db := Cfg.DBConnect()
	db = db.Begin()
	user2 := user
	notFound := db.First(&user, "email=?", user.Email).RecordNotFound()
	if notFound {
		err := db.Create(&user2).Error
		if err != nil {
			log.Println(err.Error())
			db.Rollback()
			defer db.Close()
			return user, false
		}

		user2.Wallet.UserID = user2.ID
		err = db.Create(&user2.Wallet).Error
		if err != nil {
			db.Rollback()
			defer db.Close()
			return user, false
		}

		user = user2
	} else if user.RegType > 1 {
		user.SocialAuthID = user2.SocialAuthID
		err := db.Save(&user).Error
		if err != nil {
			db.Rollback()
			defer db.Close()
			return user, false
		}
	} else {
		db.Rollback()
		defer db.Close()
		return user, false
	}
	db.Commit()
	defer db.Close()
	return user, true
}

func ReadUserWithPhone(user M.User, where ...interface{}) (M.User, bool) {
	db := Cfg.DBConnect()
	var notFound bool
	if where != nil {
		notFound = db.First(&user, "phone=? and "+where[0].(string), user.Phone).RecordNotFound()
	} else {
		notFound = db.First(&user, "phone=?", user.Phone).RecordNotFound()
	}

	if notFound {
		defer db.Close()
		return user, false
	} else {
		db.First(&user.Wallet, "user_id=?", user.ID)
		db.First(&user.DefaultAddress, "user_id=? and status=?", user.ID, 1)
		defer db.Close()
		return user, true
	}
}

func Register(user M.User) (M.User, bool) {
	db := Cfg.DBConnect()
	db = db.Begin()
	err := db.Create(&user).Error
	if err != nil {
		db.Rollback()
		defer db.Close()
		return user, false
	}

	user.Wallet.UserID = user.ID
	err = db.Create(&user.Wallet).Error
	if err != nil {
		db.Rollback()
		defer db.Close()
		return user, false
	}

	db.Commit()
	defer db.Close()
	return user, true
}

func Login(user M.User) (M.User, bool) {
	var success bool
	user, success = ReadUserWithEmail(user)
	if success {
		return user, true
	} else {
		return user, false
	}
}

/*func SetRememberToken(user M.User) bool {
	db := Cfg.DBConnect()

	if db.Model(&user).Where("email=?", user.Email).Updates(
		map[string]interface{}{"remember_token":user.RememberToken.String}).RowsAffected == 0 {
		defer db.Close()
		return false
	}

	defer db.Close()
	return true
}*/

func ActivateAccount(user M.User, query string, args ...interface{}) (M.User, bool) {
	db := Cfg.DBConnect()

	if db.Model(&user).Where(query, args...).Updates(
		map[string]interface{}{"active_status": 1, "email_verification": "", "email_verify_code": ""}).RowsAffected == 0 {
		defer db.Close()
		return user, false
	}

	db.First(&user, "email=?", user.Email)

	defer db.Close()
	return user, true
}

func SendPasswordResetLink(ps M.PasswordReset) bool {
	db := Cfg.DBConnect()
	err := db.Create(&ps).Error
	if err != nil {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func ResetPasswordGet(ps M.PasswordReset) bool {
	db := Cfg.DBConnect()

	notFound := db.First(&ps, "email=? and token=? and status=0", ps.Email, ps.Token).RecordNotFound()
	if notFound {
		defer db.Close()
		return false
	}

	defer db.Close()
	return true
}

func IsResetCodeValid(ps M.PasswordReset) bool {
	db := Cfg.DBConnect()
	if db.First(&ps, "email=? and code=? and status=?", ps.Email, ps.Code, 0).RecordNotFound() {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func IsResetTokenValid(ps M.PasswordReset) bool {
	db := Cfg.DBConnect()
	if db.First(&ps, "email=? and token=? and status=? and re", ps.Email, ps.Token, 0).RecordNotFound() {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func ResetPasswordPost(user M.User, ps M.PasswordReset, query string, args ...interface{}) bool {
	db := Cfg.DBConnect()
	tx := db.Begin()
	err := tx.Model(&user).Where("email=? and reg_type=1", user.Email).Updates(
		map[string]interface{}{"password": user.Password, "active_status": 1, "email_verification": "", "email_verify_code": ""}).Error
	if err != nil {
		tx.Rollback()
		defer db.Close()
		return false
	}

	err = tx.Model(&ps).Where(query, args...).Updates(
		map[string]interface{}{"status": 1}).Error
	if err != nil {
		//log.Println("AuthRepo.go Log8", err.Error())
		tx.Rollback()
		defer db.Close()
		return false
	}

	err = tx.Model(&ps).Where("email=?", ps.Email).Updates(
		map[string]interface{}{"token": "", "code": ""}).Error
	if err != nil {
		//log.Println("AuthRepo.go Log8", err.Error())
		tx.Rollback()
		defer db.Close()
		return false
	}

	err = tx.Commit().Error
	if err != nil {
		//log.Println("AuthRepo.go Log8", err.Error())
		tx.Rollback()
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

/*func Logout(user M.User) {
	db := Cfg.DBConnect()
	user.RememberToken.Valid = false
	db.Model(&user).Where("email=?", user.Email).Update(
		"remember_token", user.RememberToken)
	defer db.Close()

}*/
