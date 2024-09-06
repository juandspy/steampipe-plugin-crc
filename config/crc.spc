connection "crc" {
  plugin = "local/crc"

  # The baseUrl (prod or stage) for the console.redhat.com APIs
  # Can also be set with the CRC_URL environment variable.
  base_url = "https://console.redhat.com/"

  # The tokenUrl (prod or stage) for updating the token used to communicate
  # with console.redhat.com
  # Can also be set with the CRC_TOKEN_URL environment variable.
  token_url = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"

  # The client ID to access the console.redhat.com cloud instance
  # Can also be set with the `CRC_CLIENT_ID` environment variable.
  # client_id = "XXXXX"

  # The client secret to access the console.redhat.com cloud instance
  # Can also be set with the `CRC_CLIENT_SECRET` environment variable.
  # client_secret = "XXXXX"
}