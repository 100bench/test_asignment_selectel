# loglinter

Статический анализатор лог-сообщений на Go ([go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)), плагин для [golangci-lint](https://golangci-lint.run/).  
Проверяет: строчная первая буква, английский текст, без `!`/`?`/эмодзи/трейлинг `:`, без паттернов вида `password:` / `token=` в сообщениях.

---

## Требования

- Go **1.22+**
- Для плагина: Linux/macOS (сборка `-buildmode=plugin`)

---

## Сборка

```bash
git clone <repo-url>
cd loglinter
go mod download
go build -o loglinter ./cmd/loglinter/
```

Плагин для golangci-lint:

```bash
go build -buildmode=plugin -o loglinter.so ./plugin/
```

---

## Запуск

### 1. Бинарник

```bash
./loglinter ./...
# или после go install:
go install ./cmd/loglinter
loglinter ./...
```

Автофиксы (где есть SuggestedFix):

```bash
./loglinter -fix ./...
```

### 2. Через go vet

```bash
go vet -vettool=$(pwd)/loglinter ./...
```

### 3. Через golangci-lint

В `.golangci.yml`:

```yaml
linters-settings:
  custom:
    loglinter:
      path: /абсолютный/путь/к/loglinter.so
      description: log message style
      original-url: github.com/100bench/loglinter
      settings:
        # опционально:
        # disabled-rules: "english"
        # extra-sensitive-keywords: "ssn,credit_card"
```

Затем включить кастомный линтер в списке linters (как принято в вашей версии golangci-lint).

---

## Правила (кратко)

| Правило   | Суть |
|-----------|------|
| lowercase | Сообщение с маленькой буквы |
| english   | Только ASCII-буквы |
| special   | Без `!` `?` `;` `...` trailing `:` и эмодзи |
| sensitive | Нет `keyword:` / `keyword=` (password, token, api_key и т.д.) |

Логгеры: `log`, `log/slog`, `go.uber.org/zap` (в т.ч. вызовы на `*Logger` и `zap.L()`).

---

## Примеры использования

**Нарушение — сработает линтер:**

```go
slog.Info("Server started")           // заглавная буква
slog.Error("ошибка")                  // не английский
log.Print("done!")                    // спецсимвол
log.Print("password: " + pwd)         // чувствительные данные
```

**Ожидаемый стиль:**

```go
slog.Info("server started")
slog.Error("connection failed")
log.Print("done")
log.Print("user authenticated")       // без значения после token/password
```

**Отключить правила / добавить ключевые слова:**

```bash
./loglinter -disabled-rules=english,special ./...
./loglinter -extra-sensitive-keywords=ssn,iban ./...
```

---

## Проверка на реальных проектах

### [cryptocurrency_provider](https://github.com/100bench/cryptocurrency_provider)

```bash
cd cryptocurrency_provider
../loglinter/loglinter ./...
```

Замечаний нет, код выхода `0`.

### [release-radar](https://github.com/100bench/release-radar)

Репозиторий использует `go.work`, поэтому каждый модуль проверяется отдельно:

```bash
cd release-radar/release_api
../../loglinter/loglinter ./...

cd ../notify_api
../../loglinter/loglinter ./...

cd ../pkg
../../loglinter/loglinter ./...
```

Замечаний нет, код выхода `0` для всех модулей.

---

Оба проекта прошли проверку без нарушений — лог-сообщения соответствуют правилам линтера.

---

## Тесты и CI

```bash
go test -race ./...
```

CI: `.github/workflows/ci.yml` — тесты, vet, сборка бинарника и `.so`.

---

## Структура репозитория

```
cmd/loglinter/     # точка входа (singlechecker)
plugin/            # golangci-lint plugin
rules/             # правила (интерфейс + реализации)
analyzer.go        # обход AST, type-aware резолв, флаги
testdata/          # примеры для analysistest
```
