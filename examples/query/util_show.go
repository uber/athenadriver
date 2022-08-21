// Copyright (c) 2022 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"database/sql"
	"log"

	secret "github.com/uber/athenadriver/examples/constants"
	drv "github.com/uber/athenadriver/go"
)

func main() {
	// 1. Set AWS Credential in Driver Config.
	conf, err := drv.NewDefaultConfig(secret.OutputBucket, secret.Region,
		secret.AccessID, secret.SecretAccessKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 2. Open Connection.
	dsn := conf.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	for _, q := range []string{
		"SHOW FUNCTIONS",
		"show TBLPROPERTIES sampledb.elb_logs;",
		"SHOW COLUMNS IN sampledb.elb_logs;",
		"SHOW CREATE TABLE sampledb.elb_logs",
		"SHOW CREATE VIEW henrywu",
		"SHOW SCHEMAS",
		"SHOW DATABASES LIKE '[a-z0-9]*db'",
		"SHOW PARTITIONS sampledb.elb_logs",
		"SHOW TABLES;",
		"SHOW VIEWS;",
	} {
		rows, err := db.Query(q)
		if err != nil {
			println(err.Error())
			println("========================================")
			continue
		}
		defer rows.Close()
		println(drv.ColsRowsToCSV(rows))
		println("========================================")
	}
}

/*
Sample Output:
Function,Return Type,Argument Types,Function Type,Deterministic,Description
abs,bigint,bigint,scalar,true,absolute value
abs,decimal(p,s),decimal(p,s),scalar,true,absolute value
abs,double,double,scalar,true,absolute value
abs,integer,integer,scalar,true,absolute value
abs,real,real,scalar,true,absolute value
abs,smallint,smallint,scalar,true,absolute value
abs,tinyint,tinyint,scalar,true,absolute value
acos,double,double,scalar,true,arc cosine
approx_distinct,bigint,bigint,aggregate,true,
approx_distinct,bigint,bigint, double,aggregate,true,
approx_distinct,bigint,double,aggregate,true,
approx_distinct,bigint,double, double,aggregate,true,
approx_distinct,bigint,varbinary,aggregate,true,
approx_distinct,bigint,varbinary, double,aggregate,true,
approx_distinct,bigint,varchar(x),aggregate,true,
approx_distinct,bigint,varchar(x), double,aggregate,true,
approx_percentile,array(bigint),bigint, array(double),aggregate,true,
approx_percentile,array(bigint),bigint, bigint, array(double),aggregate,true,
approx_percentile,array(double),double, array(double),aggregate,true,
approx_percentile,array(double),double, bigint, array(double),aggregate,true,
approx_percentile,array(real),real, array(double),aggregate,true,
approx_percentile,array(real),real, bigint, array(double),aggregate,true,
approx_percentile,bigint,bigint, bigint, double,aggregate,true,
approx_percentile,bigint,bigint, bigint, double, double,aggregate,true,
approx_percentile,bigint,bigint, double,aggregate,true,
approx_percentile,double,double, bigint, double,aggregate,true,
approx_percentile,double,double, bigint, double, double,aggregate,true,
approx_percentile,double,double, double,aggregate,true,
approx_percentile,real,real, bigint, double,aggregate,true,
approx_percentile,real,real, bigint, double, double,aggregate,true,
approx_percentile,real,real, double,aggregate,true,
approx_set,HyperLogLog,bigint,aggregate,true,
approx_set,HyperLogLog,double,aggregate,true,
approx_set,HyperLogLog,varchar(x),aggregate,true,
arbitrary,T,T,aggregate,true,return an arbitrary non-null input value
array_agg,array(T),T,aggregate,true,return an array of values
array_distinct,array(E),array(E),scalar,true,Remove duplicate values from the given array
array_except,array(E),array(E), array(E),scalar,true,Returns an array of elements that are in the first array but not the second, without duplicates.
array_intersect,array(E),array(E), array(E),scalar,true,Intersects elements of the two given arrays
array_join,varchar,array(T), varchar,scalar,true,Concatenates the elements of the given array using a delimiter and an optional string to replace nulls
array_join,varchar,array(T), varchar, varchar,scalar,true,Concatenates the elements of the given array using a delimiter and an optional string to replace nulls
array_max,T,array(T),scalar,true,Get maximum value of array
array_min,T,array(T),scalar,true,Get minimum value of array
array_position,bigint,array(T), T,scalar,true,Returns the position of the first occurrence of the given value in array (or 0 if not found)
array_remove,array(E),array(E), E,scalar,true,Remove specified values from the given array
array_sort,array(E),array(E),scalar,true,Sorts the given array in ascending order according to the natural ordering of its elements.
array_union,array(E),array(E), array(E),scalar,true,Union elements of the two given arrays
arrays_overlap,boolean,array(E), array(E),scalar,true,Returns true if arrays have common elements
asin,double,double,scalar,true,arc sine
atan,double,double,scalar,true,arc tangent
atan2,double,double, double,scalar,true,arc tangent of given fraction
avg,decimal(p,s),decimal(p,s),aggregate,true,Calculates the average value
avg,double,bigint,aggregate,true,
avg,double,double,aggregate,true,
avg,real,real,aggregate,true,
bar,varchar,double, bigint,scalar,true,
bar,varchar,double, bigint, color, color,scalar,true,
bit_count,bigint,bigint, bigint,scalar,true,count number of set bits in 2's complement representation
bitwise_and,bigint,bigint, bigint,scalar,true,bitwise AND in 2's complement arithmetic
bitwise_and_agg,bigint,bigint,aggregate,true,
bitwise_not,bigint,bigint,scalar,true,bitwise NOT in 2's complement arithmetic
bitwise_or,bigint,bigint, bigint,scalar,true,bitwise OR in 2's complement arithmetic
bitwise_or_agg,bigint,bigint,aggregate,true,
bitwise_xor,bigint,bigint, bigint,scalar,true,bitwise XOR in 2's complement arithmetic
bool_and,boolean,boolean,aggregate,true,
bool_or,boolean,boolean,aggregate,true,
build_geo_index,varchar,varchar, varbinary,aggregate,true,
build_geo_index,varchar,varchar, varchar,aggregate,true,
cardinality,bigint,HyperLogLog,scalar,true,compute the cardinality of a HyperLogLog instance
cardinality,bigint,array(E),scalar,true,Returns the cardinality (length) of the array
cardinality,bigint,map(K,V),scalar,true,Returns the cardinality (the number of key-value pairs) of the map
cbrt,double,double,scalar,true,cube root
ceil,bigint,bigint,scalar,true,round up to nearest integer
ceil,decimal(rp,0),decimal(p,s),scalar,true,round up to nearest integer
ceil,double,double,scalar,true,round up to nearest integer
ceil,integer,integer,scalar,true,round up to nearest integer
ceil,real,real,scalar,true,round up to nearest integer
ceil,smallint,smallint,scalar,true,round up to nearest integer
ceil,tinyint,tinyint,scalar,true,round up to nearest integer
ceiling,bigint,bigint,scalar,true,round up to nearest integer
ceiling,decimal(rp,0),decimal(p,s),scalar,true,round up to nearest integer
ceiling,double,double,scalar,true,round up to nearest integer
ceiling,integer,integer,scalar,true,round up to nearest integer
ceiling,real,real,scalar,true,round up to nearest integer
ceiling,smallint,smallint,scalar,true,round up to nearest integer
ceiling,tinyint,tinyint,scalar,true,round up to nearest integer
char2hexint,varchar,varchar,scalar,true,Returns the hexadecimal representation of the UTF-16BE encoding of the argument
checksum,varbinary,T,aggregate,true,Checksum of the given values
chr,varchar(1),bigint,scalar,true,convert Unicode code point to a string
classify,bigint,map(bigint,double), Classifier(bigint),scalar,true,
classify,varchar,map(bigint,double), Classifier(varchar),scalar,true,
codepoint,integer,varchar(1),scalar,true,returns Unicode code point of a single character string
color,color,double, color, color,scalar,true,
color,color,double, double, double, color, color,scalar,true,
color,color,varchar(x),scalar,true,
concat,array(E),E, array(E),scalar,true,Concatenates an element to an array
concat,array(E),array(E),scalar,true,Concatenates given arrays
concat,array(E),array(E), E,scalar,true,Concatenates an array to an element
concat,varchar,varchar,scalar,true,concatenates given strings
contains,boolean,array(T), T,scalar,true,Determines whether given value exists in the array
corr,double,double, double,aggregate,true,
corr,real,real, real,aggregate,true,
cos,double,double,scalar,true,cosine
cosh,double,double,scalar,true,hyperbolic cosine
cosine_similarity,double,map(varchar,double), map(varchar,double),scalar,true,cosine similarity between the given sparse vectors
count,bigint,,aggregate,true,
count,bigint,T,aggregate,true,Counts the non-null values
count_if,bigint,boolean,aggregate,true,
covar_pop,double,double, double,aggregate,true,
covar_pop,real,real, real,aggregate,true,
covar_samp,double,double, double,aggregate,true,
covar_samp,real,real, real,aggregate,true,
cume_dist,double,,window,true,
current_date,date,,scalar,true,current date
current_time,time with time zone,,scalar,true,current time with time zone
current_timestamp,timestamp with time zone,,scalar,true,current timestamp with time zone
current_timezone,varchar,,scalar,true,current time zone
date,date,timestamp,scalar,true,
date,date,timestamp with time zone,scalar,true,
date,date,varchar(x),scalar,true,
date_add,date,varchar(x), bigint, date,scalar,true,add the specified amount of date to the given date
date_add,time,varchar(x), bigint, time,scalar,true,add the specified amount of time to the given time
date_add,time with time zone,varchar(x), bigint, time with time zone,scalar,true,add the specified amount of time to the given time
date_add,timestamp,varchar(x), bigint, timestamp,scalar,true,add the specified amount of time to the given timestamp
date_add,timestamp with time zone,varchar(x), bigint, timestamp with time zone,scalar,true,add the specified amount of time to the given timestamp
date_diff,bigint,varchar(x), date, date,scalar,true,difference of the given dates in the given unit
date_diff,bigint,varchar(x), time with time zone, time with time zone,scalar,true,difference of the given times in the given unit
date_diff,bigint,varchar(x), time, time,scalar,true,difference of the given times in the given unit
date_diff,bigint,varchar(x), timestamp with time zone, timestamp with time zone,scalar,true,difference of the given times in the given unit
date_diff,bigint,varchar(x), timestamp, timestamp,scalar,true,difference of the given times in the given unit
date_format,varchar,timestamp with time zone, varchar(x),scalar,true,
date_format,varchar,timestamp, varchar(x),scalar,true,
date_parse,timestamp,varchar(x), varchar(y),scalar,true,
date_trunc,date,varchar(x), date,scalar,true,truncate to the specified precision in the session timezone
date_trunc,time,varchar(x), time,scalar,true,truncate to the specified precision in the session timezone
date_trunc,time with time zone,varchar(x), time with time zone,scalar,true,truncate to the specified precision
date_trunc,timestamp,varchar(x), timestamp,scalar,true,truncate to the specified precision in the session timezone
date_trunc,timestamp with time zone,varchar(x), timestamp with time zone,scalar,true,truncate to the specified precision
day,bigint,date,scalar,true,day of the month of the given date
day,bigint,interval day to second,scalar,true,day of the month of the given interval
day,bigint,timestamp,scalar,true,day of the month of the given timestamp
day,bigint,timestamp with time zone,scalar,true,day of the month of the given timestamp
day_of_month,bigint,date,scalar,true,day of the month of the given date
day_of_month,bigint,interval day to second,scalar,true,day of the month of the given interval
day_of_month,bigint,timestamp,scalar,true,day of the month of the given timestamp
day_of_month,bigint,timestamp with time zone,scalar,true,day of the month of the given timestamp
day_of_week,bigint,date,scalar,true,day of the week of the given date
day_of_week,bigint,timestamp,scalar,true,day of the week of the given timestamp
day_of_week,bigint,timestamp with time zone,scalar,true,day of the week of the given timestamp
day_of_year,bigint,date,scalar,true,day of the year of the given date
day_of_year,bigint,timestamp,scalar,true,day of the year of the given timestamp
day_of_year,bigint,timestamp with time zone,scalar,true,day of the year of the given timestamp
degrees,double,double,scalar,true,converts an angle in radians to degrees
dense_rank,bigint,,window,true,
dow,bigint,date,scalar,true,day of the week of the given date
dow,bigint,timestamp,scalar,true,day of the week of the given timestamp
dow,bigint,timestamp with time zone,scalar,true,day of the week of the given timestamp
doy,bigint,date,scalar,true,day of the year of the given date
doy,bigint,timestamp,scalar,true,day of the year of the given timestamp
doy,bigint,timestamp with time zone,scalar,true,day of the year of the given timestamp
e,double,,scalar,true,Euler's number
element_at,E,array(E), bigint,scalar,true,Get element of array at given index
element_at,V,map(K,V), K,scalar,true,Get value for the given key, or null if it does not exist
empty_approx_set,HyperLogLog,,scalar,true,an empty HyperLogLog instance
evaluate_classifier_predictions,varchar,bigint, bigint,aggregate,true,
evaluate_classifier_predictions,varchar,varchar(x), varchar(y),aggregate,true,
every,boolean,boolean,aggregate,true,
exp,double,double,scalar,true,Euler's number raised to the given power
features,map(bigint,double),double,scalar,true,
features,map(bigint,double),double, double,scalar,true,
features,map(bigint,double),double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double, double, double, double, double,scalar,true,
features,map(bigint,double),double, double, double, double, double, double, double, double, double, double,scalar,true,
filter,array(T),array(T), function(T,boolean),scalar,false,return array containing elements that match the given predicate
first_value,T,T,window,true,
flatten,array(E),array(array(E)),scalar,true,Flattens the given array
floor,bigint,bigint,scalar,true,round down to nearest integer
floor,decimal(rp,0),decimal(p,s),scalar,true,round down to nearest integer
floor,double,double,scalar,true,round down to nearest integer
floor,integer,integer,scalar,true,round down to nearest integer
floor,real,real,scalar,true,round down to nearest integer
floor,smallint,smallint,scalar,true,round down to nearest integer
floor,tinyint,tinyint,scalar,true,round down to nearest integer
format_datetime,varchar,timestamp with time zone, varchar(x),scalar,true,formats the given time by the given format
format_datetime,varchar,timestamp, varchar(x),scalar,true,formats the given time by the given format
from_base,bigint,varchar(x), bigint,scalar,true,convert a string in the given base to a number
from_base64,varbinary,varbinary,scalar,true,decode base64 encoded binary data
from_base64,varbinary,varchar(x),scalar,true,decode base64 encoded binary data
from_base64url,varbinary,varbinary,scalar,true,decode URL safe base64 encoded binary data
from_base64url,varbinary,varchar(x),scalar,true,decode URL safe base64 encoded binary data
from_big_endian_64,bigint,varbinary,scalar,true,decode bigint value from a 64-bit 2's complement big endian varbinary
from_hex,varbinary,varbinary,scalar,true,decode hex encoded binary data
from_hex,varbinary,varchar(x),scalar,true,decode hex encoded binary data
from_iso8601_date,date,varchar(x),scalar,true,
from_iso8601_timestamp,timestamp with time zone,varchar(x),scalar,true,
from_unixtime,timestamp,double,scalar,true,
from_unixtime,timestamp with time zone,double, bigint, bigint,scalar,true,
from_unixtime,timestamp with time zone,double, varchar(x),scalar,true,
from_utf8,varchar,varbinary,scalar,true,decodes the UTF-8 encoded string
from_utf8,varchar,varbinary, bigint,scalar,true,decodes the UTF-8 encoded string
from_utf8,varchar,varbinary, varchar(x),scalar,true,decodes the UTF-8 encoded string
geo_contains,varchar,varbinary, varchar,scalar,true,Returns the first key whose corresponding geometry contains geo shape
geo_contains,varchar,varchar, varchar,scalar,true,Returns the first key whose corresponding geometry contains geo shape
geo_contains_all,array(varchar),varbinary, varchar,scalar,true,Returns all keys whose corresponding geometry contains geo shape
geo_contains_all,array(varchar),varchar, varchar,scalar,true,Returns all keys whose corresponding geometry contains geo shape
geo_intersects,varchar,varbinary, varchar,scalar,true,Returns the first key whose corresponding geometry intersects with geo shape
geo_intersects,varchar,varchar, varchar,scalar,true,Returns the first key whose corresponding geometry intersects with geo shape
geo_intersects_all,array(varchar),varbinary, varchar,scalar,true,Returns all keys whose corresponding geometry intersects with geo shape
geo_intersects_all,array(varchar),varchar, varchar,scalar,true,Returns all keys whose corresponding geometry intersects with geo shape
geometric_mean,double,bigint,aggregate,true,
geometric_mean,double,double,aggregate,true,
geometric_mean,real,real,aggregate,true,
greatest,E,E,scalar,true,get the largest of the given values
histogram,map(K,bigint),K,aggregate,true,Count the number of times each value occurs
hour,bigint,interval day to second,scalar,true,hour of the day of the given interval
hour,bigint,time,scalar,true,hour of the day of the given time
hour,bigint,time with time zone,scalar,true,hour of the day of the given time
hour,bigint,timestamp,scalar,true,hour of the day of the given timestamp
hour,bigint,timestamp with time zone,scalar,true,hour of the day of the given timestamp
index,bigint,varchar, varchar,scalar,true,Returns index of first occurrence of a substring (or 0 if not found)
infinity,double,,scalar,true,Infinity
is_finite,boolean,double,scalar,true,test if value is finite
is_infinite,boolean,double,scalar,true,test if value is infinite
is_nan,boolean,double,scalar,true,test if value is not-a-number
json_array_contains,boolean,json, bigint,scalar,true,
json_array_contains,boolean,json, boolean,scalar,true,
json_array_contains,boolean,json, double,scalar,true,
json_array_contains,boolean,json, varchar(x),scalar,true,
json_array_contains,boolean,varchar(x), bigint,scalar,true,
json_array_contains,boolean,varchar(x), boolean,scalar,true,
json_array_contains,boolean,varchar(x), double,scalar,true,
json_array_contains,boolean,varchar(x), varchar(y),scalar,true,
json_array_get,json,json, bigint,scalar,true,
json_array_get,json,varchar(x), bigint,scalar,true,
json_array_length,bigint,json,scalar,true,
json_array_length,bigint,varchar(x),scalar,true,
json_extract,json,json, JsonPath,scalar,true,
json_extract,json,varchar(x), JsonPath,scalar,true,
json_extract_scalar,varchar,json, JsonPath,scalar,true,
json_extract_scalar,varchar(x),varchar(x), JsonPath,scalar,true,
json_format,varchar,json,scalar,true,
json_parse,json,varchar(x),scalar,true,
json_size,bigint,json, JsonPath,scalar,true,
json_size,bigint,varchar(x), JsonPath,scalar,true,
kurtosis,double,bigint,aggregate,true,Returns the (excess) kurtosis of the argument
kurtosis,double,double,aggregate,true,Returns the (excess) kurtosis of the argument
lag,T,T,window,true,
lag,T,T, bigint,window,true,
lag,T,T, bigint, T,window,true,
last_value,T,T,window,true,
lead,T,T,window,true,
lead,T,T, bigint,window,true,
lead,T,T, bigint, T,window,true,
learn_classifier,Classifier(bigint),bigint, map(bigint,double),aggregate,true,
learn_classifier,Classifier(bigint),double, map(bigint,double),aggregate,true,
learn_classifier,Classifier(varchar),varchar, map(bigint,double),aggregate,true,
learn_libsvm_classifier,Classifier(bigint),bigint, map(bigint,double), varchar(x),aggregate,true,
learn_libsvm_classifier,Classifier(bigint),double, map(bigint,double), varchar,aggregate,true,
learn_libsvm_classifier,Classifier(varchar),varchar, map(bigint,double), varchar,aggregate,true,
learn_libsvm_regressor,Regressor,bigint, map(bigint,double), varchar,aggregate,true,
learn_libsvm_regressor,Regressor,double, map(bigint,double), varchar,aggregate,true,
learn_regressor,Regressor,bigint, map(bigint,double),aggregate,true,
learn_regressor,Regressor,double, map(bigint,double),aggregate,true,
least,E,E,scalar,true,get the smallest of the given values
length,bigint,char(x),scalar,true,count of code points of the given string
length,bigint,varbinary,scalar,true,length of the given binary
length,bigint,varchar(x),scalar,true,count of code points of the given string
levenshtein_distance,bigint,varchar(x), varchar(y),scalar,true,computes Levenshtein distance between two strings
like_pattern,LikePattern,varchar(x), varchar(y),scalar,true,
ln,double,double,scalar,true,natural logarithm
localtime,time,,scalar,true,current time without time zone
localtimestamp,timestamp,,scalar,true,current timestamp without time zone
log,double,double, double,scalar,true,logarithm to given base
log10,double,double,scalar,true,logarithm to base 10
log2,double,double,scalar,true,logarithm to base 2
lower,char(x),char(x),scalar,true,converts the string to lower case
lower,varchar(x),varchar(x),scalar,true,converts the string to lower case
lpad,varchar,varchar(x), bigint, varchar(y),scalar,true,pads a string on the left
ltrim,char(x),char(x),scalar,true,removes whitespace from the beginning of a string
ltrim,char(x),char(x), CodePoints,scalar,true,remove the longest string containing only given characters from the beginning of a string
ltrim,varchar(x),varchar(x),scalar,true,removes whitespace from the beginning of a string
ltrim,varchar(x),varchar(x), CodePoints,scalar,true,remove the longest string containing only given characters from the beginning of a string
map,map(K,V),array(K), array(V),scalar,true,Constructs a map from the given key/value arrays
map,map(unknown,unknown),,scalar,true,Creates an empty map
map_agg,map(K,V),K, V,aggregate,true,Aggregates all the rows (key/value pairs) into a single map
map_concat,map(K,V),map(K,V),scalar,true,Concatenates given maps
map_filter,map(K,V),map(K,V), function(K,V,boolean),scalar,false,return map containing entries that match the given predicate
map_keys,array(K),map(K,V),scalar,true,Returns the keys of the given map(K,V) as an array
map_union,map(K,V),map(K,V),aggregate,true,Aggregate all the maps into a single map
map_values,array(V),map(K,V),scalar,true,Returns the values of the given map(K,V) as an array
max,E,E,aggregate,true,Returns the maximum value of the argument
max,array(E),E, bigint,aggregate,true,Returns the maximum values of the argument
max_by,V,V, K,aggregate,true,Returns the value of the first argument, associated with the maximum value of the second argument
max_by,array(V),V, K, bigint,aggregate,true,Returns the values of the first argument associated with the maximum values of the second argument
md5,varbinary,varbinary,scalar,true,compute md5 hash
merge,HyperLogLog,HyperLogLog,aggregate,true,
min,E,E,aggregate,true,Returns the minimum value of the argument
min,array(E),E, bigint,aggregate,true,Returns the minimum values of the argument
min_by,V,V, K,aggregate,true,Returns the value of the first argument, associated with the minimum value of the second argument
min_by,array(V),V, K, bigint,aggregate,true,Returns the values of the first argument associated with the minimum values of the second argument
minute,bigint,interval day to second,scalar,true,minute of the hour of the given interval
minute,bigint,time,scalar,true,minute of the hour of the given time
minute,bigint,time with time zone,scalar,true,minute of the hour of the given time
minute,bigint,timestamp,scalar,true,minute of the hour of the given timestamp
minute,bigint,timestamp with time zone,scalar,true,minute of the hour of the given timestamp
mod,bigint,bigint, bigint,scalar,true,remainder of given quotient
mod,decimal(r_precision,r_scale),decimal(a_precision,a_scale), decimal(b_precision,b_scale),scalar,false,
mod,double,double, double,scalar,true,remainder of given quotient
mod,integer,integer, integer,scalar,true,remainder of given quotient
mod,real,real, real,scalar,true,remainder of given quotient
mod,smallint,smallint, smallint,scalar,true,remainder of given quotient
mod,tinyint,tinyint, tinyint,scalar,true,remainder of given quotient
month,bigint,date,scalar,true,month of the year of the given date
month,bigint,interval year to month,scalar,true,month of the year of the given interval
month,bigint,timestamp,scalar,true,month of the year of the given timestamp
month,bigint,timestamp with time zone,scalar,true,month of the year of the given timestamp
multimap_agg,map(K,array(V)),K, V,aggregate,true,Aggregates all the rows (key/value pairs) into a single multimap
nan,double,,scalar,true,constant representing not-a-number
normalize,varchar,varchar(x), varchar(y),scalar,true,transforms the string to normalized form
now,timestamp with time zone,,scalar,true,current timestamp with time zone
nth_value,T,T, bigint,window,true,
ntile,bigint,bigint,window,true,
numeric_histogram,map(double,double),bigint, double,aggregate,true,
numeric_histogram,map(double,double),bigint, double, double,aggregate,true,
numeric_histogram,map(real,real),bigint, real,aggregate,true,
numeric_histogram,map(real,real),bigint, real, double,aggregate,true,
objectid,ObjectId,,scalar,true,mongodb ObjectId
objectid,ObjectId,varchar,scalar,true,mongodb ObjectId from the given string
parse_datetime,timestamp with time zone,varchar(x), varchar(y),scalar,true,parses the specified date/time by the given format
percent_rank,double,,window,true,
pi,double,,scalar,true,the constant Pi
pow,double,double, double,scalar,true,value raised to the power of exponent
power,double,double, double,scalar,true,value raised to the power of exponent
quarter,bigint,date,scalar,true,quarter of the year of the given date
quarter,bigint,timestamp,scalar,true,quarter of the year of the given timestamp
quarter,bigint,timestamp with time zone,scalar,true,quarter of the year of the given timestamp
radians,double,double,scalar,true,converts an angle in degrees to radians
rand,bigint,bigint,scalar,false,a pseudo-random number between 0 and value (exclusive)
rand,double,,scalar,false,a pseudo-random value
rand,integer,integer,scalar,false,a pseudo-random number between 0 and value (exclusive)
rand,smallint,smallint,scalar,false,a pseudo-random number between 0 and value (exclusive)
rand,tinyint,tinyint,scalar,false,a pseudo-random number between 0 and value (exclusive)
random,bigint,bigint,scalar,false,a pseudo-random number between 0 and value (exclusive)
random,double,,scalar,false,a pseudo-random value
random,integer,integer,scalar,false,a pseudo-random number between 0 and value (exclusive)
random,smallint,smallint,scalar,false,a pseudo-random number between 0 and value (exclusive)
random,tinyint,tinyint,scalar,false,a pseudo-random number between 0 and value (exclusive)
rank,bigint,,window,true,
reduce,R,array(T), S, function(S,T,S), function(S,R),scalar,false,Reduce elements of the array into a single value
regexp_extract,varchar(x),varchar(x), JoniRegExp,scalar,true,string extracted using the given pattern
regexp_extract,varchar(x),varchar(x), JoniRegExp, bigint,scalar,true,returns regex group of extracted string with a pattern
regexp_extract_all,array(varchar(x)),varchar(x), JoniRegExp,scalar,true,string(s) extracted using the given pattern
regexp_extract_all,array(varchar(x)),varchar(x), JoniRegExp, bigint,scalar,true,group(s) extracted using the given pattern
regexp_like,boolean,varchar(x), JoniRegExp,scalar,true,returns whether the pattern is contained within the string
regexp_replace,varchar(x),varchar(x), JoniRegExp,scalar,true,removes substrings matching a regular expression
regexp_replace,varchar(z),varchar(x), JoniRegExp, varchar(y),scalar,true,replaces substrings matching a regular expression by given string
regexp_split,array(varchar(x)),varchar(x), JoniRegExp,scalar,true,returns array of strings split by pattern
regr_intercept,double,double, double,aggregate,true,
regr_intercept,real,real, real,aggregate,true,
regr_slope,double,double, double,aggregate,true,
regr_slope,real,real, real,aggregate,true,
regress,double,map(bigint,double), Regressor,scalar,true,
render,varchar(16),boolean,scalar,true,
render,varchar(35),bigint, color,scalar,true,
render,varchar(41),double, color,scalar,true,
render,varchar(y),varchar(x), color,scalar,true,
replace,varchar(u),varchar(x), varchar(y), varchar(z),scalar,true,greedily replaces occurrences of a pattern with a string
replace,varchar(x),varchar(x), varchar(y),scalar,true,greedily removes occurrences of a pattern in a string
reverse,array(E),array(E),scalar,true,Returns an array which has the reversed order of the given array.
reverse,varchar(x),varchar(x),scalar,true,reverse all code points in a given string
rgb,color,bigint, bigint, bigint,scalar,true,
round,bigint,bigint,scalar,true,round to nearest integer
round,bigint,bigint, bigint,scalar,true,round to nearest integer
round,decimal(rp,rs),decimal(p,s),scalar,true,round to nearest integer
round,decimal(rp,s),decimal(p,s), bigint,scalar,true,round to given number of decimal places
round,double,double,scalar,true,round to nearest integer
round,double,double, bigint,scalar,true,round to given number of decimal places
round,integer,integer,scalar,true,round to nearest integer
round,integer,integer, bigint,scalar,true,round to nearest integer
round,real,real,scalar,true,round to given number of decimal places
round,real,real, bigint,scalar,true,round to given number of decimal places
round,smallint,smallint,scalar,true,round to nearest integer
round,smallint,smallint, bigint,scalar,true,round to nearest integer
round,tinyint,tinyint,scalar,true,round to nearest integer
round,tinyint,tinyint, bigint,scalar,true,round to nearest integer
row_number,bigint,,window,true,
rpad,varchar,varchar(x), bigint, varchar(y),scalar,true,pads a string on the right
rtrim,char(x),char(x),scalar,true,removes whitespace from the end of a string
rtrim,char(x),char(x), CodePoints,scalar,true,remove the longest string containing only given characters from the end of a string
rtrim,varchar(x),varchar(x),scalar,true,removes whitespace from the end of a string
rtrim,varchar(x),varchar(x), CodePoints,scalar,true,remove the longest string containing only given characters from the end of a string
second,bigint,interval day to second,scalar,true,second of the minute of the given interval
second,bigint,time,scalar,true,second of the minute of the given time
second,bigint,time with time zone,scalar,true,second of the minute of the given time
second,bigint,timestamp,scalar,true,second of the minute of the given timestamp
second,bigint,timestamp with time zone,scalar,true,second of the minute of the given timestamp
sequence,array(bigint),bigint, bigint,scalar,true,
sequence,array(bigint),bigint, bigint, bigint,scalar,true,Sequence function to generate synthetic arrays
sequence,array(timestamp),timestamp, timestamp, interval day to second,scalar,true,
sequence,array(timestamp),timestamp, timestamp, interval year to month,scalar,true,
sha1,varbinary,varbinary,scalar,true,compute sha1 hash
sha256,varbinary,varbinary,scalar,true,compute sha256 hash
sha512,varbinary,varbinary,scalar,true,compute sha512 hash
shuffle,array(E),array(E),scalar,false,Generates a random permutation of the given array.
sign,bigint,bigint,scalar,true,
sign,decimal(1,0),decimal(p,s),scalar,true,signum
sign,double,double,scalar,true,signum
sign,integer,integer,scalar,true,signum
sign,real,real,scalar,true,signum
sign,smallint,smallint,scalar,true,signum
sign,tinyint,tinyint,scalar,true,signum
sin,double,double,scalar,true,sine
skewness,double,bigint,aggregate,true,Returns the skewness of the argument
skewness,double,double,aggregate,true,Returns the skewness of the argument
slice,array(E),array(E), bigint, bigint,scalar,true,Subsets an array given an offset (1-indexed) and length
split,array(varchar(x)),varchar(x), varchar(y),scalar,true,
split,array(varchar(x)),varchar(x), varchar(y), bigint,scalar,true,
split_part,varchar(x),varchar(x), varchar(y), bigint,scalar,true,splits a string by a delimiter and returns the specified field (counting from one)
split_to_map,map(varchar,varchar),varchar, varchar, varchar,scalar,true,creates a map using entryDelimiter and keyValueDelimiter
sqrt,double,double,scalar,true,square root
st_area,double,varbinary,scalar,true,Returns area of input polygon
st_area,double,varchar,scalar,true,Returns area of input polygon
st_boundary,varbinary,varbinary,scalar,true,Returns string representation of the boundary geometry of input geometry
st_boundary,varbinary,varchar,scalar,true,Returns string representation of the boundary geometry of input geometry
st_buffer,varbinary,varbinary, double,scalar,true,Returns string representation of the geometry buffered by distance
st_buffer,varbinary,varchar, double,scalar,true,Returns string representation of the geometry buffered by distance
st_centroid,varbinary,varbinary,scalar,true,Returns point that is the center of the polygon's envelope
st_centroid,varbinary,varchar,scalar,true,Returns point that is the center of the polygon's envelope
st_contains,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry contains right geometry
st_contains,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry contains right geometry
st_contains,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry contains right geometry
st_contains,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry contains right geometry
st_coordinate_dimension,bigint,varbinary,scalar,true,Returns count of coordinate components
st_coordinate_dimension,bigint,varchar,scalar,true,Returns count of coordinate components
st_crosses,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry crosses right geometry
st_crosses,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry crosses right geometry
st_crosses,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry crosses right geometry
st_crosses,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry crosses right geometry
st_difference,varbinary,varbinary, varbinary,scalar,true,Returns string representation of the geometry difference of left geometry and right geometry
st_difference,varbinary,varbinary, varchar,scalar,true,Returns string representation of the geometry difference of left geometry and right geometry
st_difference,varbinary,varchar, varbinary,scalar,true,Returns string representation of the geometry difference of left geometry and right geometry
st_difference,varbinary,varchar, varchar,scalar,true,Returns string representation of the geometry difference of left geometry and right geometry
st_dimension,bigint,varbinary,scalar,true,Returns spatial dimension of geometry
st_dimension,bigint,varchar,scalar,true,Returns spatial dimension of geometry
st_disjoint,boolean,varbinary, varbinary,scalar,true,Returns true if and only if the intersection of left geometry and right geometry is empty
st_disjoint,boolean,varbinary, varchar,scalar,true,Returns true if and only if the intersection of left geometry and right geometry is empty
st_disjoint,boolean,varchar, varbinary,scalar,true,Returns true if and only if the intersection of left geometry and right geometry is empty
st_disjoint,boolean,varchar, varchar,scalar,true,Returns true if and only if the intersection of left geometry and right geometry is empty
st_distance,double,varbinary, varbinary,scalar,true,Returns distance between left geometry and right geometry
st_distance,double,varbinary, varchar,scalar,true,Returns distance between left geometry and right geometry
st_distance,double,varchar, varbinary,scalar,true,Returns distance between left geometry and right geometry
st_distance,double,varchar, varchar,scalar,true,Returns distance between left geometry and right geometry
st_end_point,varbinary,varbinary,scalar,true,Returns the last point of an line
st_end_point,varbinary,varchar,scalar,true,Returns the last point of an line
st_envelope,varbinary,varbinary,scalar,true,Returns string representation of envelope of the input geometry
st_envelope,varbinary,varchar,scalar,true,Returns string representation of envelope of the input geometry
st_envelope_intersect,boolean,varbinary, varbinary,scalar,true,Returns true if and only if the envelopes of left geometry and right geometry intersect
st_envelope_intersect,boolean,varbinary, varchar,scalar,true,Returns true if and only if the envelopes of left geometry and right geometry intersect
st_envelope_intersect,boolean,varchar, varbinary,scalar,true,Returns true if and only if the envelopes of left geometry and right geometry intersect
st_envelope_intersect,boolean,varchar, varchar,scalar,true,Returns true if and only if the envelopes of left geometry and right geometry intersect
st_equals,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry equals right geometry
st_equals,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry equals right geometry
st_equals,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry equals right geometry
st_equals,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry equals right geometry
st_exterior_ring,varbinary,varbinary,scalar,true,Returns string representation of the exterior ring of the polygon
st_exterior_ring,varbinary,varchar,scalar,true,Returns string representation of the exterior ring of the polygon
st_geohash,varchar,double, double,scalar,true,Returns geo hash of a point
st_geohash,varchar,double, double, integer,scalar,true,Returns geo hash of a point
st_geohash,varchar,varbinary,scalar,true,Returns geo hash of a point
st_geohash,varchar,varbinary, bigint,scalar,true,Returns geo hash of a point
st_geometry_from_text,varbinary,varchar,scalar,true,Returns binary of geometry
st_geometry_to_text,varchar,varbinary,scalar,true,Returns text of geometry
st_interior_ring_number,bigint,varbinary,scalar,true,Returns the number of interior rings in the polygon
st_interior_ring_number,bigint,varchar,scalar,true,Returns the number of interior rings in the polygon
st_intersection,varbinary,varbinary, varbinary,scalar,true,Returns string representation of the geometry intersection of left geometry and right geometry
st_intersection,varbinary,varbinary, varchar,scalar,true,Returns string representation of the geometry intersection of left geometry and right geometry
st_intersection,varbinary,varchar, varbinary,scalar,true,Returns string representation of the geometry intersection of left geometry and right geometry
st_intersection,varbinary,varchar, varchar,scalar,true,Returns string representation of the geometry intersection of left geometry and right geometry
st_intersects,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry intersects right geometry
st_intersects,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry intersects right geometry
st_intersects,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry intersects right geometry
st_intersects,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry intersects right geometry
st_is_closed,boolean,varbinary,scalar,true,Returns true if and only if the line is closed
st_is_closed,boolean,varchar,scalar,true,Returns true if and only if the line is closed
st_is_empty,boolean,varbinary,scalar,true,Returns true if and only if the geometry is empty
st_is_empty,boolean,varchar,scalar,true,Returns true if and only if the geometry is empty
st_is_ring,boolean,varbinary,scalar,true,Returns true if and only if the line is closed and simple
st_is_ring,boolean,varchar,scalar,true,Returns true if and only if the line is closed and simple
st_length,double,varbinary,scalar,true,Returns the length of line
st_length,double,varchar,scalar,true,Returns the length of line
st_line,varbinary,varchar,scalar,true,Returns binary representation of a Line
st_max_x,double,varbinary,scalar,true,Returns the maximum X coordinate of geometry
st_max_x,double,varchar,scalar,true,Returns the maximum X coordinate of geometry
st_max_y,double,varbinary,scalar,true,Returns the maximum Y coordinate of geometry
st_max_y,double,varchar,scalar,true,Returns the maximum Y coordinate of geometry
st_min_X,double,varbinary,scalar,true,Returns the minimum X coordinate of geometry
st_min_X,double,varchar,scalar,true,Returns the minimum X coordinate of geometry
st_min_y,double,varbinary,scalar,true,Returns the minimum Y coordinate of geometry
st_min_y,double,varchar,scalar,true,Returns the minimum Y coordinate of geometry
st_overlaps,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry overlaps right geometry
st_overlaps,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry overlaps right geometry
st_overlaps,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry overlaps right geometry
st_overlaps,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry overlaps right geometry
st_point,varbinary,double, double,scalar,true,Returns binary representation of a Point
st_point_number,bigint,varbinary,scalar,true,Returns the number of points in the geometry
st_point_number,bigint,varchar,scalar,true,Returns the number of points in the geometry
st_polygon,varbinary,varchar,scalar,true,Returns binary representation of a Polygon
st_relate,boolean,varbinary, varbinary, varchar,scalar,true,Returns true if and only if left geometry has the specified DE-9IM relationship with right geometry
st_relate,boolean,varbinary, varchar, varchar,scalar,true,Returns true if and only if left geometry has the specified DE-9IM relationship with right geometry
st_relate,boolean,varchar, varbinary, varchar,scalar,true,Returns true if and only if left geometry has the specified DE-9IM relationship with right geometry
st_relate,boolean,varchar, varchar, varchar,scalar,true,Returns true if and only if left geometry has the specified DE-9IM relationship with right geometry
st_start_point,varbinary,varbinary,scalar,true,Returns the first point of an line
st_start_point,varbinary,varchar,scalar,true,Returns the first point of an line
st_symmetric_difference,varbinary,varbinary, varbinary,scalar,true,Returns string representation of the geometry symmetric difference of left geometry and right geometry
st_symmetric_difference,varbinary,varbinary, varchar,scalar,true,Returns string representation of the geometry symmetric difference of left geometry and right geometry
st_symmetric_difference,varbinary,varchar, varbinary,scalar,true,Returns string representation of the geometry symmetric difference of left geometry and right geometry
st_symmetric_difference,varbinary,varchar, varchar,scalar,true,Returns string representation of the geometry symmetric difference of left geometry and right geometry
st_touches,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry touches right geometry
st_touches,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry touches right geometry
st_touches,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry touches right geometry
st_touches,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry touches right geometry
st_within,boolean,varbinary, varbinary,scalar,true,Returns true if and only if left geometry is within right geometry
st_within,boolean,varbinary, varchar,scalar,true,Returns true if and only if left geometry is within right geometry
st_within,boolean,varchar, varbinary,scalar,true,Returns true if and only if left geometry is within right geometry
st_within,boolean,varchar, varchar,scalar,true,Returns true if and only if left geometry is within right geometry
st_x,double,varbinary,scalar,true,Returns the X coordinate of point
st_x,double,varchar,scalar,true,Returns the X coordinate of point
st_y,double,varbinary,scalar,true,Returns the Y coordinate of point
st_y,double,varchar,scalar,true,Returns the Y coordinate of point
stddev,double,bigint,aggregate,true,Returns the variance of the argument
stddev,double,double,aggregate,true,Returns the variance of the argument
stddev_pop,double,bigint,aggregate,true,Returns the variance of the argument
stddev_pop,double,double,aggregate,true,Returns the variance of the argument
stddev_samp,double,bigint,aggregate,true,Returns the variance of the argument
stddev_samp,double,double,aggregate,true,Returns the variance of the argument
strpos,bigint,varchar(x), varchar(y),scalar,true,returns index of first occurrence of a substring (or 0 if not found)
substr,char(x),char(x), bigint,scalar,true,suffix starting at given index
substr,char(x),char(x), bigint, bigint,scalar,true,substring of given length starting at an index
substr,varchar(x),varchar(x), bigint,scalar,true,suffix starting at given index
substr,varchar(x),varchar(x), bigint, bigint,scalar,true,substring of given length starting at an index
substring,varchar(x),varchar(x), bigint,scalar,true,suffix starting at given index
substring,varchar(x),varchar(x), bigint, bigint,scalar,true,substring of given length starting at an index
sum,bigint,bigint,aggregate,true,
sum,decimal(38,s),decimal(p,s),aggregate,true,Calculates the sum over the input values
sum,double,double,aggregate,true,
sum,real,real,aggregate,true,
tan,double,double,scalar,true,tangent
tanh,double,double,scalar,true,hyperbolic tangent
timezone_hour,bigint,timestamp with time zone,scalar,true,time zone hour of the given timestamp
timezone_minute,bigint,timestamp with time zone,scalar,true,time zone minute of the given timestamp
to_base,varchar(64),bigint, bigint,scalar,true,convert a number to a string in the given base
to_base64,varchar,varbinary,scalar,true,encode binary data as base64
to_base64url,varchar,varbinary,scalar,true,encode binary data as base64 using the URL safe alphabet
to_big_endian_64,varbinary,bigint,scalar,true,encode value as a 64-bit 2's complement big endian varbinary
to_char,varchar,timestamp with time zone, varchar,scalar,true,Formats a timestamp
to_date,date,varchar, varchar,scalar,true,Converts a string to a DATE data type
to_hex,varchar,varbinary,scalar,true,encode binary data as hex
to_iso8601,varchar(16),date,scalar,true,
to_iso8601,varchar(35),timestamp,scalar,true,
to_iso8601,varchar(35),timestamp with time zone,scalar,true,
to_timestamp,timestamp,varchar, varchar,scalar,true,Converts a string to a TIMESTAMP data type
to_unixtime,double,timestamp,scalar,true,
to_unixtime,double,timestamp with time zone,scalar,true,
to_utf8,varbinary,varchar(x),scalar,true,encodes the string to UTF-8
transform,array(U),array(T), function(T,U),scalar,false,apply lambda to each element of the array
transform_keys,map(K2,V),map(K1,V), function(K1,V,K2),scalar,false,apply lambda to each entry of the map and transform the key
transform_values,map(K,V2),map(K,V1), function(K,V1,V2),scalar,false,apply lambda to each entry of the map and transform the value
trim,char(x),char(x),scalar,true,removes whitespace from the beginning and end of a string
trim,char(x),char(x), CodePoints,scalar,true,remove the longest string containing only given characters from the beginning and end of a string
trim,varchar(x),varchar(x),scalar,true,removes whitespace from the beginning and end of a string
trim,varchar(x),varchar(x), CodePoints,scalar,true,remove the longest string containing only given characters from the beginning and end of a string
truncate,decimal(p,s),decimal(p,s), bigint,scalar,true,round to integer by dropping given number of digits after decimal point
truncate,decimal(rp,0),decimal(p,s),scalar,true,round to integer by dropping digits after decimal point
truncate,double,double,scalar,true,round to integer by dropping digits after decimal point
truncate,real,real,scalar,true,round to integer by dropping digits after decimal point
typeof,varchar,T,scalar,true,textual representation of expression type
upper,char(x),char(x),scalar,true,converts the string to upper case
upper,varchar(x),varchar(x),scalar,true,converts the string to upper case
url_decode,varchar(x),varchar(x),scalar,true,unescape a URL-encoded string
url_encode,varchar(y),varchar(x),scalar,true,escape a string for use in URL query parameter names and values
url_extract_fragment,varchar(x),varchar(x),scalar,true,extract fragment from url
url_extract_host,varchar(x),varchar(x),scalar,true,extract host from url
url_extract_parameter,varchar(x),varchar(x), varchar(y),scalar,true,extract query parameter from url
url_extract_path,varchar(x),varchar(x),scalar,true,extract part from url
url_extract_port,bigint,varchar(x),scalar,true,extract port from url
url_extract_protocol,varchar(x),varchar(x),scalar,true,extract protocol from url
url_extract_query,varchar(x),varchar(x),scalar,true,extract query from url
uuid,varchar,,scalar,false,Returns a randomly generated UUID
var_pop,double,bigint,aggregate,true,Returns the population variance of the argument
var_pop,double,double,aggregate,true,Returns the population variance of the argument
var_samp,double,bigint,aggregate,true,Returns the sample variance of the argument
var_samp,double,double,aggregate,true,Returns the sample variance of the argument
variance,double,bigint,aggregate,true,Returns the sample variance of the argument
variance,double,double,aggregate,true,Returns the sample variance of the argument
week,bigint,date,scalar,true,week of the year of the given date
week,bigint,timestamp,scalar,true,week of the year of the given timestamp
week,bigint,timestamp with time zone,scalar,true,week of the year of the given timestamp
week_of_year,bigint,date,scalar,true,week of the year of the given date
week_of_year,bigint,timestamp,scalar,true,week of the year of the given timestamp
week_of_year,bigint,timestamp with time zone,scalar,true,week of the year of the given timestamp
width_bucket,bigint,double, array(double),scalar,true,The bucket number of a value given an array of bins
width_bucket,bigint,double, double, double, bigint,scalar,true,The bucket number of a value given a lower and upper bound and the number of buckets
xxhash64,varbinary,varbinary,scalar,true,compute xxhash64 hash
year,bigint,date,scalar,true,year of the given date
year,bigint,interval year to month,scalar,true,year of the given interval
year,bigint,timestamp,scalar,true,year of the given timestamp
year,bigint,timestamp with time zone,scalar,true,year of the given timestamp
year_of_week,bigint,date,scalar,true,year of the ISO week of the given date
year_of_week,bigint,timestamp,scalar,true,year of the ISO week of the given timestamp
year_of_week,bigint,timestamp with time zone,scalar,true,year of the ISO week of the given timestamp
yow,bigint,date,scalar,true,year of the ISO week of the given date
yow,bigint,timestamp,scalar,true,year of the ISO week of the given timestamp
yow,bigint,timestamp with time zone,scalar,true,year of the ISO week of the given timestamp
zip,array(row(field0 T1,field1 T2)),array(T1), array(T2),scalar,true,Merges the given arrays, element-wise, into a single array of rows.
zip,array(row(field0 T1,field1 T2,field2 T3)),array(T1), array(T2), array(T3),scalar,true,Merges the given arrays, element-wise, into a single array of rows.
zip,array(row(field0 T1,field1 T2,field2 T3,field3 T4)),array(T1), array(T2), array(T3), array(T4),scalar,true,Merges the given arrays, element-wise, into a single array of rows.
zip_with,array(R),array(T), array(U), function(T,U,R),scalar,false,merge two arrays, element-wise, into a single array using the lambda function

========================================
prpt_name,prpt_value
EXTERNAL,TRUE
transient_lastDdlTime,1480278335

========================================
field
request_timestamp
elb_name
request_ip
request_port
backend_ip
backend_port
request_processing_time
backend_processing_time
client_response_time
elb_response_code
backend_response_code
received_bytes
sent_bytes
request_verb
url
protocol
user_agent
ssl_cipher
ssl_protocol

========================================
createtab_stmt
CREATE EXTERNAL TABLE `sampledb.elb_logs`(
  `request_timestamp` string COMMENT '',
  `elb_name` string COMMENT '',
  `request_ip` string COMMENT '',
  `request_port` int COMMENT '',
  `backend_ip` string COMMENT '',
  `backend_port` int COMMENT '',
  `request_processing_time` double COMMENT '',
  `backend_processing_time` double COMMENT '',
  `client_response_time` double COMMENT '',
  `elb_response_code` string COMMENT '',
  `backend_response_code` string COMMENT '',
  `received_bytes` bigint COMMENT '',
  `sent_bytes` bigint COMMENT '',
  `request_verb` string COMMENT '',
  `url` string COMMENT '',
  `protocol` string COMMENT '',
  `user_agent` string COMMENT '',
  `ssl_cipher` string COMMENT '',
  `ssl_protocol` string COMMENT '')
ROW FORMAT SERDE
  'org.apache.hadoop.hive.serde2.RegexSerDe'
WITH SERDEPROPERTIES (
  'input.regex'='([^ ]*) ([^ ]*) ([^ ]*):([0-9]*) ([^ ]*):([0-9]*) ([.0-9]*) ([.0-9]*) ([.0-9]*) (-|[0-9]*) (-|[0-9]*) ([-0-9]*) ([-0-9]*) \\\"([^ ]*) ([^ ]*) (- |[^ ]*)\\\" (\"[^\"]*\") ([A-Z0-9-]+) ([A-Za-z0-9.-]*)$')
STORED AS INPUTFORMAT
  'org.apache.hadoop.mapred.TextInputFormat'
OUTPUTFORMAT
  'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION
  's3://athena-examples-us-east-2/elb/plaintext'
TBLPROPERTIES (
  'transient_lastDdlTime'='1480278335')

========================================
View not found or not a valid presto view: henrywu
========================================
database_name
clickstreams
default
sampledb

========================================
database_name
sampledb

========================================
FAILED: Execution Error, return code 1 from org.apache.hadoop.hive.ql.exec.DDLTask. Table sampledb.elb_logs is not a partitioned table
========================================
tab_name
elb_logs_henrywu
elb_logs_new2
testme
testme2

========================================
views

========================================
*/
