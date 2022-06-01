// package main

// import (
// 	"io/ioutil"
// 	"os"
// )

// func main() {
// 	os.Chdir("./downloaded-repos")

// 	file, err := os.Open("./dyte-sample-app-backend/package.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	bytes, err := ioutil.ReadAll(file)
// 	println(string(bytes))

// 	// var buffer bytes.Buffer

// 	// cmd := exec.Command("node", "--version")

// 	// cmd.Stdout = &buffer

// 	// if err := cmd.Run(); err != nil {
// 	// 	fmt.Println(err.Error())
// 	// 	return
// 	// }

// 	// fmt.Printf("%s\n", buffer.String())

// }
