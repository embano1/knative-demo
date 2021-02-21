package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
)

const (
	listen  = ":8080"
	version = "v1"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
}

func run(args []string, stdout io.Writer) error {
	ctx := signals.NewContext()

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	srv := http.Server{
		Addr:    listen,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		logging.FromContext(ctx).Infof("shutting down")
		_ = srv.Shutdown(ctx)
	}()

	logging.FromContext(ctx).Infof("app version: %s", version)
	logging.FromContext(ctx).Infof("listening on %s", listen)
	return srv.ListenAndServe()
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	name := "stranger"

	if v := req.URL.Query().Get("name"); v != "" {
		name = v
	}

	greeting := fmt.Sprintf("Hello %s!", name)
	_, _ = w.Write([]byte(greeting))
}
