# Examples

This directory contains example programs demonstrating how to use the Ghion Go SDK.

## Prerequisites

Before running any example, set your API credentials as environment variables:

```bash
export GHION_API_KEY=your_api_key_here
export GHION_API_SECRET=your_api_secret_here
export GHION_API_PASSPHRASE=your_passphrase_here
```

## Available Examples

### quick-start.go

A simple example showing basic SDK usage:

- Initializing the client
- Creating a payment
- Getting payment status

**Run:**
```bash
go run examples/quick-start.go
```

### http-server.go

A complete HTTP server example with webhook handling:

- REST API endpoints for payment operations
- Webhook signature verification
- Event handling for different webhook types

**Run:**
```bash
export PORT=8080  # Optional, defaults to 8080
go run examples/http-server.go
```

**Test endpoints:**
```bash
# Initialize payment
curl -X POST http://localhost:8080/api/payments/initialize \
  -H "Content-Type: application/json" \
  -d '{"amount":100,"currency":"ETB","reference":"order_123"}'

# Get payment status
curl http://localhost:8080/api/payments/{payment_id}
```

### otp-flow.go

Complete OTP payment flow example:

- Initialize payment
- Send OTP to customer's phone
- Validate OTP code
- Get final payment status

**Run:**
```bash
go run examples/otp-flow.go
```

### qr-payment.go

QR code payment example:

- Initialize payment
- Generate QR code
- Display QR information
- Get checkout details

**Run:**
```bash
go run examples/qr-payment.go
```

## Testing Webhooks Locally

To test webhooks locally, use a tool like ngrok to expose your local server:

```bash
# Install ngrok
# Then run:
ngrok http 8080

# Use the ngrok URL in your webhook configuration
```

## Common Issues

### Invalid Credentials

Make sure your API credentials are correct and set as environment variables:

```bash
echo $GHION_API_KEY
echo $GHION_API_SECRET
echo $GHION_API_PASSPHRASE
```

### Network Errors

If you encounter network errors, check:
- Your internet connection
- Firewall settings
- API endpoint availability

### Signature Verification Failures

For webhook signature verification:
- Ensure you're reading the raw request body
- Check the X-Ghion-Signature header
- Verify your API secret is correct

## Additional Resources

- [Main README](../README.md)
- [Developer Guide](../docs/DEVELOPER_GUIDE.md)
- [Ghion API Documentation](https://ghion.financial/docs)
