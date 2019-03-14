CREATE TABLE team (
  id SERIAL PRIMARY KEY,
  name        VARCHAR(50) NOT NULL,
  conference  VARCHAR(10) NOT NULL,
  division    VARCHAR(5) NOT NULL,
  draft_order INTEGER  NOT NULL
);

INSERT INTO team(name,conference,division,draft_order) VALUES ('Arizona Cardinals','NFC','West',1);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Atlanta Falcons','NFC','South',14);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Baltimore Ravens','AFC','North',22);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Buffalo Bills','AFC','East',9);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Carolina Panthers','NFC','South',16);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Chicago Bears','NFC','North',24);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Cincinnati Bengals','AFC','North',11);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Cleveland Browns','AFC','North',17);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Dallas Cowboys','NFC','East',27);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Denver Broncos','AFC','West',10);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Detroit Lions','NFC','North',8);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Green Bay Packers','NFC','North',12);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Houston Texans','AFC','South',23);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Indianapolis Colts','AFC','South',26);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Jacksonville Jaguars','AFC','South',7);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Kansas City Chiefs','AFC','West',29);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Miami Dolphins','AFC','East',13);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Los Angeles Chargers','AFC','West',28);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Los Angeles Rams','NFC','West',31);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Minnesota Vikings','NFC','North',18);
INSERT INTO team(name,conference,division,draft_order) VALUES ('New England Patriots','AFC','East',32);
INSERT INTO team(name,conference,division,draft_order) VALUES ('New Orleans Saints','NFC','South',30);
INSERT INTO team(name,conference,division,draft_order) VALUES ('NY Giants','NFC','East',6);
INSERT INTO team(name,conference,division,draft_order) VALUES ('NY Jets','AFC','East',3);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Oakland Raiders','AFC','West',4);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Philadelphia Eagles','NFC','East',25);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Pittsburgh Steelers','AFC','North',20);
INSERT INTO team(name,conference,division,draft_order) VALUES ('San Francisco 49ers','NFC','West',2);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Seattle Seahawks','NFC','West',21);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Tampa Bay Buccaneers','NFC','South',5);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Tennessee Titans','AFC','South',19);
INSERT INTO team(name,conference,division,draft_order) VALUES ('Washington Redskins','NFC','East',15);
