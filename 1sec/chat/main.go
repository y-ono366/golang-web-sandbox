package main

import(
	"log"
	"net/http"
)


func main() {
	http.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte (
		`<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title></title>
		</head>
		<body>
			チャットしようぜ!!
		</body>
		</html>
		`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
