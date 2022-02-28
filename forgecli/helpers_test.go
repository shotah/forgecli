package forgecli

import (
	"testing"
)

func TestOSTargetDirectory(t *testing.T) {
	var app appEnv
	app.osTargetDirectory()
	expected := "mods"
	if app.destination[len(app.destination)-4:] != expected {
		t.Errorf("Test failed, expected: '%s', got:  '%s'", expected, app.destination)
	}
}

// func TestEnsureDestination(t *testing.T) {
// 	var app appEnv
// 	app.destination = "folder/test"
// 	app.ensureDestination()
// 	if _, err := os.Stat(app.destination); os.IsNotExist(err) {
// 		t.Errorf("Test failed, expected to create: '%s'", app.destination)
// 	}
// 	rootFolder := strings.Split(app.destination, "/")[0]
// 	err := os.RemoveAll(rootFolder)
// 	check(err)
// }
