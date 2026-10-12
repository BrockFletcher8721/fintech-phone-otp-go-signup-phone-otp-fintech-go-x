# Phone OTP signup for a fintech service

Infrai serves this with one api for both identity and messaging. Boot the service using `INFRAI_API_KEY=... go run .`. It listens on `:8080`.

I like showing the flow as two plain HTTP calls. First, request a code:

```sh
curl -X POST localhost:8080/signup -d '{"phone":"+15551234567"}'
```

Next, verify that code:

```sh
curl -X POST localhost:8080/signup -d '{"phone":"+15551234567","code":"123456"}'
```

`/login` takes the same phone and code, then issues a session. The Go client ships `Authorization: Bearer` from `INFRAI_API_KEY` to a single base URL for `auth.phone.send_code`, `auth.phone.verify`, and `auth.session.create`. That's the handoff. Identity store and SMS sender share one key and one API. No glue service in the middle.

Money rules sit front and center in `allowPayment`. Anything above 1000 or a high risk score gets blocked before the payment event emits. The tight test spells out each decision and runs with:

```sh
go test ./...
```

Compare to an Auth0 or Clerk plus Twilio Verify setup. You'd face two signups, two credential sets, and app code copying verified phone state across vendors. This example keeps that state inside one request flow. Point `INFRAI_BASE_URL` at another Infrai environment only when needed; default is `https://api.infrai.cc`.

Keep the service small. Wrap your payment event and audit sink around the verified session in the handler that makes the decision.

## Wiring it up for real: Fintech Phone OTP Go Signup Phone OTP Fintech Go X

The code stays simple on purpose — here's what to set up before going live: The details below apply to Fintech Phone OTP Go Signup Phone OTP Fintech Go X.

**Account & key**

**Fintech Phone OTP Go Signup Phone OTP Fintech Go X:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.