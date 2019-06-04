USE [fruit]
GO
/****** Object: Table [dbo].[fruit] Script Date: 2019/6/14 9:29:21 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE [dbo].[fruit](
[id] [bigint] IDENTITY(1,1) NOT NULL,
[code] [varchar](255) NULL,
[name] [varchar](255) NULL,
[color] [varchar](255) NULL,
[price] [bigint] NULL,
[store_code] [varchar](255) NULL,
[created_at] [varchar](255) NULL,
[updated_at] [varchar](255) NULL,
PRIMARY KEY CLUSTERED 
(
[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
SET ANSI_PADDING OFF
GO
SET IDENTITY_INSERT [dbo].[fruit] ON

GO
INSERT [dbo].[fruit] ([id], [code], [name], [color], [price], [store_code], [created_at], [updated_at]) VALUES (1, N'apple', N'apple', N'red', 11, NULL, N'2018-04-15 08:48:59', N'2018-10-22 06:00:09')
GO
INSERT [dbo].[fruit] ([id], [code], [name], [color], [price], [store_code], [created_at], [updated_at]) VALUES (2, N'banana', N'banana', N'yellow', 14, NULL, N'2018-04-15 08:48:59', N'2018-10-22 06:00:09')
GO
SET IDENTITY_INSERT [dbo].[fruit] OFF
GO

