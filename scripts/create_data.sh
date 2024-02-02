#!/bin/bash

source .env

generate_pg_array() {
    arr=("$@")
    pg_array="{"  # 开始一个 PostgreSQL 数组

    for ((i=0; i<${#arr[@]}; i++)); do
        if [ $i -gt 0 ]; then
            pg_array+=","  # 添加逗号分隔符
        fi
        pg_array+="\"${arr[i]}\""  # 添加元素到数组
    done

    pg_array+="}"  # 结束 PostgreSQL 数组
    echo "$pg_array"
}

RANDOM=$$$(date +%s)

for i in {1..10}; do
    title="Ad #$i"
    current_time=$(date -Iseconds)
    start_at=$(date -Iseconds -d "$current_time -$((RANDOM % 30)) days")
    end_at=$(date -Iseconds -d "$current_time +$((RANDOM % 30)) days")
    age_start=$((RANDOM % 80 + 1))
    age_end=$((age_start + RANDOM % (100 - age_start) + 1))

    genders=('M' 'F')
    countries=('US' 'CA' 'JP' 'KR' 'DE' 'FR' 'CN' 'TW' 'GB')
    platforms=('android' 'ios' 'web')

    gender=$(generate_pg_array "${genders[@]}")
    country=$(generate_pg_array "${countries[@]}")
    platform=$(generate_pg_array "${platforms[@]}")

    query="INSERT INTO advertisements (title, start_at, end_at, age_start, age_end, gender, country, platform) VALUES ('$title', '$start_at', '$end_at', $age_start, $age_end, '$gender', '$country', '$platform');"

    PGPASSWORD="$POSTGRES_PASSWORD" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -h "$POSTGRES_HOSTNAME" -c "$query"
done

echo "Done"
