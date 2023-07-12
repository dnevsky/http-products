import psycopg2
from psycopg2 import sql
from psycopg2.extras import execute_values
import random
import time

###
# Скрипт, добавляющий в таблицу products нужное количество записей. Настраивается ниже.
###

conn = psycopg2.connect(
    dbname = 'testdb',
    user = 'postgres',
    password = 'qwerty',
    host = '127.0.0.1',
    port = '5432'
)

cursor = conn.cursor()

total_records = 1_000_000 # общее кол-во записей для вставки

random.seed(time.time())

count = 0

records = []
for _ in range(total_records):

    product = "product" + str(count).zfill(6)
    price = random.randint(0, 100_000)

    records.append((product, price))
    count = count + 1

execute_values(cursor, "INSERT INTO products (id, price) VALUES %s", records)

conn.commit()

print(records[0][0], records[len(records)-1][0], "ok")

cursor.close()
conn.close()