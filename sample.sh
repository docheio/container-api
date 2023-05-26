curl -s -X GET localhost:8081/v1/
curl -s -X DELETE localhost:8081/v1/
curl -s -X POST -d '{"cpu":1000,"mem":1000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"mount":"/root/minecraft","size":5}]}' 127.0.0.1:8081/v1/
curl -s -X PUT -d '{"cpu":2000,"mem":2000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"id":"","mount":"/root/minecraft","size":6}]}' 127.0.0.1:8081/v1/
curl -s -X PUT -d '{"cpu":2000,"mem":2000,"ports":[],"pvcs":[]}' 127.0.0.1:8081/v1/