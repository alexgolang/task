# ⚠️ DEVELOPMENT CERTIFICATES ONLY

**DO NOT USE THESE CERTIFICATES IN PRODUCTION!**

These certificates are provided for development and testing purposes only.

## For Production Deployment:

1. **Generate your own certificates** using the commands in the main README.md
2. **Replace these files** with your own certificates
3. **Never commit real certificates** to version control

## Files in this directory:

- `server.key` - Server private key (for JWT signing)
- `server.crt` - Server certificate 
- `client.key` - Client private key (for testing client assertions)
- `client.crt` - Client certificate (for testing client assertions)
- `client.crt.der` - Client certificate in DER format

## Security Notes:

- These are **self-signed test certificates**
- They have **no security value** for production
- Anyone with access to this repository can see these keys
- **Always generate fresh certificates** for any real deployment

For certificate generation instructions, see the main README.md file.