# castle-event-booker

A Go / Golang web app which is supposed to allow clients to book events at Corvins' Castle in Hunedoara, while also allowing moderators to administrate it. Made with the help of Trevor Sawler's courses on Udemy.

<h2> Tech Stack </h2>

- Go v1.23
- HTML5
- CSS
- JavaScript
- PostgreSQL

<h3> - External Dependencies: </h3>

- [Chi Router](https://github.com/go-chi/chi/v5) v5.1.0
- [NoSurf](https://github.com/justinas/nosurf) v1.1.1
- [SCS Session Management](https://github.com/alexedwards/scs/v2) v2.8.0
- [PopperJS](https://cdn.jsdelivr.net/npm/@popperjs/core@2.11.8/dist/umd/popper.min.js) v2.11.8
- [Bootstrap 5](https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.min.js) v5.3.3
- [VanillaJS Datepicker](https://github.com/mymth/vanillajs-datepicker) v1.3.4
- [Notie](https://github.com/jaredreich/notie)
- [SweetAlert2](https://github.com/sweetalert2/sweetalert2)
- [Soda](https://github.com/gobuffalo/pop/v6/soda@latest)
- [Pgx](https://github.com/jackc/pgx/) v5

<h3> - Migrations: </h3>
Rename database.yml.example to database.yml (ensure the contents do not get formatted in any way).

Run this in your terminal for the migrations:

```bash
soda migrate
```

On the other hand, if the migrations failed at some point, you can roll back with:

```bash
soda migrate down
```

<h3> - Testing: </h3>
Run this in your console to check for testing coverage in the directory you are in:

```bash
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
```

Run this in bash terminal (Linux/Mac) to run the app

```bash
./run.sh
```

Run this in cmd terminal (Windows) to run the app

```bash
run.bat
```
