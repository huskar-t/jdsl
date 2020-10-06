CREATE DATABASE IF NOT EXISTS db KEEP 365
CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value Double) TAGS (type binary(64))
CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (type binary(64))
IMPORT INTO db.t_zone USING db.stb_string TAGS ("temperature") VALUES (now,'15')
IMPORT INTO db.t_zone USING db.stb_string TAGS ("humidity") VALUES (now,'17')