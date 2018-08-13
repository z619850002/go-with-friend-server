 -- --------------------------------------------------
    --  Table Structure for `myproject/models.Player`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`player` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `nick_name` varchar(50) NOT NULL DEFAULT '' ,
        `avatar` varchar(255) NOT NULL DEFAULT '' ,
        `avatar_edge` varchar(255) NOT NULL DEFAULT ''
        `isonline` bool NOT NULL DEFAULT FALSE ,
        `dan` integer NOT NULL DEFAULT 1 ,
    ) ENGINE=InnoDB;


 -- --------------------------------------------------
    --  Table Structure for `myproject/models.Account`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`account` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `player_id` integer NOT NULL,
        `identity_type` varchar(50) NOT NULL DEFAULT '' ,
        `identifier` varchar(100) NOT NULL DEFAULT '' ,
        `credential` varchar(100) NOT NULL DEFAULT '' ,
        FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;



-- --------------------------------------------------
    --  Table Structure for `go-with-friend-server/src/Server/DB/models.Level`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`level` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `player_id` integer NOT NULL UNIQUE,
        `experience` bigint NOT NULL DEFAULT 0,
        `fall_number` bigint NOT NULL DEFAULT 0 ,
        `take_number` bigint NOT NULL DEFAULT 0 ,
        `practice_info` varchar(500) NOT NULL DEFAULT '' ,
        `friends_battle` integer NOT NULL DEFAULT 0 ,
        `smart_battle` integer NOT NULL DEFAULT 0 ,
        `circumstance` integer NOT NULL DEFAULT 0 ,
        FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;

    -- --------------------------------------------------
    --  Table Structure for `myproject/models.History`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`history` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `player1_id` integer NOT NULL,
        `player2_id` integer NOT NULL,
        `url` varchar(100) NOT NULL DEFAULT '' ,
        `commentsnumber` INTEGER NOT NULL DEFAULT 0,
        `bystanders` integer NOT NULL DEFAULT 0 ,
        `like` integer NOT NULL DEFAULT 0 ,
        `step` integer NOT NULL DEFAULT 0 ,
        `winner` integer NOT NULL DEFAULT 0 ,
        `size` integer NOT NULL DEFAULT 0 ,
        `type` integer NOT NULL DEFAULT 0 ,
        `time` datetime NOT NULL,
        `player1_score` integer NOT NULL default 0,
        `player2_score` integer NOT NULL default 0,
        -- 对战中的相关信息
        `chess_board` varchar(500) NOT NULL DEFAULT '' ,
        `turn` bool NOT NULL DEFAULT FALSE ,
        `last_turn_time` datetime NOT NULL,         -- 结束之后作为棋局结束的时间
        FOREIGN KEY(player1_id) REFERENCES player(id) ON DELETE CASCADE,
        FOREIGN KEY(player2_id) REFERENCES player(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;

 
    -- --------------------------------------------------
    --  Table Structure for `myproject/models.GoodChoice`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`good_choice` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `histories_id` integer NOT NULL,
        `pos` integer NOT NULL DEFAULT 0 ,
        `praise_number` integer NOT NULL DEFAULT 0,
        FOREIGN KEY(histories_id) REFERENCES history(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;

    -- --------------------------------------------------
    --  Table Structure for `myproject/models.Message`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`message` (
        `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `history_id` integer NOT NULL,
        `player_id` integer NOT NULL,
        `data` varchar(100) NOT NULL DEFAULT '',
        FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE,
        FOREIGN KEY(history_id)  REFERENCES history(id)  ON DELETE CASCADE
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `myproject/models.PlayerHistorys`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`collected_history` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `player_id` integer NOT NULL,
        `history_id` integer NOT NULL,
        FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE,
        FOREIGN KEY(history_id) REFERENCES history(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;


    -- --------------------------------------------------
    --  Table Structure for `myproject/models.PlayerGoodChoices`
    -- --------------------------------------------------
    CREATE TABLE IF NOT EXISTS `test`.`praised_goodchoice` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `player_id` integer NOT NULL,
        `good_choice_id` integer NOT NULL,
        FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE,
        FOREIGN KEY(good_choice_id) REFERENCES good_choice(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;
    
    
    
	CREATE TABLE IF NOT EXISTS `test`.`invitation` (
        `id` bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `inviter_id` integer NOT NULL,
        `invitee_id` integer NOT NULL,
        `turn` bool NOT NULL DEFAULT FALSE ,
        `time` datetime NOT NULL DEFAULT NOW(),
        `size` integer NOT NULL DEFAULT 0 ,
        `firststep` VARCHAR(10) NOT NULL DEFAULT '0',         -- 黑子的位置
        `type` VARCHAR(50) NOT NULL DEFAULT '',               --邀请方式
        FOREIGN KEY(inviter_id) REFERENCES player(id) ON DELETE CASCADE,
        FOREIGN KEY(invitee_id) REFERENCES player(id) ON DELETE CASCADE
    ) ENGINE=InnoDB;



   