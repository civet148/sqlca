package postgres

/*
-- postgres表名查询
select * from pg_tables where schemaname='public'

-- postgres表结构查询
SELECT a.attnum,
a.attname AS field,
t.typname AS type,
a.attlen AS length,
a.atttypmod AS lengthvar,
a.attnotnull AS notnull,
b.description AS comment
FROM pg_class c,
pg_attribute a
LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,
pg_type t
WHERE c.relname = 'users'
and a.attnum > 0
and a.attrelid = c.oid
and a.atttypid = t.oid
ORDER BY a.attnum;
*/
