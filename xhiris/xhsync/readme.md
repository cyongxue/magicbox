## 说明
有些并发场景，但是需要基于关键key，进行同步执行，即保证串行执行。

## 问题
1. 扩容、缩容的时候如何处理？
    - 利用一致性hash保证尽量少的迁移，但是依然存在迁移问题？该如何处理？