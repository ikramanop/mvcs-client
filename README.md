### Запуск для разработки под виндой

```shell
cd test_project

cmd /V /C 'set "CGO_ENABLED=1" && go run ..\main.go *command* *subcommand*'
```

### Запуск для разработки под линуксом/маком

```shell
cd test_project

CGO_ENABLED=1 go run ../main.go *command* *subcommand*
```