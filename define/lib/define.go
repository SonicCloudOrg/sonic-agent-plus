package lib

import _ "embed"

var (
	// from: https://github.com/aoliaoaoaojiao/AndroidTouch
	//go:embed jar/AndroidTouch.jar
	AndroidTouchJarBytes []byte
)

const (
	RemoteTouchToolPath = "/data/local/tmp/AndroidTouch.jar"
)
