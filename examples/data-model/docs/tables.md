# Tables

A table holds rows with one schema. Every table sits in exactly one database,
and belongs to the team that owns that database:

![](../definitions/table.yaml#Table)

![](../definitions/table.yaml#tableTeamFromDatabase)

Teams kept breaking each other by reading internal tables, so reading across
applications needs a contract — see
[crossAppReadNeedsContract](../definitions/reading-across-applications.yaml#crossAppReadNeedsContract).
