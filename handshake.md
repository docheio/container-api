# Handshake List
| endpoint | method | description | type                                                                                                         | response                                                                                                                                     | 
| -------- | ------ | ----------- | ------------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- | 
| /        | POST   | create      | {cpu:uint16,mem:uint16,ports:[]{protocol:string,internal:uint16},pvc:[]{id:string,mount:string,size:uint16}} | []{id:string,name:string,status:string,ports:[]{protocol:string,internal:uint16,external:uint16},pvc:[]{id:string,mount:string,size:uint16}} | 
| /        | GET    | get all     | {}                                                                                                           | []{id:string,name:string,status:string,ports:[]{protocol:string,internal:uint16,external:uint16},pvc:[]{id:string,mount:string,size:uint16}} | 
| /:name   | GET    | get one     | {}                                                                                                           | []{id:string,name:string,status:string,ports:[]{protocol:string,internal:uint16,external:uint16},pvc:[]{id:string,mount:string,size:uint16}} | 
| /:name   | PUT    | update      | {cpu:uint16,mem:uint16,ports:[]{protocol:string,internal:uint16},pvc:[]{id:string,mount:string,size:uint16}} | []{id:string,name:string,status:string,ports:[]{protocol:string,internal:uint16,external:uint16},pvc:[]{id:string,mount:string,size:uint16}} | 
| /:name   | DELETE | delete      | {}                                                                                                           | []{id:string,name:string,status:string,ports:[]{protocol:string,internal:uint16,external:uint16},pvc:[]{id:string,mount:string,size:uint16}} |

# Usage (curl)
```
curl -s -X GET 127.0.0.1:8081/v1/
curl -s -X DELETE 127.0.0.1:8081/v1/
curl -s -X POST -d '{"cpu":1000,"mem":1000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"mount":"/root/minecraft","size":5}]}' 127.0.0.1:8081/v1/
curl -s -X PUT -d '{"cpu":2000,"mem":2000,"ports":[{"protocol":"UDP","internal":19132}],"pvcs":[{"id":"","mount":"/root/minecraft","size":6}]}' 127.0.0.1:8081/v1/
curl -s -X PUT -d '{"cpu":2000,"mem":2000,"ports":[],"pvcs":[]}' 127.0.0.1:8081/v1/
```