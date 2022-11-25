package thelper

import (
	"github.com/mrasu/GravityR/injection"
	"testing/fstest"
)

func InjectClientDist() {
	injection.ClientDist = fstest.MapFS{
		"client/dist/assets/main.js": {
			Data: []byte("console.log('hello')"),
		},
		"client/dist/assets/main.css": {
			Data: []byte("body{}"),
		},
		"client/dist/assets/mermaid.js": {
			Data: []byte("console.log('hello')"),
		},
		"client/dist/assets/mermaid.css": {
			Data: []byte("body{}"),
		},
	}
}
