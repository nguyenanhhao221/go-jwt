package main

func main() {
	apiSrv := NewAPIServer("8080")
	apiSrv.Run()
}
