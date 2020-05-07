/*
Copyright © 2020 si9ma <si9ma@si9ma.com>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
)

var port string

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "a fis receiver with golang",
	Long:  `a fis receiver build with golang`,
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", uploadHandler)

		hostname, _ := os.Hostname()
		fmt.Println("start receiver at http://" + hostname + ":" + port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal("create receiver fail : ", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)

	receiverCmd.Flags().StringVarP(&port, "port", "p", "8080", "listen port")
}

// reference from https://github.com/lrenc/fis-receiver/blob/master/main.go
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		io.WriteString(w, "<p>Hey guys, I'm ready for that, you know.</p>")
		return
	}
	if r.Method == "POST" {
		f, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		defer f.Close()

		to := r.FormValue("to")
		_, filePath := getFileInfo(to)
		_, err = os.Open(filePath)

		// create directory if directory not exist
		if err != nil && os.IsNotExist(err) {
			oldMask := syscall.Umask(0)
			_ = os.MkdirAll(filePath, os.ModePerm)
			syscall.Umask(oldMask)
		}

		t, err := os.Create(to)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
		fmt.Println(r.FormValue("to"))

		// success
		io.WriteString(w, "0")
	}
}

// get file name and directory of file
func getFileInfo(file string) (fileName, filePath string) {
	index := strings.LastIndex(file, "/")
	return file[index+1:], file[0 : index+1]
}
