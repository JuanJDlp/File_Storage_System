package main

import (
	"fmt"

	"github.com/JuanJDlp/File_Storage_System/internal"
)

func main() {
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello world")
	// })
	// e.Logger.Fatal(e.Start(":8080"))
	fmt.Print(internal.HashString("Final Exam.pdf"))
}
