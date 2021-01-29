package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

/*
 https://zhuanlan.zhihu.com/p/80213099
*/

//获取网页
func fetch(url string) string {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

//根据正则获取想要的内容
func parseUrls(url string) {
	body := fetch(url)
	body = strings.Replace(body, "\n", "", -1)
	rp := regexp.MustCompile(`<li class="clearfix">(.*?)</li>`)
	titleRe := regexp.MustCompile(`<a href="(.*?)" class="h4 link-gray-dark mb-1">(.*?)</a>`)
	items := rp.FindAllStringSubmatch(body, -1)
	for _, item := range items {
		if strings.Contains(item[1], "merge") {
			if strings.HasPrefix(titleRe.FindStringSubmatch(item[1])[2], "rgw:") {
				writefile(titleRe.FindStringSubmatch(item[1])[1], titleRe.FindStringSubmatch(item[1])[2])
			}
		}
	}
}

//追加写入文件
func writefile(url, name string) {
	var content string
	if strings.HasPrefix(url, "/") {
		url = "https://github.com" + url
	}
	if url == "##" {
		content = url + " " + name + "\r\n"
	} else {
		content = "- " + "[" + name + "]" + "(" + url + ")" + "\r\n"
	}

	file, err := os.OpenFile("rgw_monthly.md", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(content)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()

}

func parseUrls2(url string) {
	if strings.HasPrefix(url, "/") {
		url = "https://github.com" + url
	}
	body := fetch(url)
	body = strings.Replace(body, "\n", "", -1)
	rp := regexp.MustCompile(`<div class="TimelineItem-body">(.*?)</div>`)
	//rrp := regexp.MustCompile(`<div class="TimelineItem-body">(.*?)</div>`)
	//titleRe := regexp.MustCompile(`<a href="(.*?)"data-name="(.*?)"(.*?)>(.*?)</a>`)
	items := rp.FindAllStringSubmatch(body, -1)
	for _, item := range items {
		fmt.Println(item[1])
		//fmt.Println(rrp.FindStringSubmatch(item[2])[1])
		// if titleRe.FindStringSubmatch(item[1])[2] == "feature" {
		// 	fmt.Println(url + "/commits")
		// } else if titleRe.FindStringSubmatch(item[1])[2] == "bug fix" {
		// 	fmt.Println(url + "/commits")
		// } else {
		// 	fmt.Println(url + "/commits")
		// }
	}
}

func main() {
	start := time.Now()
	// for i := 0; i < 10; i++ {
	// 	parseUrls("https://movie.douban.com/top250?start=" + strconv.Itoa(25*i))
	// }
	//parseUrls2("https://github.com/ceph/ceph/pull/38861")
	writefile("##", start.Format("2006-01"))
	parseUrls("https://github.com/ceph/ceph/pulse/monthly")
	fmt.Println("success")
	elapsed := time.Since(start)
	fmt.Printf("Took %s", elapsed)
}
