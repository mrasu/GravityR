/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"embed"
	"github.com/mrasu/GravityR/cmd"
	"github.com/mrasu/GravityR/injection"
	"github.com/mrasu/GravityR/lib"
	"net/http"
)

//go:embed client/dist/assets/*
var clientDist embed.FS

func main() {
	injection.ClientDist = clientDist
	http.DefaultTransport = lib.NewHttpTransport()

	cmd.Execute()
}
