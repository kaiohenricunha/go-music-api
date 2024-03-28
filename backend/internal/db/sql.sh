# create database
create database infnet_music_db;

# use database
use infnet_music_db;

# insert into table songs
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Burn It Down", "Linkin Park");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Earth Song", "Michael Jackson");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Hey Jude", "The Beatles");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Sound of Silence", "Simon and Garfunkel");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Hotel California", "The Eagles");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Comfortably Numb", "Pink Floyd");
INSERT INTO songs (created_at, updated_at, name, artist) VALUES (now(), now(), "Borbulhas de Amor", "Fagner");

# list all songs
SELECT * FROM songs;
