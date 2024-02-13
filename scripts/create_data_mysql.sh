#!/bin/bash

source .env

# 生成具有隨機元素的字符串
generate_string_from_array() {
    local all_elements=("$@")  # 接收所有可能的元素
    local selected_elements=()  # 存儲將要選擇的元素
    local num_elements=${#all_elements[@]}  # 元素總數
    local num_select=$((RANDOM % num_elements + 1))  # 隨機決定選取的元素數量

    # 隨機選擇元素
    for ((i=0; i<num_select; i++)); do
        local rand_index=$((RANDOM % num_elements))
        selected_elements+=("${all_elements[rand_index]}")
        # 防止選取重複元素
        all_elements=( "${all_elements[@]:0:rand_index}" "${all_elements[@]:rand_index+1}" )
        ((num_elements--))
    done

    # 構建逗號分隔的字符串
    local result=$(printf "%s," "${selected_elements[@]}")
    result=${result%,}  # 移除末尾的逗號
    echo "$result"
}

RANDOM=$$$(date +%s)

for i in {1..100}; do
    title="Ad #$i"
    # 調整時間格式以符合MySQL標準
    current_time=$(date +%Y-%m-%d\ %H:%M:%S)
    start_at=$(date -v-"$((RANDOM % 30))"d "+%Y-%m-%d %H:%M:%S")
    end_at=$(date -v+"$((RANDOM % 30))"d "+%Y-%m-%d %H:%M:%S")
    age_start=$((RANDOM % 80 + 1))
    age_end=$((age_start + RANDOM % (100 - age_start) + 1))

    genders=('M' 'F')
    countries=('US' 'CA' 'JP' 'KR' 'DE' 'FR' 'CN' 'TW' 'GB')
    platforms=('android' 'ios' 'web')

    # 生成逗號分隔的字符串
    gender=$(generate_string_from_array "${genders[@]}")
    country=$(generate_string_from_array "${countries[@]}")
    platform=$(generate_string_from_array "${platforms[@]}")

    query="INSERT INTO advertisements (title, start_at, end_at, age_start, age_end, gender, country, platform) VALUES ('$title', '$start_at', '$end_at', $age_start, $age_end, '$gender', '$country', '$platform');"

    mysql -u "$MYSQL_USER" --password="$MYSQL_PASSWORD" -h "$MYSQL_HOST" "$MYSQL_DATABASE" -e "$query"
done

echo "Done"
