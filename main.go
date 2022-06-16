package main

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

var baseUrl string = "https://www.cifraclub.com.br"

type SongListItem struct {
	Titulo, Link string
}

type SongItem struct {
	Nome, Artista, Tom, Afinacao, Capo, Cifra, Versao string
}

type SongVersion struct {
	Nome, Artista, Versao, Dificuldade, Link string
}

func main() {
	artist := "/queen"
	for _, song := range OpenSongList(artist) {
		for _, version := range GetSongVersions(song.Link) {
			fmt.Println(version.Nome)
			fmt.Println(version.Versao)
			chords, err := OpenSong(song.Link)
			if err == nil {
				fmt.Println(chords.Cifra)
			}
		}
	}
}

func GetSongVersions(songLink string) []SongVersion {
	songResp, err := soup.Get(baseUrl + songLink)
	response := []SongVersion{}
	if err == nil {
		songDoc := soup.HTMLParse(songResp)
		songVersionsDiv := songDoc.Find("div", "class", "list-versions")
		songVersions := songVersionsDiv.FindAll("a")
		for _, version := range songVersions {
			names, _ := GetNameAndArtist(songDoc)
			name := names[0]
			artist := names[1]
			versionAttr := version.Attrs()
			infos := version.FindAll("span")
			link := versionAttr["href"]
			versionName := infos[0].Text()
			dificuldade := infos[1].Text()
			response = append(response, SongVersion{
				Nome:        strings.TrimSpace(name),
				Artista:     strings.TrimSpace(artist),
				Versao:      strings.TrimSpace(versionName),
				Dificuldade: strings.TrimSpace(dificuldade),
				Link:        link,
			})
		}
	}
	return response
}

func GetNameAndArtist(songDoc soup.Root) ([]string, error) {
	name := songDoc.Find("h1", "class", "t1").Text()
	artistH2 := songDoc.Find("h2", "class", "t3")
	artist := artistH2.Find("a").Text()
	return []string{name, artist}, nil
}

func OpenSong(songLink string) (SongItem, error) {
	songResp, err := soup.Get(baseUrl + songLink)
	if err == nil {
		songDoc := soup.HTMLParse(songResp)

		names, _ := GetNameAndArtist(songDoc)
		name := names[0]
		artist := names[1]
		version := songDoc.Find("a", "id", "js-c-versions").Text()

		cifraDiv := songDoc.Find("div", "class", "cifra-mono")
		tomSpan := cifraDiv.Find("span", "id", "cifra_tom")
		tom := tomSpan.Find("a")
		afinacao := cifraDiv.Find("span", "id", "cifra_afi")
		capo := cifraDiv.Find("span", "id", "cifra_capo")
		cifra := cifraDiv.Find("pre")
		return SongItem{
			Nome:     strings.TrimSpace(name),
			Artista:  strings.TrimSpace(artist),
			Tom:      strings.TrimSpace(tom.Text()),
			Afinacao: strings.TrimSpace(afinacao.Text()),
			Capo:     strings.TrimSpace(capo.Text()),
			Cifra:    cifra.HTML(),
			Versao:   strings.TrimSpace(version),
		}, nil
	} else {
		return SongItem{}, err
	}
}

func OpenSongList(artist string) []SongListItem {
	response := []SongListItem{}
	resp, err := soup.Get(baseUrl + artist)
	if err != nil {
		return response
	}
	doc := soup.HTMLParse(resp)
	ul := doc.Find("ul", "id", "js-a-songs")
	songs := ul.FindAll("li")
	for _, songLi := range songs {
		a := songLi.Find("a")
		songSpan := songLi.Find("span")
		guitar := songSpan.FindAll("a")[0]
		guitarAttributes := guitar.Attrs()
		hasGuitar := guitarAttributes["data-ajax"]
		if hasGuitar != "false" {
			attributes := a.Attrs()
			title := attributes["title"]
			songLink := attributes["href"]
			response = append(response, SongListItem{
				Titulo: strings.TrimSpace(title),
				Link:   songLink,
			})
		}
	}
	return response
}
