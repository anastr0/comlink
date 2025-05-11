-- TODO : convert to migration file -->
-- TODO : better indexing in db -->
CREATE TABLE IF NOT EXISTS users (
  id bigserial primary key,
  name varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  id bigserial primary key,
  content varchar(256) NOT NULL,
  sender bigserial,
  receiver bigserial,
  read boolean DEFAULT false,
  timestamp timestamp NOT NULL,
  conversation varchar(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversations (
  id bigserial,
  key varchar(64) primary key,
  user1 bigserial references users(id),
  user2 bigserial references users(id)
)

CREATE INDEX conv_index ON messages(conversation);

UPDATE messages
SET read = 'false'
WHERE id = 129 returning *;
