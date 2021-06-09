package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"new-pac-subscribe/src/model"
	"new-pac-subscribe/src/utils"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func uploadVmess(picBed string, path string, key string) (url string, err error) {
	if picBed == "qiniu" {
		prefix := model.QiniuGetResourcePrefix()
		joinKey := fmt.Sprintf("%s/%s/%s", prefix, "vmess", key)
		// 如果 prefix 为空需删除开头的 /
		joinKey = strings.TrimPrefix(joinKey, "/")
		retQiniu, err := model.QiniuUpload(path, joinKey)
		if err != nil {
			return "", err
		}
		url = retQiniu.URL
		os.Remove(path)
	}
	return url, nil
}

func fetchVmess() string {
	fmt.Printf("fetch begin\n")
	res, err := http.Get("https://github.com/Alvin9999/new-pac/wiki/v2ray%E5%85%8D%E8%B4%B9%E8%B4%A6%E5%8F%B7")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	vmess := ""

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		vmessReg := regexp.MustCompile(`vmess://(.*?)$`)
		if vmessReg.MatchString(s.Text()) {
			vmess += s.Text() + "\n"
		}
	})
	fmt.Printf("%s\n", vmess)
	fmt.Printf("fetch end \n")

	return vmess
}

func main() {
	dir, _ := os.Getwd()
	vmessDir := path.Join(dir, "./raw/vmess") + "/"
	vmess := []byte(fetchVmess())
	fileName := "m5ldCI6ICJ3cyIsDQogICJ0eXBlI"
	encoded := []byte(base64.StdEncoding.EncodeToString(vmess))
	utils.WriteJSONFile(vmessDir, fileName, encoded)
	uploadVmess("qiniu", vmessDir+fileName, fileName)
}
