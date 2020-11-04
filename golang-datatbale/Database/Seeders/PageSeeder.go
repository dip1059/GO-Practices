package Seeders

import (
	Cfg "gold-store/Config"
	H "gold-store/Helpers"
	Mod "gold-store/Models"
)

var pages = make([]Mod.Page,0)

func PageSeeder() {
	db := Cfg.DBConnect()

	page1()
	page2()
	for i,_ := range pages {
		db.FirstOrCreate(&pages[i],&Mod.Page{
			Url:pages[i].Url},
		)
	}
	defer db.Close()
}

func page1() {
	var page = Mod.Page{
		Title: "About Us",
		Url: H.MakeUrl("About Us"),
		TextContent: `<div class="gold-about-section">
						<div class="about-banner-img">
							<img src="/assets/Public/User/images/about1.jpg" class="img-fluid" alt="">
						</div>
						<div class="container">
							<div class="row">
								<div class="col-lg-4">
									<h2>
										About us
										<span class="subtitle">trust gold 999</span>
									</h2>
									<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut
										labore et dolore magna aliqua.</p>
									<p class="text-gray">Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
										pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt.
									</p>
								</div>
							</div>
						</div>
					</div>
					
					
					<!-- =============================================
								(02) End About Banner section 
					============================================== -->
					
					
					<!-- =============================================
								(03) Statr About one section 
					============================================== -->
					
					<div class="gold-about-one" style="background: #202022 !important;">
						<div class="gold-about-one-img">
							<img src="/assets/Public/User/images/about2.jpg" class="img-fluid" alt="">
						</div>
						<div class="container">
							<div class="row">
								<div class="col-lg-5 offset-lg-6">
									<div class="gold-about-one-content">
										<h3>Lorem Delirium Versus</h3>
										<p>
											Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
											labore et dolore magna aliqua. Quis ipsum suspendisse ultrices gravida. Risus commodo
											viverra maecenas accumsan lacus vel facilisis.
										</p>
										<p>
											Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt
											ut labore et dolore magna aliqua.
										</p>
										<p>
											Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
											commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
											dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.
										</p>
									</div>
								</div>
							</div>
						</div>
					</div>
					
					<!-- =============================================
								(03) End About one section 
					============================================== -->
					
					
					<!-- =============================================
								(04) Statr About two section 
					============================================== -->
					
					<div class="gold-about-two" style="background: #202022 !important;">
						<div class="container">
							<div class="row align-items-center">
								<div class="col-lg-5 offset-lg-1">
									<div class="gold-about-two-content">
										<h3>
											<span>Lorem Delirium</span>
											Versus Insteroym Elit
											<span>Sed do Labore</span>
										</h3>
										<p>
											Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
											labore et dolore magna aliqua. Quis ipsum suspendisse ultrices gravida. Risus commodo
											viverra maecenas accumsan lacus vel facilisis.
										</p>
										<p>
											Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt
											ut labore et dolore magna aliqua.
										</p>
										<p>
											Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
											commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum
											dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident.
										</p>
									</div>
								</div>
								<div class="col-lg-6">
									<div class="gold-about-two-img">
										<img src="/assets/Public/User/images/about3.png" class="img-fluid" alt="">
									</div>
								</div>
							</div>
						</div>
					</div>
					
					<!-- =============================================
								(04) End About two section 
					============================================== -->
				
				
					
					<!-- =============================================
							(05) End About banner two section 
					============================================== -->
				
					<div class="gold-about-banner-two">
						<div class="container">
							<div class="row">
								<div class="col-lg-7">
									<h2>Lorem Ipsum Dolor Sit Amet Delirium</h2>
									<span class="subtitle">Consectetur adipisicing elit, sed do eiusmod</span>
									<p>
										Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut
										labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco
										laboris nisi ut aliquip ex ea commodo consequat.
									</p>
									<p>
										Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
										pariatur. Excepteur sint occaecat cupidatat non proident.
									</p>
								</div>
							</div>
						</div>
					</div>`,
		Status: 1,
	}
	pages = append(pages, page)
}

func page2() {
	var page = Mod.Page{
		Title: "Contact Us",
		Url: H.MakeUrl("Contact Us"),
					TextContent: `<h1>
									Contact Us
									<span>trust gold 999</span>
								</h1>
								<p>
									Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut
									labore et dolore magna aliqua.
								</p>
				
								<div class="gold-contact-address">
									<ul>
										<li class="address">451 Wall Street, Lisbon, Portugal</li>
										<li class="email">contact@trustgold999.com</li>
										<li class="phone">
											<span>+12 345 678 9000</span>
											<span>+12 345 678 9000</span>
										</li>
									</ul>
								</div>`,
		Status: 1,
	}
	pages = append(pages, page)
}