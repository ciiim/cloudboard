#!/bin/bash

# 提取 RSS（Resident Set Size）
RSS=$(ps -p 203929 -o rss=)

# 将 RSS 转换为 MB（以 KB 为单位）
RSS_MB=$(echo "scale=2; $RSS / 1024" | bc)

echo "进程占用的内存为：$RSS_MB MB"

