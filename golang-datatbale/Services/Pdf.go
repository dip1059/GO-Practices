package Services

import (
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool) {
	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	// write whole the body
	err1 := ioutil.WriteFile("Storage/Temp/"+fileName+".html", []byte(r.body), 0777)
	if err1 != nil {
		log.Println(err1.Error())
		//RemoveTempHTML(fileName)
		return false
	}

	f, err := os.Open("Storage/Temp/"+fileName+".html")
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Println(err.Error())
		//RemoveTempHTML(fileName)
		return false
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Println(err.Error())
		//RemoveTempHTML(fileName)
		return false
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Println(err.Error())
		//RemoveTempHTML(fileName)
		return false
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Println(err.Error())
		//RemoveTempHTML(fileName)
		return false
	}
	//RemoveTempHTML(fileName)
	return true
}

func RemoveTempHTML(fileName string) {
	err := os.Remove("Storage/Temp/"+fileName+".html")
	if err != nil {
		log.Println(err.Error())
		//RemoveTempHTML(fileName)
	}
}
