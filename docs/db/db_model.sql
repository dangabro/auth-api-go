-- =============================================================================
-- Diagram Name: db_model
-- Created on: 2/11/2019 4:06:05 PM
-- Diagram Version: 
-- =============================================================================
SET FOREIGN_KEY_CHECKS=0;

CREATE TABLE `right` (
  `right_cd` varchar(64) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY(`right_cd`)
)
ENGINE=INNODB;

CREATE TABLE `user` (
  `user_id` varchar(64) NOT NULL,
  `name` varchar(128) NOT NULL,
  `login` varchar(64) NOT NULL,
  `password` varchar(128),
  PRIMARY KEY(`user_id`),
  UNIQUE INDEX `ix_user_login_unique`(`login`)
)
ENGINE=INNODB;

CREATE TABLE `token` (
  `token_id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `token` varchar(64) NOT NULL,
  `added_dt` datetime NOT NULL,
  `expired_ind` varchar(1) NOT NULL,
  `expiry_dt` datetime NOT NULL,
  PRIMARY KEY(`token_id`),
  UNIQUE INDEX `ix_token`(`token`),
  CONSTRAINT `fk_token_user` FOREIGN KEY (`user_id`)
    REFERENCES `user`(`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE=INNODB;

CREATE TABLE `user_right` (
  `user_right_id` varchar(64) NOT NULL,
  `user_id` varchar(64) NOT NULL,
  `right_cd` varchar(64) NOT NULL,
  PRIMARY KEY(`user_right_id`),
  UNIQUE INDEX `ix_user_right_unique`(`user_id`, `right_cd`),
  CONSTRAINT `fk_user_right_user` FOREIGN KEY (`user_id`)
    REFERENCES `user`(`user_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_right_right` FOREIGN KEY (`right_cd`)
    REFERENCES `right`(`right_cd`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION
)
ENGINE=INNODB;

SET FOREIGN_KEY_CHECKS=1;
