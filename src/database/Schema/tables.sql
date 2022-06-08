
create database projet-pc3r;
use database projet-pc3r;

CREATE TABLE Utilisateur(
                            login varchar(10),
                            id bigint(5) ,
                            password varchar(15),
                            mail varchar(30),
                            cagnotte float,
                            derniereConnexion timestamp,
                            CONSTRAINT pk_login PRIMARY KEY(login)
);

CREATE TABLE Connexion(
                          login varchar(10),
                          id bigint(5) NOT NULL AUTO_INCREMENT,
                          idSession varchar(16),
                          date timestamp,
                          CONSTRAINT pk_session PRIMARY KEY(id),
                          CONSTRAINT fk_login_utilisateur FOREIGN KEY (login) references Utilisateur(login)
);

CREATE TABLE Matchs(
                       sport varchar(10),
                       league varchar(10),
                       equipeA varchar(60),
                       equipeB varchar(60),
                       date timestamp,
                       id bigint(5),
                       cote float,
                       statut enum('canceled','finished','not_started','postponed','running'),
                       vainqueur varchar(15),
                       CONSTRAINT unique_matchs_id UNIQUE(id),
                       CONSTRAINT pk_login PRIMARY KEY(sport,league,equipeA,equipeB,date)
);


CREATE TABLE Pari(
                     idMatch bigint(5),
                     equipeGagnante varchar(60),
                     login varchar(10),
                     id bigint(5),
                     cote float,
                     montant float,
                     resultat enum('win','loose','coming') default 'coming',
                     date timestamp,
                     CONSTRAINT unique_pari_id UNIQUE(id),
                     CONSTRAINT pk_pari PRIMARY KEY(idMatch, equipeGagnante, login),
                     CONSTRAINT fk_id_match_pari FOREIGN KEY (idMatch) REFERENCES Matchs(id),
                     CONSTRAINT fk_login_parieur FOREIGN KEY (login) references Utilisateur(login)
);


