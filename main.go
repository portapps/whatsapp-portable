//go:generate go install -v github.com/kevinburke/go-bindata/go-bindata
//go:generate go-bindata -prefix res/ -pkg assets -o assets/assets.go res/WhatsApp.lnk
//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico
package main

import (
	"io/ioutil"
	"os"
	"path"

	_ "github.com/kevinburke/go-bindata"
	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/shortcut"
	"github.com/portapps/portapps/pkg/utl"
	"github.com/portapps/whatsapp-portable/assets"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("whatsapp-portable", "WhatsApp"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	electronBinPath := utl.PathJoin(app.AppPath, utl.FindElectronAppFolder("app-", app.AppPath))

	app.Process = utl.PathJoin(app.AppPath, "WhatsApp.exe")
	app.WorkingDir = electronBinPath

	// Copy default shortcut
	shortcutPath := path.Join(utl.StartMenuPath(), "WhatsApp Portable.lnk")
	defaultShortcut, err := assets.Asset("WhatsApp.lnk")
	if err != nil {
		Log.Error().Err(err).Msg("Cannot load asset WhatsApp.lnk")
	}
	err = ioutil.WriteFile(shortcutPath, defaultShortcut, 0644)
	if err != nil {
		Log.Error().Err(err).Msg("Cannot write default shortcut")
	}

	// Update default shortcut
	err = shortcut.Create(shortcut.Shortcut{
		ShortcutPath:     shortcutPath,
		TargetPath:       app.Process,
		Arguments:        shortcut.Property{Clear: true},
		Description:      shortcut.Property{Value: "WhatsApp Portable by Portapps"},
		IconLocation:     shortcut.Property{Value: app.Process},
		WorkingDirectory: shortcut.Property{Value: app.AppPath},
	})
	if err != nil {
		Log.Error().Err(err).Msg("Cannot create shortcut")
	}
	defer func() {
		if err := os.Remove(shortcutPath); err != nil {
			Log.Error().Err(err).Msg("Cannot remove shortcut")
		}
	}()

	app.Launch(os.Args[1:])
}
