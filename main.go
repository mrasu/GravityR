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

var (
	version = "0.0.0"
	commit  = "none"
)

func main() {
	injection.ClientDist = clientDist
	injection.SetBinInfo(version, commit)

	http.DefaultTransport = lib.NewHttpTransport()

	cmd.Execute()
}
