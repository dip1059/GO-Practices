package Config

import (
	G "gold-store/Globals"
	Mod "gold-store/Models"
)

func LoadAdminSettings() {
	db := DBConnect()
	var admSets []Mod.AdminSetting
	db.Find(&admSets)
	for i, _ := range admSets {
		G.Adm[admSets[i].Slug] = admSets[i]
	}
	defer db.Close()
}
