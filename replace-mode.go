package main

type ReplaceMode byte

const (
	NormalReplace ReplaceMode = iota
	NoReplace
	PropertiesReplace
)

// fileReplaceModes contains the list of files with modified ReplaceMode values
var fileReplaceModes = map[string]ReplaceMode{
	".github/FUNDING.yml":     NoReplace,
	"forge/gradle.properties": NoReplace,

	"gradle/wrapper/gradle-wrapper.jar":        NoReplace,
	"gradle/wrapper/gradle-wrapper.properties": NoReplace,

	".gitattributes":    NoReplace,
	".gitignore":        NoReplace,
	"gradle.properties": PropertiesReplace,
	"gradlew":           NoReplace,
	"gradlew.bat":       NoReplace,
	"LICENSE.md":        NoReplace,
}
