
USE master;
CREATE DATABASE fruit;
USE fruit;
IF  NOT EXISTS (SELECT * FROM sysobjects WHERE id = object_id(N'fruit'))
    BEGIN
        CREATE TABLE fruit(
        [id] [bigint] NULL,
        [code] [varchar](50) NULL,
        [name] [varchar](50) NULL,
        [color] [varchar](50) NULL,
        [price] [decimal](18, 2) NULL,
        [store_code] [varchar](50) NULL,
        [created_at] [datetime] NULL,
        [updated_at] [datetime] NULL
        )
    END;
-- 导出  表 fruit.fruit 结构
TRUNCATE TABLE fruit;

INSERT INTO fruit (id, created_at, updated_at, code, name, color, price, store_code) VALUES
	(1, '2018-04-15 08:48:59', '2018-10-22 06:00:09', 'apple', 'apple', 'red', 11, NULL);
INSERT INTO [fruit] (id, created_at, updated_at, code, name, color, price, store_code) VALUES
	(2, '2018-04-15 08:48:59', '2018-10-22 06:00:09', 'banana', 'banana', 'yellow', 14, NULL);

