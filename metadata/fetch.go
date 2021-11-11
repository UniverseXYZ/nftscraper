package metadata

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/universexyz/nftscraper/conf"
)

type ctxKeyHttpClient struct{}

const (
	fetchRetryInitialInterval = 200 * time.Millisecond
	fetchRetryMaxInterval     = 30 * time.Second
	fetchRetryMaxElapsedTime  = 3 * time.Minute
)

func ContextWithHttpClient(ctx context.Context, httpClient *http.Client) context.Context {
	return context.WithValue(ctx, ctxKeyHttpClient{}, httpClient)
}

func ReadExternalResource(ctx context.Context, resourceURI string, dst io.Writer) (int64, error) {
	var err error

	resourceURI = strings.TrimSpace(resourceURI)
	if len(resourceURI) == 0 {
		return 0, nil
	}

	resourceURI, err = useIPFSGwIfNeeded(resourceURI)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	hc := http.DefaultClient
	if v, ok := ctx.Value(ctxKeyHttpClient{}).(*http.Client); ok && hc != nil {
		hc = v
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, resourceURI, nil)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = fetchRetryInitialInterval
	retry.MaxInterval = fetchRetryMaxInterval
	retry.MaxElapsedTime = fetchRetryMaxElapsedTime

	var resp *http.Response

	err = backoff.Retry(func() (callErr error) {
		resp, callErr = hc.Do(req)
		if callErr != nil {
			// check if err implements Timeout()
			var timeoutErr interface{ Timeout() bool }
			timeout := errors.As(callErr, &timeoutErr) && timeoutErr.Timeout()

			// check if err implements Temporary()
			var tempErr interface{ Temporary() bool }
			temporary := errors.As(callErr, &tempErr) && tempErr.Temporary()

			// timeout means it is temporary
			temporary = temporary || timeout

			if !temporary {
				return backoff.Permanent(err)
			}

			return callErr
		}

		return nil
	}, backoff.WithContext(retry, ctx))

	if err != nil {
		return 0, errors.WithStack(err)
	}

	defer resp.Body.Close()

	n, err := io.Copy(dst, resp.Body)
	if err != nil {
		return n, errors.WithStack(err)
	}

	return n, nil
}

func useIPFSGwIfNeeded(resourceURI string) (string, error) {
	resURL, err := url.Parse(resourceURI)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if resURL.Scheme != "ipfs" {
		return resourceURI, nil
	}

	cfg := conf.Conf()

	if cfg.IPFSHost == "" {
		return "", errors.New("IPFS hostname is not configured")
	}

	ipfsGw, err := url.Parse(cfg.IPFSHost)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse IPFS host configuration: `%s` - %s", cfg.IPFSHost, err.Error())
	}

	ipfsGw.Path = "/api/v0/cat"

	ipfsGw.RawQuery = url.Values{
		"arg": []string{path.Join(resURL.Host, resURL.Path)},
	}.Encode()

	if cfg.IPFSUser != "" || cfg.IPFSPass != "" {
		ipfsGw.User = url.UserPassword(cfg.IPFSUser, cfg.IPFSPass)
	}

	return ipfsGw.String(), nil
}
