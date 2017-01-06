CREATE TABLE apple_store_app (
    id bigint NOT NULL,
    bundleId varchar(256) NOT NULL,
    trackCensoredName varchar(256) NOT NULL,
    trackViewUrl varchar(256) NOT NULL,
    genre1 varchar(64) NOT NULL,
    genre2 varchar(64) NOT NULL,
    genre3 varchar(64) NOT NULL,
    genre4 varchar(64) NOT NULL,
    currency varchar(8) NOT NULL,
    price real NOT NULL,
    artistId bigint NOT NULL,
    artistName varchar(128) NOT NULL,
    sellerName varchar(128) NOT NULL,
    trackContentRating varchar(8) NOT NULL,
    averageUserRating real NOT NULL,
    userRatingCount bigint NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE google_play_app (
    bundleId varchar(256) NOT NULL,
    trackCensoredName varchar(256) NOT NULL,
    trackViewUrl varchar(256) NOT NULL,
    genre varchar(64) NOT NULL,
    artistName varchar(128) NOT NULL,
    trackContentRating varchar(32) NOT NULL,
    averageUserRating real NOT NULL,
    userRatingCount bigint NOT NULL,
    id bigint NOT NULL,
    currency varchar(8) NOT NULL,
    price real NOT NULL,
    PRIMARY KEY(bundleId)
);
