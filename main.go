package main

import "os"

func main() {
	a := App{}
	a.Initialazer(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		port)
	a.Run(":8010")
}
