package Repositories

import (
	Cfg "gold-store/Config"
	M "gold-store/Models"
)

type CashData struct {
	Rate float64
	CashPoint float64
}


func CheckCashPoint(user M.User) CashData {
	db2 := Cfg.DBConnect()
	var data CashData
	db2.Raw("Select value `rate` from admin_settings where slug='euro_per_cash_point';").Scan(&data)
	db2.Raw("Select points `cash_point` from current_cash_points where user_id=?;", user.ID).Scan(&data)
	defer db2.Close()
	return data
}

func DeductCashPoint(user M.User, cashP float64) bool {
	db2 := Cfg.DBConnect()
	if db2.Exec("Update current_cash_points set points = points - ? where user_id=?;", cashP, user.ID).RowsAffected == 1 {
		defer db2.Close()
		return true
	} else {
		defer db2.Close()
		return false
	}
}