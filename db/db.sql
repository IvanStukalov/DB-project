CREATE EXTENSION IF NOT EXISTS CITEXT; -- eliminate calls to lower

CREATE UNLOGGED TABLE users
(
    Nickname   CITEXT PRIMARY KEY,
    FullName   TEXT NOT NULL,
    About      TEXT NOT NULL DEFAULT '',
    Email      CITEXT UNIQUE
);

CREATE UNLOGGED TABLE forum
(
    Title    TEXT   NOT NULL,
    "user"   CITEXT,
    Slug     CITEXT PRIMARY KEY,
    Posts    INT    DEFAULT 0,
    Threads  INT    DEFAULT 0
);

CREATE UNLOGGED TABLE threads
(
    Id      SERIAL    PRIMARY KEY,
    Title   TEXT      NOT NULL,
    Author  CITEXT    REFERENCES "users"(Nickname),
    Forum   CITEXT    REFERENCES "forum"(Slug),
    Message TEXT      NOT NULL,
    Votes   INT       DEFAULT 0,
    Slug    CITEXT,
    Created TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE UNLOGGED TABLE posts
(
    Id        SERIAL      PRIMARY KEY,
    Author    CITEXT,
    Created   TIMESTAMP   WITH TIME ZONE DEFAULT now(),
    Forum     CITEXT,
    IsEdited  BOOLEAN     DEFAULT FALSE,
    Message   CITEXT      NOT NULL,
    Parent    INT         DEFAULT 0,
    Thread    INT,
    Path      INTEGER[],
    FOREIGN KEY (thread) REFERENCES "threads" (id),
    FOREIGN KEY (author) REFERENCES "users"  (nickname)
);

CREATE UNLOGGED TABLE votes
(
    ID       SERIAL PRIMARY KEY,
    Author   CITEXT    REFERENCES "users" (Nickname),
    Voice    INT       NOT NULL,
    Thread   INT,
    FOREIGN KEY (thread) REFERENCES "threads" (id),
    UNIQUE (Author, Thread)
);


CREATE UNLOGGED TABLE users_forum
(
    Nickname  CITEXT  NOT NULL,
    FullName  TEXT    NOT NULL,
    About     TEXT,
    Email     CITEXT,
    Slug      CITEXT  NOT NULL,
    FOREIGN KEY (Nickname) REFERENCES "users" (Nickname),
    FOREIGN KEY (Slug) REFERENCES "forum" (Slug),
    UNIQUE (Nickname, Slug)
);


CREATE OR REPLACE FUNCTION addThread() RETURNS TRIGGER AS
$update_forum$
BEGIN
    UPDATE forum SET Threads=(Threads+1) WHERE Slug = NEW.Forum;
    return NEW;
END
$update_forum$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_thread
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE addThread();


CREATE OR REPLACE FUNCTION addPost() RETURNS TRIGGER AS
$update_forum$
DECLARE
    parent_path   INTEGER[];
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT path FROM posts WHERE id = NEW.parent INTO parent_path;
        NEW.path := parent_path || NEW.id;
    END IF;
    UPDATE forum SET Posts=(Posts+1) WHERE Slug = NEW.Forum;
    return NEW;
END
$update_forum$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_post
    BEFORE INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE addPost();


CREATE OR REPLACE FUNCTION addVote() RETURNS TRIGGER AS
$update_forum$
BEGIN
    UPDATE threads SET Votes=(Votes+New.Voice) WHERE Id = NEW.Thread;
    return NEW;
END
$update_forum$ LANGUAGE plpgsql;

CREATE TRIGGER on_insert_vote
    AFTER INSERT
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE addVote();


CREATE OR REPLACE FUNCTION changeVote() RETURNS TRIGGER AS
$update_forum$
BEGIN
UPDATE threads SET Votes=(Votes+2*New.Voice) WHERE Id = NEW.Thread;
return NEW;
END
$update_forum$ LANGUAGE plpgsql;

CREATE TRIGGER on_update_vote
    AFTER UPDATE
    ON votes
    FOR EACH ROW
    EXECUTE PROCEDURE changeVote();


CREATE OR REPLACE FUNCTION UpdateUserOnPost() RETURNS TRIGGER AS
$update_users_on_post$
DECLARE
    author_fullname CITEXT;
    author_about    CITEXT;
    author_email    CITEXT;
BEGIN
    SELECT Fullname, About, Email FROM users WHERE Nickname = NEW.Author INTO author_fullname, author_about, author_email;
    INSERT INTO users_forum (Nickname, Fullname, About, Email, Slug)
        VALUES (NEW.Author, author_fullname, author_about, author_email, NEW.Forum)
        ON CONFLICT DO NOTHING;
    return NEW;
END
$update_users_on_post$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_on_post
    AFTER INSERT
    ON posts
    FOR EACH ROW
    EXECUTE PROCEDURE UpdateUserOnPost();


CREATE OR REPLACE FUNCTION UpdateUserForum() RETURNS TRIGGER AS
$update_uf$
DECLARE
    author_fullname CITEXT;
    author_about    CITEXT;
    author_email    CITEXT;
BEGIN
    SELECT Fullname, About, Email FROM users WHERE Nickname = NEW.Author INTO author_fullname, author_about, author_email;
    INSERT INTO users_forum (Nickname, Fullname, About, Email, Slug)
        VALUES (NEW.Author, author_fullname, author_about, author_email, NEW.Forum)
        ON CONFLICT DO NOTHING;
    return NEW;
END
$update_uf$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
    EXECUTE PROCEDURE UpdateUserForum();