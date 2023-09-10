# logtrace

### Is the Pixel8labs library to generate the trace and inject it into a log

#### How to use

```bash
go get github.com/pixel8labs/logtrace
```

```bash
go mod tidy
```

#### Example on using in Cron
```go
package cron

import (
  "github.com/pixel8labs/logtrace/trace"
  log "github.com/pixel8labs/logtrace/log"
  ...
)

func Run(ctx context.Context, app *application.App) {
  h := cron.NewHandler(app)
  
  err := trace.InitTrace(ctx)
  if err != nil {
    panic(err)
  }
  
  trace.Tracer("Health-Squad-Cron")
  
  if _, err := c.AddFunc("@every 1m", func() {
    ctx, span := trace.StartSpan(ctx, "Health-Squad-Cron", "DeleteExpiredRefreshTokens")
    defer span.End()
  
    h.DeleteExpiredRefreshTokens(ctx)
  }); err != nil {
    log.Error(err, log.Fields{}, "Error on adding CRON job for DeleteExpiredRefreshTokens.")
    panic(err)
  }

  ...
}

func (h *Handler) DeleteExpiredRefreshTokens(ctx context.Context) {
  log.Info(ctx, log.Fields{}, "Deleting expired refresh tokens")
  if err := h.RefreshTokenRepo.DeleteExpired(ctx); err != nil {
    log.Error(ctx, err, log.Fields{}, "Error when deleting expired refresh tokens")
  }
  log.Info(ctx, log.Fields{}, "Expired refresh tokens deleted")
}
```

<img width="984" alt="Screenshot 2023-09-08 at 13 20 22" src="https://github.com/pixel8labs/logtrace/assets/79161142/276eec19-3b68-4226-8a4f-6ff2b4c761d4">


#### Example on using in Worker
```go
import (
  "github.com/pixel8labs/logtrace/trace"
  log "github.com/pixel8labs/logtrace/log"
  ...
)

func Run(ctx context.Context, app *application.App) {
  config := config.LoadGoogleCloud()
  
  err := trace.InitTrace(ctx)
  if err != nil {
    panic(err)
  }
  
  trace.Tracer("Health-Squad-Worker")
  
  ctx, span := trace.StartSpan(ctx, "Health-Squad-Worker", "FetchFitbit")
  defer span.End()

  ...

  subscriber := client.Subscription(config.SubscriptionID)
  err = subscriber.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
    ...
  
    ...
  
    log.Info(ctx, log.Fields{"payload": payload}, "user fitbit data fetched successfully")
    msg.Ack()
  })
  if err != nil {
    log.Error(ctx, err, log.Fields{}, "unable to receive messages")
  }
  
  log.Info(ctx, log.Fields{}, "Stopping WORKERS...")
}
```
<img width="959" alt="Screenshot 2023-09-08 at 13 55 20" src="https://github.com/pixel8labs/logtrace/assets/79161142/b70671d6-a720-459d-90e8-0f5ba2ec364d">

#### Example on using in Server

```go
import (
	trace "github.com/pixel8labs/logtrace/trace"
	traceMiddleware "github.com/pixel8labs/logtrace/middleware"
)

func Run(ctx context.Context, app *application.App) {
  port := os.Getenv("PORT")
  if port == "" {
    panic("No PORT Set")
  }
  
  err := trace.InitTrace(ctx)
  if err != nil {
    panic(err)
  }
  e := echo.New()
  
  // Middlewares
  e.Use(
    traceMiddleware.TracingMiddleware("Health-Squad-Server"),
    traceMiddleware.Logger(),
    ...
  )
  ...
}
```
```go
package config

import (
  "context"
  
  "github.com/kelseyhightower/envconfig"
  "golang.org/x/exp/slices"
  
  log "github.com/pixel8labs/logtrace/log"
)

type Twilio struct {
  VerificationSID string   `envconfig:"VERIFICATION_SID" required:"true"`
  Whitelists      []string `envconfig:"WHITELISTS" required:"true"`
}

func (t Twilio) IsWhitelisted(ctx context.Context, phone string) bool {
  whitelisted := slices.Contains(t.Whitelists, phone)
  if whitelisted {
    log.Info(ctx, log.Fields{"phone": phone}, "Phone is whitelisted, skipping Twilio API call")
  }
  return whitelisted
}

func LoadTwilio() Twilio {
  var config Twilio
  envconfig.MustProcess("TWILIO", &config)
  
  return config
}
```
![image](https://github.com/pixel8labs/logtrace/assets/79161142/482661c2-9b8d-40a9-9c23-d43df99bc5a6)
