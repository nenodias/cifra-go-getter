# cifra-go-getter
Projeto em go para baixar cifras

Nesse exemplo o programa irá acessar a url https://www.cifraclub.com.br/queen/ buscando todas as músicas para violão/guitarra iterar as musicas e verificar todas as cifras de todas as versões.

A Struct que irá conter os dados das cifras
```main.go
type SongItem struct {
	Nome, Artista, Tom, Afinacao, Capo, Cifra, Versao string
}
```

```main.go
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
```
