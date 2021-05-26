# Panacea TX Generator

for the Incentivized Testnet


## Building

To get some dependencies from GitBub Packages maven repo, your GitHub username and token must be set as environment variables.
```bash
GPR_USER=<your-github-username> \
GPR_API_KEY=<your-github-api-key> \
./gradlew build
```

## Running

```bash
MNEMONIC=<mnemonic> \
LCD_ENDPOINT=<lcd_http_addr> \
./gradlew run
```
