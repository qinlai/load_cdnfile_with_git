package main

import (
    "crypto/sha1"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strings"
)

const (
	cdnUrl = "http://xxx.com/"
	baseTag = "2015113003"
	newTag = "2015120301"
)

func main() {
    gitIndex := cdnUrl + "index/" + baseTag
    gitDiff := cdnUrl + "diff/" + baseTag + ".." + newTag

    index := getGitData(gitIndex)
    diff := getGitData(gitDiff)

    f1, f2, f3 := "Iphone/version.txt", "Iphone/Version/1.zip", "Iphone/Version/66_70.zip"

    d1 := loadData(cdnUrl, index, diff, f1)
    d2 := loadData(cdnUrl, index, diff, f2)
    d3 := loadData(cdnUrl, index, diff, f3)

    fmt.Println(f1, len(d1))
    fmt.Println(f2, len(d2))
    fmt.Println(f3, len(d3))
}

func loadData(cndUrl string, index, diff *map[string]string, file string) []byte {
    f := fmt.Sprintf("%x", sha1.Sum([]byte(file)))[0:10]

    dirName, fileName := getFileInfo(file)

    dir, ok := (*diff)[f]
    if ok {
        return getHttp(cndUrl + "file/" + dir + "/" + fileName)
    } else {
        f = fmt.Sprintf("%x", sha1.Sum([]byte(dirName)))[0:10]
        dir, ok = (*index)[f]
        if ok {
            return getHttp(cndUrl + "tree/" + dir + "/" + fileName)
        }
    }

    return nil
}

func getGitData(uri string) *map[string]string {
    return formatGitData(changeHexToString(getHttp(uri)))
}

func changeHexToString(data []byte) string {
    return fmt.Sprintf("%x", data)
}

func getHttp(uri string) []byte {
    fmt.Println("http get :", uri)
    u, _ := url.Parse(uri)
    res, _ := http.Get(u.String())
    data, _ := ioutil.ReadAll(res.Body)

    return data
}

func formatGitData(s string) *map[string]string {
    n := 50
    result := make(map[string]string)
    i := 0
    l := len(s)
    start := 0

    for {
        start = i * n
        result[s[start:start+10]] = s[start+10 : start+n]
        i++

        if i*n >= l {
            break
        }
    }

    return &result
}

func getFileInfo(file string) (dir, fileName string) {
    arr := strings.Split(file, "/")
    l := len(arr)
    return strings.Join(arr[0:l-1], "/"), strings.Join(arr[l-1:l], "")
}
