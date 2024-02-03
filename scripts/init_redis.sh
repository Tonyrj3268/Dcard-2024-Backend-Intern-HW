#!/bin/bash

REDIS_HOST=redis
REDIS_PORT=6379

# 獲取當前日期的年份、月份和日
current_year=$(date +%Y)
current_month=$(date +%m)
current_day=$(date +%d)

# 設定當晚12點的時間
midnight_time="00:00:00"

# 構建完整的日期時間字符串
target_datetime="$current_year-$current_month-$current_day $midnight_time"

# 將目標日期時間轉換為時間戳
timestamp=$(date -d "$target_datetime" +%s)

# 使用redis-cli設定key的過期時間
redis-cli -h $REDIS_HOST -p $REDIS_PORT SETEX "CreateAd" $timestamp 0
