CREATE DATABASE IF NOT EXISTS db KEEP 365
CREATE TABLE IF NOT EXISTS db.stb_double (ts timestamp, value Double) TAGS (device binary(64))
CREATE TABLE IF NOT EXISTS db.stb_string (ts timestamp, value NCHAR(64)) TAGS (device binary(64))
IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,92)
IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,93)
IMPORT INTO db.t_sunshine USING db.stb_double TAGS ("d1") VALUES (now,94)