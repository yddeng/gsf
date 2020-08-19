# rank

基于 `Size Balanced Tree` 的排行榜，提供读写 `O(log n)` 的时间复杂度

> 基本规则是：由大到小排序，且得分相同的情况下，最新更新的 key 排名更靠前

> key,score 的数据格式根据不同需求做相应修改，现有格式为 key:int64, score:int64
