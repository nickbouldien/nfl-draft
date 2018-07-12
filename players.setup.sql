-- TODO: make non null
CREATE TABLE players (
  id SERIAL PRIMARY KEY,
  full_name VARCHAR(50) NOT NULL,
  school VARCHAR(50),
  position VARCHAR(2),
  drafted boolean NOT NULL DEFAULT FALSE
);

-- \copy zip_codes(full_name,school,position,drafted) FROM '/path/to/csv/ZIP_CODES.txt' DELIMITER ',' CSV


-- ENUM
-- https://www.postgresql.org/docs/9.1/static/datatype-enum.html
-- CREATE TYPE positon AS ENUM ('QB, RB, WR, DB, DL, OL, K, LB, P');

-- insert data:
INSERT INTO players(full_name,school,position,drafted) VALUES ('Coznesster Smiff','Rutgers University','QB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Elipses Corter','University of Alabama','OL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Nyquillus Dillwad','LSU','LB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Bismo Funyuns','Florida State University','DB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Decatholac Mango','Georgia Tech University','OL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Mergatroid Skittle','University of Louisville','TE',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Quiznatodd Bidness','University of Tennessee','WR',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('D''Pez Poopsie','Ole Miss','DL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Quackadilly Blip','Auburn University','LB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Goolius Boozler','The U','K',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Bisquiteen Trisket','University of Michigan','RB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Fartrell Cluggins','Arkansas University','OL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Blyrone Blashinton','Syracuse University','DL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Cartoons Plural','Virginia Tech University','QB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Jammie Jammie-Jammie','The Ohio State University','DB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Fudge','N/A','WR',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Equine Ducklings','Indiana University Purdue University Indianapolis','P',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Dahistorius Lamystorius','Utah State University','RB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Ewokoniad Sigourneth JuniorStein','Oklahoma State University','OL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Eqqsnuizitine Buble-Schwinslow','University of Nebraska','WR',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Huka''lakanaka Hakanakaheekalucka''hukahakafaka','University of Hawaii','LB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('King Prince Chambermaid','Baylor University','DL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Ladennifer Jadaniston','University of Colorado','DB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Ladadadaladadadadada Dala-Dadaladaladalada','University of Arizona','DB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Harvard University','DeVry University','DL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Morse Code','Army Navay Surplus Store','QB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Wingdings','Online Classes','LB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Firstname Lastname','College University','DL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('God','Home schooled','RB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Squeeeeeeeeeeps','Santa Monica College','RB',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('Benedict Cumberbatch','Oxford University','OL',FALSE);
INSERT INTO players(full_name,school,position,drafted) VALUES ('A.A. Ron Balakay','Morehouse College','WR',FALSE);
