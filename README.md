# accounts

This repository includes a go service called accounts that using embedded OPA for policy decisions.

## Download & Run
- Grab the latest [binary](https://github.com/ronenl1/accounts/releases) from the releases page and run it using the following command:
### Linux
```bash
./accounts_linux_amd64
```
### Windows
```powershell
.\accounts_windows.exe
```
### MacOS
```bash
./accounts_darwin_amd64
```
- Clone the repo sources and run the service using [Go CLI](https://golang.org/dl/):
```bash
git clone https://github.com/ronenl1/accounts.git
```
Run the following command from the **root** of the repo:
```bash
go run ./cmd/main.go
```

When the service is ready to accept requests you will see: "OPA engine is up!"

The service uses `data.json` as the source of truth, if you want to add more accounts just edit this file and restart the service.

The service will listen on port :7777 in `/accounts/{id}`. You will get response only from account IDs that specified in `data.json`.

The policies will allow users to access their own account ID or any user with costumer-service role only to the accounts in his region.

Example for authorized request:
```bash
curl -H username:alice -H region:EU -H roles:customer-service http://localhost:7777/accounts/2
```

Example for unauthorized request:
```bash
curl -H username:bob -H region:US http://localhost:7777/accounts/3
```
