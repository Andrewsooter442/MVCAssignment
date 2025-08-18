**If you didn't change the JWT secret. **

**Adming JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MSwiTmFtZSI6ImFkbWluIiwiSXNBZG1pbiI6dHJ1ZSwiSXNDaGVmZiI6ZmFsc2UsImlzcyI6ImtpdGNoZW4tYXBwIiwiZXhwIjoxNzU1NTQxMjU5fQ.gvNc02q8qq5cLuaARp6CF8MOcCiyigldXz-Mb8yOl5c**

**Chef JWT: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MiwiTmFtZSI6ImNoZWYiLCJJc0FkbWluIjpmYWxzZSwiSXNDaGVmZiI6dHJ1ZSwiaXNzIjoia2l0Y2hlbi1hcHAiLCJleHAiOjE3NTU1NDE0OTR9.ZrYHt7cCdB4_ZTqnhrE7y_pDKVUH9lZVaQ2mlXHn8gc**

The JWT of a new user will be logged in the console when he logs in 

Test command used
```
ab -n 100000 -c 150 -C "token=JWT_Token" http://localhost:8080/
```
### My test results.

**For admin**

20.778 [ms] (mean)
```
Server Software:
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        15527 bytes

Concurrency Level:      150
Time taken for tests:   13.852 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      1562300000 bytes
HTML transferred:       1552700000 bytes
Requests per second:    7219.17 [#/sec] (mean)
Time per request:       20.778 [ms] (mean)
Time per request:       0.139 [ms] (mean, across all concurrent requests)
Transfer rate:          110141.72 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  30.1      0    1019
Processing:     1   19  20.5     17     262
Waiting:        1   14   6.4     13      69
Total:          1   20  36.4     17    1234

Percentage of the requests served within a certain time (ms)
  50%     17
  66%     19
  75%     21
  80%     23
  90%     27
  95%     31
  98%     39
  99%     61
 100%   1234 (longest request)
```

**For chef**
10.903 [ms] (mean)

```
Server Software:
Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        1984 bytes

Concurrency Level:      150
Time taken for tests:   7.269 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      210200000 bytes
HTML transferred:       198400000 bytes
Requests per second:    13757.87 [#/sec] (mean)
Time per request:       10.903 [ms] (mean)
Time per request:       0.073 [ms] (mean, across all concurrent requests)
Transfer rate:          28241.26 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  10.9      0    1032
Processing:     1   10   4.2     10      56
Waiting:        1   10   4.1      9      56
Total:          1   11  11.6     10    1037

Percentage of the requests served within a certain time (ms)
  50%     10
  66%     12
  75%     12
  80%     13
  90%     15
  95%     17
  98%     20
  99%     29
 100%   1037 (longest request)
```

**For User** 

27.328 [ms] (mean)

```
Server Software:
    Server Hostname:        localhost
Server Port:            8080

Document Path:          /
Document Length:        30884 bytes

Concurrency Level:      150
Time taken for tests:   18.219 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      3098000000 bytes
HTML transferred:       3088400000 bytes
Requests per second:    5488.90 [#/sec] (mean)
Time per request:       27.328 [ms] (mean)
Time per request:       0.182 [ms] (mean, across all concurrent requests)
Transfer rate:          166060.76 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1  28.3      0    1018
Processing:     2   26   9.8     25     260
Waiting:        2   19   5.8     18      98
Total:          3   27  29.6     25    1227

Percentage of the requests served within a certain time (ms)
  50%     25
  66%     27
  75%     29
  80%     30
  90%     33
  95%     37
  98%     43
  99%     49
 100%   1227 (longest request)
```