CREATE TABLE musica (
    id integer primary key AUTO_INCREMENT not null,
    nome varchar(255) not null,
    artista varchar(255) not null,
    tom varchar(5),
    afinacao varchar(50),
    capo varchar(50),
    cifra text,
    versao varchar(255) not null
);