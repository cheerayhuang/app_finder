CREATE TABLE IF NOT EXISTS apple_store_app (
    id bigint NOT NULL,
    bundleId varchar(64) NOT NULL,
    trackCensoredName varchar(128) NOT NULL,
    trackViewUrl varchar(256) NOT NULL,
    description MEDIUMTEXT NOT NULL,
    currency varchar(8) NOT NULL,
    price float NOT NULL,
    artistId bigint NOT NULL,
    artistName varchar(128) NOT NULL,
    sellerName varchar(128) NOT NULL,
    trackContentRating varchar(8) NOT NULL,
    averageUserRating float NOT NULL,
    userRatingCount bigint NOT NULL,
    `blob` MEDIUMTEXT NOT NULL,
    PRIMARY KEY(id),
    INDEX(bundleId, id)
)ENGINE = innoDB DEFAULT CHARACTER SET = utf8;
