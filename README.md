## Borsch Playground Service
Build and set up the application:
```shell
mkdir "bin"
go build -o ./bin/borschplayground main.go
cp ./settings.json ./bin
cd ./bin
```

Migrate the database:
```shell
./borschplayground --address 127.0.0.1:8080
```

Run the server:
```shell
./borschplayground --address 127.0.0.1:8080
```

### API
* `[POST] /api/v1/run`
  
  Input:
  ```json
  {
    "lang_v": "0.1.0",
    "b64_source": "0LTRgNGD0LrRgCgi0J/RgNC40LLRltGCLCDQodCy0ZbRgtC1IDEhIik7"
  }
  ```
  Output:
  ```json
  {
    "job_id": "<uuid>"
  }
  ```

  Source code which is base64-encoded:
  ```
  друкр("Привіт, Світе 1!");
  ```

* `[GET] /api/v1/job/<uuid>`
* `[GET] /api/v1/job/<uuid>/output`
