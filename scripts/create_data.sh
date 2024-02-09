#!/bin/bash

source .env

# 生成具有随机元素的 PostgreSQL 数组
generate_pg_array() {
    local all_elements=("$@")  # 接收所有可能的元素
    local selected_elements=()  # 存储将要选择的元素
    local num_elements=${#all_elements[@]}  # 元素总数
    local num_select=$((RANDOM % num_elements + 1))  # 随机决定选取的元素数量

    # 随机选择元素
    for ((i=0; i<num_select; i++)); do
        local rand_index=$((RANDOM % num_elements))
        selected_elements+=("${all_elements[rand_index]}")
        # 防止选取重复元素
        all_elements=( "${all_elements[@]:0:rand_index}" "${all_elements[@]:rand_index+1}" )
        ((num_elements--))
    done

    # 构建 PostgreSQL 数组字符串
    local pg_array="{"
    for ((i=0; i<${#selected_elements[@]}; i++)); do
        if [ $i -gt 0 ]; then
            pg_array+=","
        fi
        pg_array+="\"${selected_elements[i]}\""
    done
    pg_array+="}"
    echo "$pg_array"
}

RANDOM=$$$(date +%s)

for i in {1..100}; do
    title="Ad #$i"
    current_time=$(date -Iseconds)
    start_at=$(date -Iseconds -d "$current_time -$((RANDOM % 30)) days")
    end_at=$(date -Iseconds -d "$current_time +$((RANDOM % 30)) days")
    age_start=$((RANDOM % 80 + 1))
    age_end=$((age_start + RANDOM % (100 - age_start) + 1))

    genders=('M' 'F')
    countries=('US' 'CA' 'JP' 'KR' 'DE' 'FR' 'CN' 'TW' 'GB')
    platforms=('android' 'ios' 'web')

    # 从数组中随机选择元素以构建 PostgreSQL 数组
    gender=$(generate_pg_array "${genders[@]}")
    country=$(generate_pg_array "${countries[@]}")
    platform=$(generate_pg_array "${platforms[@]}")

    query="INSERT INTO advertisements (title, start_at, end_at, age_start, age_end, gender, country, platform) VALUES ('$title', '$start_at', '$end_at', $age_start, $age_end, '$gender', '$country', '$platform');"

    PGPASSWORD="$POSTGRES_PASSWORD" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -h "$POSTGRES_HOSTNAME" -c "$query"
done

echo "Done"
