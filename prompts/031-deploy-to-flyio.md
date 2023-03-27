# Deploying Go app on Fly.io.

> I have a Go app that I'd like to deploy to Fly.io. The app has several endpoints listening on port 8080. It requires access to a Postgres database; the connection parameters are passed on via environment variables, as well as some other (secret) configuration, like this:
>
> export DB_HOST=localhost
> export DB_PORT=5432
> export DB_USER=usr
> export DB_PASS=pwd
> export DB_NAME=balgpt
> export AUTH_USERNAME=admin
> export AUTH_PASSWORD=admin
> export SCRAPER_URL=https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/
> export API_BASE_URL=http://localhost:8080
>
> What do I need to do to get it up and running?

To deploy your Go app on Fly.io, you'll need to follow these steps:

1. Install the Fly CLI:
Visit https://fly.io/docs/getting-started/installing-flyctl/ and follow the instructions for your operating system.

2. Create a fly.toml configuration file in your project's root directory. This file will tell Fly how to build and deploy your app. Here's a basic example for your Go app:makefileCopy codeapp = "your-app-name"

```makefile
app = "your-app-name"

[build]
builder = "paketocommunity/go:full"
buildpacks = ["gcr.io/paketo-buildpacks/go"]

[env]
DB_HOST = "localhost"
DB_PORT = "5432"
DB_USER = "usr"
DB_PASS = "pwd"
DB_NAME = "balgpt"
AUTH_USERNAME = "admin"
AUTH_PASSWORD = "admin"
SCRAPER_URL = "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"
API_BASE_URL = "http://localhost:8080"

[[services]]
internal_port = 8080
protocol = "tcp"

[services.concurrency]
hard_limit = 25
soft_limit = 20

[[services.ports]]
handlers = ["http"]
port = "80"

[[services.ports]]
handlers = ["tls", "http"]
port = "443"

[services.tls]
cert_type = "automated"
```

Replace "your-app-name" with a unique name for your app.

3. Set up your Postgres database:
Fly.io provides a managed PostgreSQL service called Fly Postgres. You can follow the instructions at https://fly.io/docs/reference/postgres/ to create a new database and obtain the connection parameters (host, port, user, password, and database name).Update the [env] section of your fly.toml with the new Postgres connection parameters.

4. Configure secret environment variables:
To securely store secret configuration, such as AUTH_USERNAME and AUTH_PASSWORD, you should use Fly's secrets management. Remove these secrets from the fly.toml file and set them using the Fly CLI:
```arduino
fly secrets set AUTH_USERNAME=admin
fly secrets set AUTH_PASSWORD=admin
```
These secrets will be available as environment variables in your app, just like the ones defined in the fly.toml.

5. Initialize your app with Fly:
```csharp
fly init
```
If prompted, select the builder and app name that match your fly.toml.

6. Deploy your app:
```
fly deploy
```
This command will build and deploy your app to Fly.io. It may take a few minutes to complete.

7. Check your app status:
```lua
fly status
```
This command will show you the current status of your app, including the deployed version and the number of instances running.

8. Visit your app:
You can now access your app using the URL provided in the output of the fly status command or the custom domain you've configured.

For more information on deploying apps to Fly.io, visit their documentation at [https://fly.io/docs/](https://fly.io/docs/).

## Follow-up prompt

> Can the build section be simplified? Also, does https work out of the box?

Yes, the build section can be simplified. If you're using Go modules for dependency management, Fly.io can automatically detect your Go app without needing a specific builder or buildpack. You can remove the `[build]` section entirely, and your `fly.toml` would look like this:

```toml
app = "your-app-name"

[env]
DB_HOST = "localhost"
DB_PORT = "5432"
DB_USER = "usr"
DB_PASS = "pwd"
DB_NAME = "balgpt"
SCRAPER_URL = "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"
API_BASE_URL = "http://localhost:8080"

[[services]]
internal_port = 8080
protocol = "tcp"

[services.concurrency]
hard_limit = 25
soft_limit = 20

[[services.ports]]
handlers = ["http"]
port = "80"

[[services.ports]]
handlers = ["tls", "http"]
port = "443"

[services.tls]
cert_type = "automated"

```

For HTTPS, Fly.io automatically handles TLS termination for your app, as shown in the `fly.toml` configuration above. The `[[services.ports]]` sections define the port mappings for HTTP (port 80) and HTTPS (port 443). The `handlers` property specifies that Fly should handle TLS termination for port 443. The `[services.tls]` section configures the certificate type to be "automated", which means Fly will automatically provision and renew Let's Encrypt certificates for your custom domain or the default Fly app domain.

So, yes, HTTPS will work out-of-the-box with the provided configuration. Just make sure to update the `API_BASE_URL` environment variable to use the correct base URL with the `https` scheme when accessing your app's endpoints.

## Follow-up prompt

> Can you give me a GitHub actions workflow for releasing to fly.io?

Certainly! Here's a GitHub Actions workflow that automatically deploys your Go app to Fly.io whenever you push changes to the main branch. Create a new file named `.github/workflows/fly-deploy.yml` in your repository and add the following content:

```yaml
name: Deploy to Fly.io

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17

    - name: Install Fly CLI
      run: |
        curl -sL https://install.fly.io/flyctl-install.sh | sh
        echo "$GITHUB_WORKSPACE/bin" >> $GITHUB_PATH

    - name: Deploy to Fly.io
      env:
        FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
      run: |
        fly login --api-token "$FLY_API_TOKEN"
        fly deploy

```

This workflow performs the following steps:

1. Triggers on push to the main branch.
2. Sets up the Go environment with the specified Go version.
3. Installs the Fly CLI.
4. Deploys your app to Fly.io using the Fly CLI.

Before you can use this workflow, you need to obtain an API token from Fly.io:

1. Visit https://fly.io/sites and log in to your account.
2. Click on your avatar in the top-right corner and go to "API Tokens."
3. Click "Create token" and give it a descriptive name. Copy the generated token.

Next, you need to add the API token as a secret in your GitHub repository:

1. Go to your repository on GitHub and click on the "Settings" tab.
2. In the left sidebar, click on "Secrets."
3. Click the "New repository secret" button.
4. Name the secret "FLY_API_TOKEN" and paste the token you copied earlier in the "Value" field. Click "Add secret."

Now, whenever you push changes to the main branch, this workflow will automatically deploy your Go app to Fly.io using the specified configuration.

## Follow-up prompt

> Awesome! Can you integrate that with the existing workflow? I want to build and test on each push regardless to what branch, and deploy only if the push was to master. Also, I don't want to repeat any logic, so the deploy jobs should depend on the existing build and test job.
>
> (snipped: original GitHub Actions workflow)

You can modify the existing workflow to build and test on each push, regardless of the branch, and deploy only if the push was to the main branch. Here's your updated workflow:

```yaml
name: CI/CD Pipeline

on:
  push:

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20.2

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        run: |
          go mod download

      - name: Build
        run: |
          go build

      - name: Test
        run: |
          go test -v ./...

  deploy:
    runs-on: ubuntu-latest
    needs: build-and-test
    if: github.ref == 'refs/heads/main'
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Install Fly CLI
        run: |
          curl -sL https://install.fly.io/flyctl-install.sh | sh
          echo "$GITHUB_WORKSPACE/bin" >> $GITHUB_PATH

      - name: Deploy to Fly.io
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
        run: |
          fly login --api-token "$FLY_API_TOKEN"
          fly deploy
```

This updated workflow has the following changes:

1. The `on` trigger has been modified to run on each push, regardless of the branch.

2. A new `deploy` job has been added, which depends on the `build-and-test` job using the `needs` keyword. The `deploy` job will only run if the push was to the main branch, as specified by the `if` condition: `if: github.ref == 'refs/heads/main'`.

Now, the workflow will build and test your Go app on each push. If the push was to the main branch, and the build-and-test job succeeded, it will also deploy the app to Fly.io.

## Note

The hyperlinks were all good!

It got the order of the steps wrong though: you have to do `fly init` (which is now called `fly launch`) before you can set secrets.

I had to create a database in the Postgres instance manually:

- `fly pg connect --app balgpt-pg`
- `CREATE DATABASE balgpt;`

Afterwards, I could use `flyctl proxy 5432 -a balgpt-pg` to locally connect to the database, for instance to run the Python script to import historical data.

Finally, the GitHub workflow for deploying to Fly.io was out of date, so I updated accordingly.
