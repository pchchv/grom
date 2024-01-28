package grom

import "os"

func routerSetupBody() string {
	fileBytes, _ := os.ReadFile(testFilename())
	return string(fileBytes)
}

func testFilename() string {
	return "router_setup.go"
}
