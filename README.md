# OAuth CLI with GitHub 

This is a command-line interface (CLI) application written in Go that authenticates users using the OAuth 2.0 protocol with GitHub. The program allows users to access resources from the GitHub API on their behalf, without needing to provide their credentials to the CLI.

## Installation

To use this program, you will need to have Go installed on your system. You can download and install Go from the official website at https://golang.org/dl/. Once you have installed Go, you can clone this repository and build the program using the following commands:

```bash
# clone the repo
git clone git@github.com:pyadav/how-to-authenticate-cli-using-oauth.git
cd how-to-authenticate-cli-using-oauth
go mod tidy
```

## Usage

To run the program, you will need to have a GitHub account and create a new GitHub application on the GitHub developer settings page. Navigate to "Settings" -> "Developer settings" -> "OAuth Apps" and click "New OAuth App". Fill out the required information, including the callback URL, which is the URL where the user will be redirected after authorizing the application.

Once you have created your GitHub application, update your .env file with github client_id and client_secret then you can run the CLI with the following command:

``` go
go run main.go
```

The program will start a local server and launch the user's default web browser for granting permission to authorize the program to access GitHub resources. Once the user grants permission, the authorization server will redirect the request to the predefined redirect URL.

The CLI will then parse the redirect request and receive the authorization code. The program will use this authorization code to exchange for an access token calling the authorization server endpoint. The program will save this token to a file named token.json with 0600 permissions and can make API requests on behalf of the user.

**Note:**  if you need to use a different scope, you will need to modify the Scopes field in the oauth2.Config struct accordingly.

## Contributing

If you find any issues with this program or would like to suggest improvements, please feel free to submit a pull request or open an issue on the GitHub repository.

## License

This program is licensed under the [Apache 2.0](LICENSE). See the [LICENSE file](LICENSE) for more information.
