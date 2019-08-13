package fileUtils

import (
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := path.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func GetAllFiles(dirPth string, ext string, files *[]string) error {
	sep := string(os.PathSeparator)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if fi.IsDir() { // 目录, 递归遍历
			if name == "res" || name == "xdoc" {
				continue
			}

			GetAllFiles(dirPth+name+sep, ext, files)
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(name, "."+ext)
			if ok {
				*files = append(*files, dirPth+name)
			}
		}
	}

	return nil
}

func GetSpecifiedFiles(scriptDir string, fileNames []string) (files []string, err error) {
	ret := make([]string, 0)

	for _, name := range fileNames {
		file := name

		if path.Ext(file) == "."+constant.SuiteExt {
			fileList := make([]string, 0)
			GetSuiteFiles(scriptDir, file, &fileList)

			for _, f := range fileList {
				ret = append(ret, f)
			}
		} else {
			ret = append(ret, file)
		}
	}

	return ret, nil
}

func GetFailedFiles(resultFile string) ([]string, string, string, error) {
	ret := make([]string, 0)
	dir := ""
	extName := ""

	content := ReadFile(resultFile)

	reg := regexp.MustCompile(`\n\sFAIL\s([^\n]+)\n`)
	arr := reg.FindAllStringSubmatch(content, -1)

	if len(arr) > 0 {
		for _, file := range arr {
			if len(file) == 1 {
				continue
			}

			caseFile := commonUtils.RemoveBlankLine(file[1])
			ret = append(ret, caseFile)

			if dir == "" {
				dir = path.Dir(caseFile)
			}
			if extName == "" {
				extName = strings.TrimLeft(path.Ext(caseFile), ".")
			}
		}
	}

	return ret, dir, extName, nil
}

func GetSuiteFiles(dirPth string, name string, fileList *[]string) {
	content := ReadFile(name)
	for _, line := range strings.Split(content, "\n") {
		file := strings.TrimSpace(line)
		if file == "" {
			return
		}

		if path.Ext(file) == "."+constant.SuiteExt {
			GetSuiteFiles(dirPth, file, fileList)
		} else {
			*fileList = append(*fileList, file)
		}
	}
}

func MkDirIfNeeded(dir string) {
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}