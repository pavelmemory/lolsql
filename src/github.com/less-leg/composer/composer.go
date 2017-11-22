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
   - IsNot - same as `<>` or `!=`
   - Like - same as `LIKE '_va%' AND LIKE '123va%'`
   - LikeOr - same as `LIKE '_va%' OR LIKE '11va%'`
   - NotLike - same as `NOT LIKE '_va%' AND NOT LIKE '123va%'`
   - NotLikeOr - same as `NOT LIKE '_va%' OR NOT LIKE '11va%'`
   - Greater - same as `>`
   - Lower - same as `<`
   - IsOrGreater - same as `>=`
   - IsOrLower - same as `<=`
   - Between - same as `BETWEEN 1 AND 4`
   - IsNull - same as `IS NULL`
   - IsNotNull - same as `IS NOT NULL`
   - In - same as `IN (...)`
   - NotIn - same as `NOT IN (...)`
*/
