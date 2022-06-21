# cifra-go-getter
Projeto em go para visualizar cifras

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

eagles
avantasia
avenged sevenfold
bobby darin
neil sakada
billy idol
scorpions
jorge e matheus
ze ramalho
caetano veloso
chico buarque
marisa monte
adriana calcanhoto
toquinho
rita lee
vanessa da mata
leoni
lenine
oswaldo montenegro
tie
banda do mar
a banda mais bonita da cidade
jota quest
coldplay
paramore
cassia eller
guns n roses
o rappa
pink floyd
oasis
barao vermelho
frejat
credence clear water revival
u2
linkin park
john mayer
led zeppelin
bee gees
dire straits
cogumelo plutao
papas da lingua
jason mraz
extreme
acdc
the police
avril lavigne
supercombo
the rolling stones
kiss
creed
perl jam
david bowie
rpm
roxxete
ramones
billy joel
kansas
the doors
counting crows
journey
paulo ricardo
