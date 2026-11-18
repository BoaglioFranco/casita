Casita is a household management app I'm developing for fun and practice that me and my and my wife use. It's running self hosted on a server in my office.  

## Local development

From the repository root:

```bash
make dev        # start PostgreSQL and all app servers, then attach to tmux
make dev-start  # start without attaching
make dev-stop   # stop app servers; leave PostgreSQL and its data running
```
