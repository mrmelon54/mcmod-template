package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/fatih/color"
	"github.com/mrmelon54/mcmodupdater"
	mcmConfig "github.com/mrmelon54/mcmodupdater/config"
	"github.com/mrmelon54/mcmodupdater/develop"
	"github.com/mrmelon54/mcmodupdater/develop/dev"
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
	platforms     = []develop.DevPlatform{
		dev.PlatformFabric,
		dev.PlatformForge,
	}
	propVersions = []develop.PropVersion{
		develop.ModVersion,
		develop.ArchitecturyVersion,
		develop.FabricLoaderVersion,
		develop.FabricApiVersion,
		develop.ForgeVersion,
		develop.QuiltLoaderVersion,
		develop.QuiltFabricApiVersion,
		develop.NeoForgeVersion,
	}
)

//go:embed all:template
var templateDir embed.FS

func prompt(s string) string {
	_, _ = questionColor.Print(s)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func fakePrompt(s string, v string) {
	_, _ = questionColor.Print(s)
	fmt.Println(v)
}

type Config struct {
	ModGroupBase  string           `json:"mod_group_base"`
	ModSiteBase   string           `json:"mod_site_base"`
	ModSourceBase string           `json:"mod_source_base"`
	McmConfig     mcmConfig.Config `json:"mcm_config"`
}

func main() {
	conf := Config{
		ModGroupBase:  "com.example",
		ModSiteBase:   "https://example.com/minecraft",
		ModSourceBase: "https://github.com/example",
		McmConfig:     mcmConfig.DefaultConfig(),
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

	mcm, err := mcmodupdater.NewMcModUpdater(&conf.McmConfig)
	if err != nil {
		log.Fatal(err)
	}

	mcVersion := prompt("[?] Minecraft Version (1.20, 1.20.4): ")
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
	modInfo["minecraft_version"] = modName
	modInfo["architectury_plugin_version"] = modName
	modInfo["architectury_loom_version"] = modName
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

	switch prompt("Is that ok [y/N]? ") {
	case "y", "Y":
		break
	default:
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

	log.Println(color.GreenString("[+] Fetching version data..."))

	// add subplatforms to architectury
	mcm.PlatArch().SubPlatforms = make(map[develop.DevPlatform]develop.Develop)
	for _, i := range platforms {
		mcm.PlatArch().SubPlatforms[i] = mcm.Platforms()[i]
	}

	// fetch architectury specific caches first
	err = fetchCalls(mcm.PlatArch())
	if err != nil {
		log.Println(color.HiRedString("Error: %s", err))
		os.Exit(1)
	}

	// fetch sub-platform caches
	for _, i := range dev.Platforms {
		if c, ok := mcm.PlatArch().SubPlatforms[i]; ok {
			err := fetchCalls(c)
			if err != nil {
				log.Println(color.HiRedString("Error: %s", err))
				os.Exit(1)
			}
		}
	}

	// rename and replace rest of template
	err = fs.WalkDir(templateDir, "template", func(tempPath string, d fs.DirEntry, err error) error {
		relPath := strings.TrimPrefix(tempPath, "template/")
		replacedPath, err := modInfo.ReplaceInString(relPath)
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
		openFile, err := templateDir.Open(tempPath)
		if err != nil {
			return err
		}

		// open output file
		createFile, err := os.Create(fullPath)
		if err != nil {
			return err
		}

		switch fileReplaceModes[relPath] {
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
		case PropertiesReplace:
			// hard code chosen Minecraft version
			infoVersions := make(map[develop.PropVersion]string)
			infoVersions[develop.MinecraftVersion] = mcVersion
			for _, i := range propVersions {
				infoVersions[i] = " "
			}

			// generate update list
			ver := mcm.VersionUpdateList(&develop.PlatformVersions{
				Platform: mcm.PlatArch(),
				Versions: infoVersions,
			})
			tempGradleProp := modInfo.ReplaceInStream(openFile)
			err = mcm.UpdateGradleProperties(createFile, tempGradleProp, ver.ChangeToLatest())
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

func fetchCalls(platform develop.Develop) error {
	for _, i := range platform.FetchCalls() {
		err := i.Call()
		if err != nil {
			return err
		}
	}
	return nil
}
