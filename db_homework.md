2. В Postgresql при помощи SQL-запроса создайте таблицу developers (запрос добавьте в db_homework.txt)
```sql
`Сначала добавим расширение для поля с uuid4 и координатами`
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; 

CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    department TEXT NOT NULL,
    geolocation GEOGRAPHY(Point, 4326),  -- широта/долгота
    last_known_ip INET,
    is_available BOOLEAN
)
```
3. Добавьте в таблицу 20 записей на основании правил ниже. 
|id                                  |name             |department|geolocation              |last_known_ip |is_available|
|------------------------------------|-----------------|----------|-------------------------|--------------|------------|
|07c11522-6d35-4c34-9c4b-04cf67f1e423|James Smith      |backend   |POINT (37.6173 55.7558)  |192.168.0.1   |true        |
|02a31a34-bd44-4534-864d-dcd7a9901144|Mary Johnson     |frontend  |POINT (-74.006 40.7128)  |10.0.0.2      |false       |
|af150a69-e347-41ed-933e-ca329f678933|John Williams    |ios       |POINT (2.3522 48.8566)   |172.16.5.5    |true        |
|e9b70b5b-8b7b-419a-8295-5e3856361ac1|Patricia Brown   |android   |POINT (139.6917 35.6895) |192.168.1.3   |true        |
|14daba93-0b8e-4382-8c30-6d536e48e39d|Robert Jones     |backend   |POINT (13.405 52.52)     |192.0.2.4     |false       |
|b1a3c838-e21f-47aa-9c38-6720a02ba0ed|James Johnson    |frontend  |POINT (-0.1276 51.5074)  |198.51.100.7  |true        |
|d2455317-645b-4356-8cb1-42bccae52857|Mary Smith       |ios       |POINT (144.9631 -37.8136)|203.0.113.9   |false       |
|e1b1986f-b4bd-4580-b14f-6f8c17bdf75f|John Brown       |android   |POINT (-58.3816 -34.6037)|10.10.10.10   |true        |
|d6067f0b-4051-47ff-a32a-eb3f79b9b60f|Patricia Jones   |backend   |POINT (30.5234 50.4501)  |192.168.100.1 |false       |
|0eeb402a-29bb-4f3b-9944-9d8ee4febbc7|Robert Williams  |frontend  |POINT (-43.1729 -22.9068)|8.8.8.8       |true        |
|9e73a681-cdf9-43f7-b27f-c8221c68583d|James Brown      |ios       |POINT (116.4074 39.9042) |203.0.113.45  |false       |
|0771cbb6-6941-4ab9-8330-765bf8d60418|Mary Jones       |android   |POINT (151.2093 -33.8688)|10.1.2.3      |true        |
|13c9a016-b65a-44d0-a197-071b5b814b1b|John Smith       |backend   |POINT (-3.7038 40.4168)  |172.20.0.5    |true        |
|ae8e1031-a74c-4b46-9214-9822fee5b4d1|Patricia Williams|frontend  |POINT (18.4241 -33.9249) |192.168.2.2   |false       |
|992d3459-ee50-42be-9708-af629448d302|Robert Johnson   |ios       |POINT (55.2708 25.2048)  |203.0.113.77  |true        |
|83d2fce5-bdc1-43f9-8358-c3073b2d3338|James Jones      |android   |POINT (103.8198 1.3521)  |192.168.4.4   |true        |
|049f9906-df7f-45e1-9bb4-c8e557f61b38|Mary Brown       |backend   |POINT (-99.1332 19.4326) |10.0.0.99     |false       |
|67245cc8-5935-4af2-8c57-bc7dfb94dc77|John Johnson     |frontend  |POINT (-70.6693 -33.4489)|192.0.2.123   |true        |
|e282c1f5-e630-4edf-b0d3-f02424ec9939|Patricia Smith   |ios       |POINT (31.2357 30.0444)  |172.16.16.16  |true        |
|fb6f34ae-00ff-40d6-a1a4-2ba484e703f7|Robert Brown     |android   |POINT (106.8456 -6.2088) |198.51.100.100|false       |

3.1 Задача со звёздочкой: 2-db-Homework/1000-devs-add.sql

4.1
Запросы исполняются на 1000 строк, сделанных с помощью скрипта

SELECT * FROM developers WHERE name LIKE 'James%'; 
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=204 width=76)|
|  Filter: (name ~~ 'James%'::text)                          |
4.2
SELECT * FROM developers WHERE department = 'backend';
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=250 width=76)|
|  Filter: (department = 'backend'::text)                    |
4.3
SELECT * FROM developers WHERE last_known_ip = '192.168.1.10';
|QUERY PLAN                                                |
|----------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=1 width=76)|
|  Filter: (last_known_ip = '192.168.1.10'::inet)          |
4.4
SELECT * FROM developers WHERE is_available = TRUE;
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..23.00 rows=494 width=76)|
|  Filter: is_available                                      |
4.5 Задача со звёздочкой
SELECT * 
FROM developers
WHERE ST_DWithin(
    geolocation,
    ST_SetSRID(ST_MakePoint(20.0, 55.0), 4326)::geography,
    10000
);
5. Подумайте, какой индекс можно создать для каждого из полей в таблице, чтобы ускорить поиск.
Вот автономные запросы для создания оптимальных индексов по каждому полю таблицы. Поле `id` уже индексируется автоматически первичным ключом, поэтому дополнительный индекс для него не нужен.

```sql
CREATE INDEX idx_developers_name ON developers (name);

CREATE INDEX idx_developers_department ON developers (department);
-- Применяем gist, т.к. b-tree будет менее эффективен
CREATE INDEX idx_developers_geolocation ON developers USING GIST (geolocation);
-- Применяем gist, т.к. b-tree будет менее эффективен
CREATE INDEX idx_developers_last_known_ip ON developers USING GIST (last_known_ip inet_ops);

CREATE INDEX idx_developers_is_available ON developers (is_available);
```
6. Повторно выполните запросы из пункта 4 
Запросы исполняются на 1000 строк, сделанных с помощью скрипта

6.1
SELECT * FROM developers WHERE name LIKE 'James%'; 
Было:
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=204 width=76)|
|  Filter: (name ~~ 'James%'::text)                          |
Стало:
SELECT * FROM developers WHERE name LIKE 'James%'; 
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=204 width=76)|
|  Filter: (name ~~ 'James%'::text)  
Индекс не используется, полагаю что из-за слишком частого и спонтанного присуствия name='James%'

6.2
SELECT * FROM developers WHERE department = 'backend';
Было:
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=250 width=76)|
|  Filter: (department = 'backend'::text)                    |
Стало:
|QUERY PLAN                                                                              |
|----------------------------------------------------------------------------------------|
|Bitmap Heap Scan on developers  (cost=6.09..22.21 rows=250 width=76)                    |
|  Recheck Cond: (department = 'backend'::text)                                          |
|  ->  Bitmap Index Scan on idx_developers_department  (cost=0.00..6.03 rows=250 width=0)|
|        Index Cond: (department = 'backend'::text)                                      |
Индекс используется, сократил стоимость более чем в 3 раза!
6.3
Было:
SELECT * FROM developers WHERE last_known_ip = '192.168.1.10';
|QUERY PLAN                                                |
|----------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..25.50 rows=1 width=76)|
|  Filter: (last_known_ip = '192.168.1.10'::inet)          |
Стало:
Seq Scan on developers  (cost=0.00..25.50 rows=1 width=76)
  Filter: (last_known_ip = '192.168.1.10'::inet)

Индекс не используется, аналогичная ситуация с обычным b-tree индексом

Seq Scan on developers  (cost=0.00..25.50 rows=1 width=76)
  Filter: (last_known_ip = '192.168.1.10'::inet)

6.4
Было:
SELECT * FROM developers WHERE is_available = TRUE;
|QUERY PLAN                                                  |
|------------------------------------------------------------|
|Seq Scan on developers  (cost=0.00..23.00 rows=494 width=76)|
|  Filter: is_available    
Стало:
Seq Scan on developers  (cost=0.00..23.00 rows=494 width=76)
  Filter: is_available
Индекс тоже не используется, данные не подходят под b-tree индекс

Также пробовал альтернативные по типу CREATE INDEX idx_developers_is_available_true ON developers (is_available) WHERE is_available = TRUE;
Результат такой же негативный.

Тестируем b-tree индекс на поле name
Запрос - SELECT * FROM DEVELOPERS для 10 000 записей: 
Seq Scan on developers  (cost=0.00..230.20 rows=10020 width=76)

Запрос - SELECT * FROM DEVELOPERS для 100 000 записей: 
Seq Scan on developers  (cost=0.00..2296.69 rows=99969 width=76)