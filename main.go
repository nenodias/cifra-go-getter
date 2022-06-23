package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/nenodias/cifra-go-getter/dados"
	"github.com/nenodias/cifra-go-getter/db"
)

var baseUrl string = "https://www.cifraclub.com.br"

type SongListItem struct {
	Titulo, Link string
}

type SongItem struct {
	Id       uint   `gorm:"column:id;primarykey"`
	Nome     string `gorm:"column:nome;type:varchar(255);not null" json:"nome"`
	Artista  string `gorm:"column:artista;type:varchar(255);not null" json:"artista"`
	Tom      string `gorm:"column:tom;type:varchar(5)" json:"tom"`
	Afinacao string `gorm:"column:afinacao;type:varchar(50)" json:"afinacao"`
	Capo     string `gorm:"column:capo;type:varchar(50)" json:"capo"`
	Cifra    string `gorm:"column:cifra;type:text" json:"cifra"`
	Versao   string `gorm:"column:versao;type:varchar(255);not null" json:"versao"`
}

func (SongItem) TableName() string {
	return "musica"
}

type SongVersion struct {
	Nome, Artista, Versao, Dificuldade, Link string
}

func main() {
	db.Init()
	artists := dados.GetDados()
	for _, artist := range artists {
		time.Sleep(10 * time.Second)
		for _, song := range OpenSongList(artist) {
			for _, version := range GetSongVersions(song.Link) {
				fmt.Println(version.Nome, version.Artista, version.Versao)
				chords, err := OpenSong(version.Link)
				if err == nil {
					tx := db.DB.Save(&chords)
					tx.Commit()
				}
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
		if songVersionsDiv.Pointer != nil {
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
		} else {
			version := songDoc.Find("a", "class", "js-modal-trigger")
			names, _ := GetNameAndArtist(songDoc)
			name := names[0]
			artist := names[1]
			infos := version.FindAll("span")
			versionName := infos[0].Text()
			var dificuldade string
			if len(infos) > 1 {
				dificuldade = infos[1].Text()
			}
			response = append(response, SongVersion{
				Nome:        strings.TrimSpace(name),
				Artista:     strings.TrimSpace(artist),
				Versao:      strings.TrimSpace(versionName),
				Dificuldade: strings.TrimSpace(dificuldade),
				Link:        songLink,
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
		if cifraDiv.Pointer == nil {
			cifraDiv = songDoc.Find("div", "class", "cifra_cnt")
		}
		tomSpan := cifraDiv.Find("span", "id", "cifra_tom")
		tom := tomSpan.Find("a")
		afinacao := cifraDiv.Find("span", "id", "cifra_afi")
		afinacaoLink := afinacao.Find("a")
		capo := cifraDiv.Find("span", "id", "cifra_capo")
		capoLink := capo.Find("a")
		cifra := cifraDiv.Find("pre")

		afinacaoTxt := afinacao.Text()
		if afinacaoLink.Pointer != nil {
			afinacaoTxt += " " + afinacaoLink.Text()
		}
		capoTxt := capo.Text()
		if capoLink.Pointer != nil {
			capoTxt += " " + capoLink.Text()
		}

		return SongItem{
			Nome:     strings.TrimSpace(name),
			Artista:  strings.TrimSpace(artist),
			Tom:      strings.TrimSpace(tom.Text()),
			Afinacao: strings.TrimSpace(afinacaoTxt),
			Capo:     strings.TrimSpace(capoTxt),
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
	if ul.Pointer != nil {
		songs := ul.FindAll("li")
		response = DealWithSongList(songs)
	} else {
		ol := doc.Find("ol", "id", "js-a-t")
		songs := ol.FindAll("li")
		response = DealWithSongList(songs)
	}
	return response
}

func DealWithSongList(songs []soup.Root) []SongListItem {
	response := []SongListItem{}
	for _, songLi := range songs {
		a := songLi.Find("a")
		songSpan := songLi.Find("span")
		if songSpan.Pointer != nil {
			guitarLink := songSpan.FindAll("a")
			if len(guitarLink) > 0 {
				guitar := guitarLink[0]
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
		}
	}
	return response
}
