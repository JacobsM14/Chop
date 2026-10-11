# Chop

Chop es una CLI para guardar listas de comandos de terminal (las *chops*) y
ejecutarlas después en orden. Sirve para esas secuencias que repites a menudo
—desplegar, levantar un entorno, limpiar cachés— y que no merecen un script
propio.

Cada chop tiene un nombre, un color y, opcionalmente, numeración. Cada item
guarda un comando, una descripción y el color con que se muestra cada uno.
Todo se guarda en una base de datos SQLite local.

## Requisitos

- Go 1.27.2 o superior
- `sh` en el `PATH` (los comandos se ejecutan con `sh -c`)

No hace falta CGO: el driver de SQLite ([modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)) es Go puro.

## Instalación

```sh
git clone <url-del-repo> chop
cd chop
make install   # instala el binario en $GOBIN o $GOPATH/bin
```

Si prefieres no instalarlo, `make build` deja el binario en `bin/chop`.

## Uso rápido

```sh
chop new deploy --color green --ordered   # crea una chop numerada
chop add deploy                           # añade un item (formulario interactivo)
chop add deploy
chop get deploy                           # muestra sus items
chop run deploy                           # pide confirmación y los ejecuta en orden
```

Salida de `chop get deploy` con dos items:

```
deploy
──────
1. go test ./...: pasar los tests
2. git push: subir los cambios
```

## Comandos

| Comando | Alias | Descripción |
|---|---|---|
| `chop new <nombre>` | | Crea una chop |
| `chop add <nombre>` | | Añade un item a una chop |
| `chop ls` | | Lista todas las chops |
| `chop get <nombre>` | `ver` | Muestra una chop y sus items |
| `chop run <nombre>` | | Ejecuta los comandos de una chop en orden |
| `chop rm-item <nombre> [n]` | | Elimina un item de una chop |
| `chop delete <nombre>` | `rm` | Elimina una chop y todos sus items |

### `new`

```sh
chop new <nombre> [--color <color>] [--ordered]
```

- Sin flags, abre un formulario que pide el color y si la lista es numerada.
- Con `--ordered` pero sin `--color`, solo pregunta el color.
- Con `--color`, no pregunta nada; la lista queda sin numerar salvo que pases `--ordered`.

Colores válidos: `red`, `green`, `blue`, `purple`, `orange`, `white`.

### `add`

```sh
chop add <nombre>
```

Siempre es interactivo. Pide el comando, la descripción (ninguno puede estar
vacío) y el color de cada uno. Por defecto, el comando sale en `purple` y la
descripción en `white`. El item se añade al final de la lista.

### `ls`

```sh
chop ls [--sort name|color|ordered]
```

Muestra el nombre de cada chop en su color. Por defecto ordena por nombre.

### `get`

```sh
chop get <nombre>
```

Muestra el título y los items. Las chops numeradas usan `1.`, `2.`…; las demás, `-`.

### `run`

```sh
chop run <nombre> [flags]
```

Antes de ejecutar, enseña la lista de comandos y pide confirmación. Cada
comando corre con `sh -c` conectado a tu terminal, así que los comandos
interactivos funcionan. Por defecto se detiene en el primer comando que falla.

| Flag | Descripción |
|---|---|
| `-i, --item N` | Ejecuta solo el item N (con la numeración de `chop get`) |
| `--dry-run` | Imprime los comandos sin ejecutarlos ni pedir confirmación |
| `--continue` | Sigue con el resto aunque un comando falle; al final informa de cuántos fallaron |
| `-y, --yes` | No pide confirmación |

`Ctrl+C` interrumpe el comando en curso y detiene la ejecución.

### `rm-item`

```sh
chop rm-item <nombre> [n]
```

Elimina el item `n` (con la numeración de `chop get`). Si no indicas `n`, abre
un selector para elegirlo.

### `delete`

```sh
chop delete <nombre> [-y]
```

Elimina la chop y todos sus items. Pide confirmación salvo que pases `-y`/`--yes`.

## Configuración

La base de datos se guarda en `~/.local/share/chop/chop.db`. El directorio se
crea automáticamente la primera vez.

Para usar otra ruta, define `CHOP_DB`:

```sh
CHOP_DB=/tmp/pruebas.db chop ls
```

Es útil para probar sin tocar tus chops reales.

## Desarrollo

```sh
make            # compila en bin/chop (objetivo por defecto)
make run ARGS="get deploy"
make test       # tests con -race
make check      # fmt + vet + test
make help       # lista todos los objetivos
```

Estructura del proyecto:

```
.
├── main.go
├── cmd/                 # un archivo por subcomando (cobra)
└── internal/
    ├── store/           # interfaz Store e implementación SQLite
    ├── render/          # colores y formato de salida (lipgloss)
    └── runner/          # ejecución de comandos con sh -c
```

Los formularios interactivos usan [huh](https://github.com/charmbracelet/huh).
