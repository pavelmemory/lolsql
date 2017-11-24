package composer

/*
 - for each entity there is a same-name function to specify fields to retrieve
 - each field represent a global read-only variable and has a bunch of methods to use in
   - WHERE
   - GROUP BY
   - HAVING
   - ORDER BY
  such as
   - Equal - same as `=`
   - NotEqual - same as `<>` or `!=`
   - Like - same as `LIKE '_va%' AND LIKE '123va%'`
   - LikeOr - same as `LIKE '_va%' OR LIKE '11va%'`
   - NotLike - same as `NOT LIKE '_va%' AND NOT LIKE '123va%'`
   - NotLikeOr - same as `NOT LIKE '_va%' OR NOT LIKE '11va%'`
   - Greater - same as `>`
   - Lower - same as `<`
   - EqualOrGreater - same as `>=`
   - EqualOrLower - same as `<=`
   - Between - same as `BETWEEN 1 AND 4`
   - IsNull - same as `IS NULL`
   - IsNotNull - same as `IS NOT NULL`
   - In - same as `IN (...)`
   - NotIn - same as `NOT IN (...)`
*/

/*
Each entity must have fields that are represent their identity(uniqueness).
Each such field must have comparable by `=` type.
Each field must be marked with tag `lol:"PK"` (case insensitive).

Each entity can have fields with supported types to map to table columns.
By default field name will be used to build SQL queries.
It can be overwritten with tag `lol:"column[PK]"` where 'PK' is a column name that will be used.

By default each entity with tag `lol:"PK"` will be mapped to table in database.
Type name will be used as table name.
It can be overwritten by implementing `TableNameOverride` interface.

Entity can have relations: 1:1, 1:m, m:n
Relation between entities described by providing tag `lol:"FK[Id->Identifier]"`
Where 'FK' signals about existing relation and
join of tables must be done via field 'Id' for current type and field 'Identifier' for type that specified for field with tag.

*/
