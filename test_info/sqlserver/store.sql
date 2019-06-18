USE [fruit]
GO
/****** Object: Table [dbo].[store] Script Date: 2019/6/14 9:29:21 ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
DROP TABLE IF EXISTS [dbo].[store]
GO

CREATE TABLE [dbo].[store](
[id] [bigint] IDENTITY(1,1) NOT NULL,
[code] [varchar](255) NULL,
[name] [nvarchar](255) NULL,
PRIMARY KEY CLUSTERED 
(
[id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
SET ANSI_PADDING OFF
GO
SET IDENTITY_INSERT [dbo].[store] ON

GO
INSERT [dbo].[store] ([id], [code], [name]) VALUES (1, N'10001', N'许仙网')
GO
INSERT [dbo].[store] ([id], [code], [name]) VALUES (2, N'10002', N'果想你')
GO
SET IDENTITY_INSERT [dbo].[store] OFF
GO

