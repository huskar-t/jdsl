CREATE DATABASE IF NOT EXISTS db KEEP 365
CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value Double) TAGS (device binary(64),gateway binary(64))
CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (device binary(64),gateway binary(64))
IMPORT INTO db.t_p1 USING db.stb_string TAGS ("d1","g2") VALUES (now,'2')
IMPORT INTO db.t_p2 USING db.stb_double TAGS ("d2","g2") VALUES (now,2)
IMPORT INTO db.t_p1 USING db.stb_string TAGS ("d1","g1") VALUES (now,'1')
IMPORT INTO db.t_p2 USING db.stb_double TAGS ("d2","g1") VALUES (now,1)