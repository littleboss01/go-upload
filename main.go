package main
//package main
import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


const configDir = "/etc/lya/config"
const exceptionFile = "asf.json"


func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// 检查请求方法是否为 POST
		if r.Method == http.MethodPost {
			// 获取文件
			file, header, err := r.FormFile("file")
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			defer file.Close()
			// 检查文件大小
			if header.Size > (2 << 20) { // 2MB
				fmt.Fprintln(w, "File size exceeded maximum limit of 2MB.")
				return
			}
			// 创建目标文件
			targetPath := filepath.Join(configDir, header.Filename)
			targetFile, err := os.Create(targetPath)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			defer targetFile.Close()
			// 将文件复制到目标位置
			if _, err := io.Copy(targetFile, file); err != nil {
				fmt.Fprintln(w, err)
				return
			}
			fmt.Fprintln(w, "File uploaded successfully.")
		} else {
			fmt.Fprintln(w, "Invalid request method. Must be POST.")
		}
	})
	http.HandleFunc("/cleanup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			files, err := ioutil.ReadDir(configDir)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			for _, file := range files {
				if strings.ToLower(file.Name())  == exceptionFile {
					continue
				}else if strings.ToLower(file.Name()) == "asf.db" {
					continue
				}
				os.Remove(filepath.Join(configDir, file.Name()))
			}
			fmt.Fprintln(w, "Cleaned up config directory successfully.")
		} else {
			fmt.Fprintln(w, "Invalid request method. Must be POST.")
		}
	})
	http.HandleFunc("/checkstartup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// 使用操作系统特定的方式检测开机启动项
			// 例如在 Ubuntu 中使用 update-rc.d
			startupCheck, err := exec.Command("update-rc.d", "lya", "status").Output()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			fmt.Fprintln(w, string(startupCheck))
		} else {
			fmt.Fprintln(w, "Invalid request method. Must be GET.")
		}
	})
	http.HandleFunc("/addstartup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 使用操作系统特定的方式添加开机启动项
			// 例如在 Ubuntu 中使用 update-rc.d
			_, err := exec.Command("update-rc.d", "lya", "defaults").Output()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			fmt.Fprintln(w, "Added lya to startup successfully.")
		} else {
			fmt.Fprintln(w, "Invalid request method. Must be POST.")
		}
	})
	http.ListenAndServe(":8080", nil)


}

