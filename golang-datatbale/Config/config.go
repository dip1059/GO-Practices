package Config

func Config(){
	AppConfig()
	SocialAuthConfig()
	DirectoryConfig()
	LoadAdminSettings()
	LoadCountryMap()
	WriteLogFile()
	WritePanicLogFile()
	go func() {
		DailyLogFile()
	}()
}
