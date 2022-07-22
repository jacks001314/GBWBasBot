package fileutils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func GetFilePaths(fpath string) []string {

	paths := strings.Split(fpath,string(os.PathSeparator))
	n := len(paths)
	start := 0
	end := n

	if n>1 {

		if paths[0] == "" {

			start = 1
		}

		if paths[n-1] == "" {
			end = n-1
		}

	}

	return paths[start:end]
}

func WriteFile(fpath string,data []byte) error {
	return ioutil.WriteFile(fpath,data,0644)
}

func DeleteFile(fname string ) {

	os.Remove(fname)
}

func DeleteFiles(dir string ){

	os.RemoveAll(dir)
}

/*copy file*/
func FileCopy(dstPath string,srcPath string) (err error) {

	if err = os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return
	}

	fr,err := os.Open(srcPath)
	if err !=nil {

		return
	}

	defer fr.Close()

	fi, _ := fr.Stat()
	perm := fi.Mode()

	fw, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC,perm)
	if err != nil {
		return
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)

	return
}

func doUnzip(file *zip.File,path string) (err error) {

	// 获取到 Reader
	fr, err := file.Open()
	if err != nil {
		return
	}
	defer fr.Close()

	// 创建要写出的文件对应的 Write
	fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if err != nil {
		return
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)

	return
}

func FileIsExisted(filename string) bool {

	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}

	return existed
}


/*unzip zip file into specify dir by @dstDir*/
func UnzipFile(fpath string,dstDir string) (err error) {


	zr, err := zip.OpenReader(fpath)
	if err != nil {
		return
	}
	defer zr.Close()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dstDir != "" {
		if err = os.MkdirAll(dstDir, 0755); err != nil {
			return
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {

		path := filepath.Join(dstDir, file.Name)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return
			}
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		if err = doUnzip(file,path); err!=nil {
			return
		}

	}

	return nil

}

/*get all files in a dir*/
func GetAllFiles(dir string) (files []string) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if !info.IsDir(){

			files = append(files,path)
		}

		return nil
	})

	return
}

func dozip(writer *zip.Writer,prefix string , fname string,filterPrefix bool) (err error) {

	fr, err := os.Open(fname)
	if err != nil {
		return
	}

	defer fr.Close()

	fn := fname
	if filterPrefix {

		n := len(prefix)

		if prefix[n-1] != os.PathSeparator {

			fn = fname[n+1:]
		}else{
			fn = fname[len(prefix):]
		}
	}

	fw, err := writer.Create(fn)
	if err != nil {
		return
	}

	if _, err = io.Copy(fw,fr); err != nil {

		return
	}

	return nil
}

/*zip dir's all file into dstdir*/
func ZipFiles(zfile string,prefix string,srcFiles []string, filterPrefix bool) (err error) {

	zf, err := os.Create(zfile)
	if err != nil {
		return
	}

	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	for _,fname := range srcFiles {

		if err = dozip(zw,prefix,fname,filterPrefix) ; err !=nil {
			return
		}
	}

	return nil
}

/*read lines*/
func ReadAllLines(fpath string) (lines []string,err error){

	file, err := os.Open(fpath)
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		line = strings.TrimSpace(line)

		if line!="" {
			lines = append(lines,line)
		}
	}

	if err = scanner.Err(); err != nil {
		return
	}

	return
}

/*get all file paths starts with specify prefix
*
*/
func GetFilesStartsWith(prefix string) (files []string) {

	if strings.HasSuffix(prefix,"*") {

		name := filepath.Base(prefix)
		namePrefix:= name[0:len(name)-1]
		pp := filepath.Dir(prefix)


		filepath.Walk(pp, func(path string, info os.FileInfo, err error) error {

			if !info.IsDir() && strings.HasPrefix(info.Name(),namePrefix) {

				files = append(files,path)
			}

			return nil
		})

	}else {

		files = append(files,prefix)
	}

	return
}

/*generate file from template file */
func GenerateFileFromTemplateFile(fpath string,tpath string,tempData interface{}) error{

	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)

	if err != nil {

		errS := fmt.Sprintf("Cannot open file:%s to store content from template file:%s,err:%v", fpath, tpath, err)
		return fmt.Errorf("%s", errS)
	}

	defer file.Close()

	t, err := template.ParseFiles(tpath)

	if err != nil {

		errS := fmt.Sprintf("Cannot parse template file:%s,err:%v", tpath, err)
		return fmt.Errorf("%s", errS)
	}

	return  t.Execute(file,tempData)
}




