package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var homeContents = make([]Mod.HomeContent,0)

func HomeContentSeeder() {
	db := Cfg.DBConnect()

	homeContent1()
	homeContent2()
	homeContent3()
	for i,_ := range homeContents {
		db.FirstOrCreate(&homeContents[i])
	}
	defer db.Close()
}

func homeContent1() {
	var homeContent = Mod.HomeContent{
		TextContent: `<p class="gold-banner-subtitle">Consectetur adipisicing elit, sed do eiusmod</p>
                    <h1 class="gold-banner-title">Gold Lorem Ipsum</h1>
                    <p class="gold-banner-text">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>`,
		Image: "/Public/User/images/bannger-gold.png",
	}
	homeContent.ID = 1
	homeContents = append(homeContents, homeContent)
}

func homeContent2() {
	var homeContent = Mod.HomeContent{
		TextContent: `<h2>Lorem Ipsum <span>Dolorium</span></h2>
                        <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Quis ipsum suspendisse ultrices gravida. Risus commodo viverra maecenas accumsan lacus vel facilisis.</p>
                        <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
                        <p>Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.</p>`,
		Image: "/Public/User/images/about-banner.jpg",
	}
	homeContent.ID = 2
	homeContents = append(homeContents, homeContent)
}

func homeContent3() {
	var homeContent = Mod.HomeContent{
		TextContent: `<p class="gold-banner-subtitle">Consectetur adipisicing elit, sed do eiusmod</p>
                    <h3>Lorem Ipsum Dolor Sit Amet Delirium</h3>
                    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>`,
		Image: "/Public/User/images/add-banner.jpg",
	}
	homeContent.ID = 3
	homeContents = append(homeContents, homeContent)
}
