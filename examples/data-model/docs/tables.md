# Tables

A table holds rows with one schema. Each table is in exactly one database.
The table belongs to the team that owns that database:

![](../defs/table.of.yaml#Table)

![](../defs/table.of.yaml#tableTeamFromDatabase)

Teams often caused failures for other teams when they read internal tables.
Thus, an application must have a contract to read a table from a different
application. Refer to
[crossAppReadNeedsContract](../defs/reading-across-applications.of.yaml#crossAppReadNeedsContract).
