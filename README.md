# goto
- master
  - go run . -http=:8080 -rpc=true
- slave
  - go run . -master=127.0.0.1:8080 -http=:8081
- 在wsl上运行，使用localhost:8081访问,但使用127.0.0.1:8081无法访问