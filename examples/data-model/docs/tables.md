# Tables

A table holds rows with one schema. Every table sits in exactly one database,
and belongs to the team that owns that database:

![](../defs/table.yaml#Table)

![](../defs/table.yaml#tableTeamFromDatabase)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../defs/reading-across-applications.yaml#crossAppReadNeedsContract).
