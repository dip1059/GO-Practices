package Admin

import (
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	Product = make(map[uint]Mod.Product)
)

func AddProductGet(c *gin.Context) {
	if user, success := M.IsAuthAdminUser(c, G.FStore); success {
		var karats []Mod.Karat
		karats = R.Karats(karats, "status=?", 1)
		
		var products []Mod.Product
		products = R.ProductsWithOthers(products, 0)
		for _, product := range products {
			Product[product.ID] = product
		}
		var position = make([]int, 0)
		lenPosition := len(Product) + 1
		for i := 1; i <= lenPosition; i++ {
			position = append(position, i)
		}

		c.HTML(http.StatusOK, "add-product.html", map[string]interface{}{"Position": position,
			"AppEnv": G.AppEnv, "User": user, "Msg": G.Msg, "Nav":"products", "Title": "Add-Product", "Karats":karats})
		G.Msg.Success = ""
		G.Msg.Fail = ""
	}
}

func AddProductPost(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var product Mod.Product
	err := c.ShouldBind(&product)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill up all the input fields with valid values."
		c.Redirect(http.StatusFound, "/add-product")
		return
	}
	if product.Type == 2 {
		min, err := strconv.ParseFloat(c.PostForm("min"), 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid minimum amount."
			c.Redirect(http.StatusFound, "/add-product")
			return
		}
		if min < 0.1 {
			G.Msg.Fail = "Minimum amount can't be less than 0.1 gram."
			c.Redirect(http.StatusFound, "/add-product")
			return
		}
	}

	position := len(Product) + 1
	var pro2 Mod.Product
	pro2 = R.Product(pro2,0, "position=?", product.Position)
	if pro2.ID != 0 {
		pro2.Position = position
		if !R.UpdateProduct(pro2) {
			G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
			c.Redirect(http.StatusFound, "/add-product")
		}
	}

	product.Description.String = c.PostForm("description")
	product.Description = H.NullStringProcess(product.Description)
	/*product.Size.String = c.PostForm("size")
	product.Size = H.NullStringProcess(product.Size)
	product.Color.String = c.PostForm("color")
	product.Color = H.NullStringProcess(product.Color)*/
	img, _ := c.FormFile("img")
	ext := filepath.Ext(img.Filename)
	imgName := H.RandomString(60) + ext
	dst := "./Storage/Images/" + imgName
	err = c.SaveUploadedFile(img, dst)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Image Upload Failed. Try Again Later."
		c.Redirect(http.StatusFound, "/add-product")
		return
	}
	imgUrl := []byte(dst)
	product.ImgUrl.String = string(imgUrl[1:])
	product.ImgUrl = H.NullStringProcess(product.ImgUrl)

	if R.AddProduct(product) {
		G.Msg.Success = "Product Added Successfully."
		c.Redirect(http.StatusFound, "/all-product")
		return
	} else {
		G.Msg.Fail = "Some Error Occured. Please Try Again Later."
		c.Redirect(http.StatusFound, "/add-product")
		return
	}
}

func AllProduct(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var products []Mod.Product
	products = R.ProductsWithOthers(products, 0)
	for _, product := range products {
		Product[product.ID] = product
	}
	c.HTML(http.StatusOK, "all-product.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user, "Nav":"products", "Title": "All-Product", "Products": products, "Msg": G.Msg})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func MakeProductInactive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var product Mod.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product = Product[uint(id)]
	product.Status = 0
	if R.UpdateProduct(product) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/all-product")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/all-product")
	}
}

func MakeProductActive(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var product Mod.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product = Product[uint(id)]
	product.Status = 1
	if R.UpdateProduct(product) {
		G.Msg.Success = "Status Updated Successfully"
		c.Redirect(http.StatusFound, "/all-product")
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/all-product")
	}
}

func EditProduct(c *gin.Context) {
	var user Mod.User
	var success bool
	user, success = M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var product Mod.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product = Product[uint(id)]

	var position = make([]int, 0)
	lenPosition := len(Product)
	for i := 1; i <= lenPosition; i++ {
		position = append(position, i)
	}

	var karats []Mod.Karat
	karats = R.Karats(karats, "status=?", 1)

	c.HTML(http.StatusOK, "edit-product.html", map[string]interface{}{"Position": position,
		"AppEnv": G.AppEnv, "User": user, "Nav":"products", "Title": "Edit-Product", "Product": product, "Msg": G.Msg, "Karats":karats})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}

func UpdateProduct(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var product Mod.Product
	id, _ := strconv.Atoi(c.PostForm("id"))
	product = Product[uint(id)]

	prevPosition := product.Position

	err := c.ShouldBind(&product)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill up all the input fields with valid values."
		c.Redirect(http.StatusFound, "/all-product")
		return
	}

	if product.Type == 2 {
		min, err := strconv.ParseFloat(c.PostForm("min"), 64)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Invalid minimum amount."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
		if min < 0.1 {
			G.Msg.Fail = "Minimum amount can't be less than 0.1 gram."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
	}

	if prevPosition != product.Position {
		var pro2 Mod.Product
		pro2 = R.Product(pro2, 0,"position=?", product.Position)
		if pro2.ID != 0 {
			pro2.Position = prevPosition
			if !R.UpdateProduct(pro2) {
				G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
				c.Redirect(http.StatusFound, "/all-product")
			}
		}
	}

	product.Description.String = c.PostForm("description")
	product.Description = H.NullStringProcess(product.Description)
	/*product.Size.String = c.PostForm("size")
	product.Size = H.NullStringProcess(product.Size)
	product.Color.String = c.PostForm("color")
	product.Color = H.NullStringProcess(product.Color)*/
	img, _ := c.FormFile("img")
	if img != nil {
		os.Remove("." + product.ImgUrl.String)

		ext := filepath.Ext(img.Filename)
		imgName := H.RandomString(60) + ext
		dst := "./Storage/Images/" + imgName
		err = c.SaveUploadedFile(img, dst)
		if err != nil {
			log.Println(err.Error())
			G.Msg.Fail = "Image Upload Failed. Try Again Later."
			c.Redirect(http.StatusFound, "/all-product")
			return
		}
		imgUrl := []byte(dst)
		product.ImgUrl.String = string(imgUrl[1:])
		product.ImgUrl = H.NullStringProcess(product.ImgUrl)
	}
	if R.UpdateProduct(product) {
		G.Msg.Success = "Updated Successfully"
		c.Redirect(http.StatusFound, "/all-product")
	} else {
		G.Msg.Fail = "Some Error Occurred, Update Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/all-product")
	}
}

func DeleteProduct(c *gin.Context) {
	if _, success := M.IsAuthAdminUser(c, G.FStore); !success {
		return
	}
	var product Mod.Product
	id, _ := strconv.Atoi(c.Param("id"))
	product = Product[uint(id)]
	if R.DeleteProduct(product) {
		G.Msg.Success = "Deleted Successfully"
		c.Redirect(http.StatusFound, "/all-product")
	} else {
		G.Msg.Fail = "Some Error Occurred, Deletion Failed. Please Try Again Later."
		c.Redirect(http.StatusFound, "/all-product")
	}
}
