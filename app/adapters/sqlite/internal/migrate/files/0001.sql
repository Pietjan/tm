CREATE TABLE [project] (
  [uid]         TEXT     NOT NULL  PRIMARY KEY,
  [ref]         TEXT,
  [name]        TEXT     NOT NULL,
  [created_at]  INTEGER  NOT NULL,
  [updated_at]  INTEGER
);

CREATE TABLE [task] (
  [uid]         TEXT     NOT NULL  PRIMARY KEY,
  [ref]         TEXT     NOT NULL,
  [name]        TEXT     NOT NULL,
  [status]      INTEGER  NOT NULL,
  [project]     TEXT,
  [created_at]  INTEGER  NOT NULL,
  [updated_at]  INTEGER
);

CREATE TABLE [timer] (
  [uid]         TEXT     NOT NULL  PRIMARY KEY,
  [task]        TEXT     NOT NULL,
  [start]       INTEGER  NOT NULL,
  [end]         INTEGER
);