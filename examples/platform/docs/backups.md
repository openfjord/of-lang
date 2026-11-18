# Backups

A backup holds the same data as its table, so it needs the same care. The
platform makes a backup of each table each day. Each backup has a day to
delete it and the sensitivity of its data:

![](../defs/backups.of.yaml#Backup)

The platform keeps each backup in a region, such as the European Union. A
region that does not allow data of that sensitivity is an error:

![](../defs/backups.of.yaml#backupInAllowedRegion)

This knowledge is only for the platform team. The package does not export
it, and a privacy review does not get it: its definitions have the tag
`backups`, not `privacy`.
