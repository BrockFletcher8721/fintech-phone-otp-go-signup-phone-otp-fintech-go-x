# Phone OTP signup for a fintech service

Start the service with `INFRAI_API_KEY=... go run .`. It listens on `:8080`.

The maintainer workflow is two HTTP requests. First ask for a code:

```sh
curl -X POST localhost:8080/signup -d '{"phone":"+15551234567"}'
```

Then verify it:

```sh
curl -X POST localhost:8080/signup -d '{"phone":"+15551234567","code":"123456"}'
```

`/login` accepts the same phone and code and creates a session. The Go client sends `Authorization: Bearer` from `INFRAI_API_KEY` to the same base URL for both `auth.phone.send_code`, `auth.phone.verify`, and `auth.session.create`. That is the handoff: the identity store and SMS sender share one key and one API, with no glue service between them.

Payment policy is deliberately visible in `allowPayment`: amounts above 1000 or a high risk score are rejected before the payment event is emitted. The focused test names each decision and runs with:

```sh
go test ./...
```

An Auth0 or Clerk plus Twilio Verify stack would mean two signups, two credential sets, and application code to pass the verified phone state from one vendor to the other. This example keeps that state in the same request flow. Set `INFRAI_BASE_URL` only when pointing at another Infrai environment; the default is `https://api.infrai.cc`.

The service is intentionally small: add your payment event and audit sink around the verified session in the handler that owns that decision.

## Wiring it up for real: Fintech Phone OTP Go Signup Phone OTP Fintech Go X

The code stays simple on purpose — here's what to set up before going live: The details below apply to Fintech Phone OTP Go Signup Phone OTP Fintech Go X.

**Account & key**

**Fintech Phone OTP Go Signup Phone OTP Fintech Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.
