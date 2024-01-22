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

func getLatestMavenVersion(url string) (T, error) {
	var t T
	resp, err := http.Get(url)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()
	dec := xml.NewDecoder(resp.Body)
	err = dec.Decode(&t)
	return t, err
}

type MavenMetadata struct {
	Versioning struct {
		Latest string `xml:"latest"`
	} `xml:"versioning"`
}

func getLatestArchitecturyPlugin() (string, error) {
	m, err := decodeMavenXml[MavenArchitecturyPlugin]("https://maven.architectury.dev/architectury-plugin/architectury-plugin.gradle.plugin/maven-metadata.xml")
	return m.Versioning.Latest, err
}

type MavenArchitecturyLoom struct {
}

func getLatestArchitecturyLoom() (string, error) {
	m, err := decodeMavenXml[MavenArchitecturyLoom]("https://maven.architectury.dev/dev/architectury/architectury-loom/maven-metadata.xml")

}
