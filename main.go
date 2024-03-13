package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/fatih/color"
	"github.com/wessie/appdirs"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	userConfigDir = appdirs.UserConfigDir("mcmod-template", "mrmelon54", "", false)
	configPath    = filepath.Join(userConfigDir, "config.json")
	questionColor = color.New(color.FgCyan)
	toModId       = regexp.MustCompile("[^a-zA-Z0-9]+")

	//go:embed all:template
	templateDir embed.FS
)

func MustSub(f fs.FS, dir string) fs.FS {
	s, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return s
}

func prompt(s string) string {
	_, _ = questionColor.Print(s)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func promptCheckbox(s string) bool {
	_, _ = questionColor.Print(s)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	switch scanner.Text() {
	case "y", "Y":
		return true
	default:
		return false
	}
}

func fakePrompt(s string, v string) {
	_, _ = questionColor.Print(s)
	fmt.Println(v)
}

type Config struct {
	ModGroupBase  string `json:"mod_group_base"`
	ModSiteBase   string `json:"mod_site_base"`
	ModSourceBase string `json:"mod_source_base"`
}

func main() {
	conf := Config{
		ModGroupBase:  "com.example",
		ModSiteBase:   "https://example.com/minecraft",
		ModSourceBase: "https://github.com/example",
	}
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	confOpen, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Config file is missing, creating default mcmConfig...")
			create, err := os.Create(configPath)
			if err != nil {
				log.Fatal(err)
			}
			enc := json.NewEncoder(create)
			enc.SetIndent("", "  ")
			err = enc.Encode(conf)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		dec := json.NewDecoder(confOpen)
		err := dec.Decode(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	var templateLayers fs.FS
	templateLayers = MustSub(templateDir, "template")

	modName := prompt("[?] Mod Name: ")
	modDesc := prompt("[?] Mod Description: ")

	modNameSafe := toModId.ReplaceAllString(modName, "_")
	modId := strings.ToLower(modNameSafe)
	modClass := strings.ReplaceAll(modNameSafe, "_", "")
	modGroup := conf.ModGroupBase + "." + modClass
	modDash := strings.ReplaceAll(modId, "_", "-")
	modWebsite := conf.ModSiteBase + "/" + modDash
	modSource := conf.ModSourceBase + "/" + modId
	modIssue := modSource + "/issues"

	fakePrompt("[+] Mod ID: ", modId)
	fakePrompt("[+] Mod Class: ", modClass)
	fakePrompt("[+] Mod Group: ", modGroup)
	fakePrompt("[+] Mod Website: ", modWebsite)
	fakePrompt("[+] Mod Source: ", modSource)
	fakePrompt("[+] Mod Issue: ", modIssue)

	modInfo := make(ModInfo)
	modInfo["modname"] = modName
	modInfo["moddesc"] = modDesc
	modInfo["modid"] = modId
	modInfo["modclass"] = modClass
	modInfo["modgroup"] = modGroup
	modInfo["moddash"] = modDash
	modInfo["modwebsite"] = modWebsite
	modInfo["modsource"] = modSource
	modInfo["modissue"] = modIssue

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wdPath := filepath.Join(cwd, modId)

	fakePrompt("[@] Mod Path: ", wdPath)

	if !promptCheckbox("[?] Is that ok [y/N]? ") {
		log.Println("Goodbye")
		os.Exit(1)
	}

	err = os.MkdirAll(wdPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(color.GreenString("[+] Finding latest Architectury versions..."))

	latestArchPlugin, err := getLatestArchitecturyPlugin()
	if err != nil {
		log.Fatal("getLatestArchitecturyPlugin", err)
	}
	latestArchLoom, err := getLatestArchitecturyLoom()
	if err != nil {
		log.Fatal("getLatestArchitecturyLoom", err)
	}
	log.Println("Latest Arch Plugin:", latestArchPlugin)
	log.Println("Latest Arch Loom:", latestArchLoom)
	modInfo["architectury_plugin_version"] = latestArchPlugin
	modInfo["architectury_loom_version"] = latestArchLoom

	log.Println(color.GreenString("[+] Fetching version data..."))

	// rename and replace rest of template
	err = fs.WalkDir(templateLayers, ".", func(tempPath string, d fs.DirEntry, err error) error {
		replacedPath, err := modInfo.ReplaceInPath(tempPath)
		if err != nil {
			return err
		}
		fullPath := filepath.Join(wdPath, replacedPath)

		// skip directories
		if d.IsDir() {
			return nil
		}

		// create directory before file
		err = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return err
		}

		// open input from template
		openFile, err := templateLayers.Open(tempPath)
		if err != nil {
			return err
		}

		// open output file
		createFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		switch fileReplaceModes[tempPath] {
		case NormalReplace:
			_, err = io.Copy(createFile, modInfo.ReplaceInStream(openFile))
			if err != nil {
				return err
			}
		case NoReplace:
			_, err = io.Copy(createFile, openFile)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to walk files in template: ", err)
	}

	log.Println(color.HiGreenString("[+] Finished generating mod from template"))
}

type MavenMetadata struct {
	Versioning struct {
		Latest string `xml:"latest"`
	} `xml:"versioning"`
}

func getLatestMavenVersion(url string) (string, error) {
	var t MavenMetadata
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&t)
	return t.Versioning.Latest, err
}

func getLatestArchitecturyPlugin() (string, error) {
	return getLatestMavenVersion("https://maven.architectury.dev/architectury-plugin/architectury-plugin.gradle.plugin/maven-metadata.xml")
}

func getLatestArchitecturyLoom() (string, error) {
	return getLatestMavenVersion("https://maven.architectury.dev/dev/architectury/architectury-loom/maven-metadata.xml")
}
