# Spoilage

Food spoils when cooks keep it too warm or too long. The kitchen receives
each ingredient in batches. Each batch has a use-by day and the highest
temperature that it can stay at:

![](../defs/spoilage.of.yaml#Batch)

The cooks keep each batch in a storage zone, such as the fridge. A zone that
is too warm for a batch is an error:

![](../defs/spoilage.of.yaml#batchKeptColdEnough)

This knowledge is only for the kitchen. The package does not export it, and
the allergen app does not get it: its definitions have the tag `spoilage`,
not `allergens`.
