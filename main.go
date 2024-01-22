package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
)

var minecraftVersion = flag.String("mc", "latest", "latest, snapshot, 1.20, 1.20.4")

func main() {
	latestArch, err := getLatestArchitecturyPlugin()
	if err != nil {
		log.Fatal("getLatestArchitecturyPlugin", err)
	}
	log.Println(latestArch)
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
