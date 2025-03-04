package envoy

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	net_url "net/url"

	envoy_bootstrap "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v2"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"

	kuma_dp "github.com/kumahq/kuma/pkg/config/app/kuma-dp"
	"github.com/kumahq/kuma/pkg/core"
	util_proto "github.com/kumahq/kuma/pkg/util/proto"
	"github.com/kumahq/kuma/pkg/xds/bootstrap/types"
)

type remoteBootstrap struct {
	client *http.Client
}

func NewRemoteBootstrapGenerator(client *http.Client) BootstrapConfigFactoryFunc {
	rb := remoteBootstrap{client: client}
	return rb.Generate
}

var (
	log           = core.Log.WithName("dataplane")
	DpNotFoundErr = errors.New("Dataplane entity not found. If you are running on Universal please create a Dataplane entity on kuma-cp before starting kuma-dp. If you are running on Kubernetes, please check the kuma-cp logs to determine why the Dataplane entity could not be created by the automatic sidecar injection.")
)

func (b *remoteBootstrap) Generate(url string, cfg kuma_dp.Config) (proto.Message, error) {
	bootstrapUrl, err := net_url.Parse(url)
	if err != nil {
		return nil, err
	}

	backoff, err := retry.NewConstant(cfg.ControlPlane.BootstrapServer.Retry.Backoff)
	if err != nil {
		return nil, errors.Wrap(err, "could not create retry backoff")
	}
	backoff = retry.WithMaxDuration(cfg.ControlPlane.BootstrapServer.Retry.MaxDuration, backoff)
	var respBytes []byte
	err = retry.Do(context.Background(), backoff, func(ctx context.Context) error {
		log.Info("trying to fetch bootstrap configuration from the Control Plane")
		respBytes, err = b.requestForBootstrap(bootstrapUrl, cfg)
		if err == nil {
			return nil
		}
		switch err {
		case DpNotFoundErr:
			log.Info("Dataplane entity is not yet found in the Control Plane. If you are running on Kubernetes, CP is most likely still in the process of converting Pod to Dataplane. Retrying.", "backoff", cfg.ControlPlane.ApiServer.Retry.Backoff)
		default:
			log.Info("could not fetch bootstrap configuration. Retrying.", "backoff", cfg.ControlPlane.BootstrapServer.Retry.Backoff, "err", err.Error())
		}
		return retry.RetryableError(err)
	})
	if err != nil {
		return nil, err
	}

	bootstrap := envoy_bootstrap.Bootstrap{}
	if err := util_proto.FromYAML(respBytes, &bootstrap); err != nil {
		return nil, errors.Wrap(err, "could not parse the bootstrap configuration")
	}

	return &bootstrap, nil
}

func (b *remoteBootstrap) requestForBootstrap(url *net_url.URL, cfg kuma_dp.Config) ([]byte, error) {
	url.Path = "/bootstrap"
	request := types.BootstrapRequest{
		Mesh: cfg.Dataplane.Mesh,
		Name: cfg.Dataplane.Name,
		// if not set in config, the 0 will be sent which will result in providing default admin port
		// that is set in the control plane bootstrap params
		AdminPort:          cfg.Dataplane.AdminPort.Lowest(),
		DataplaneTokenPath: cfg.DataplaneRuntime.TokenPath,
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal request to json")
	}
	resp, err := b.client.Post(url.String(), "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, errors.Wrap(err, "request to bootstrap server failed")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, DpNotFoundErr
		}
		if resp.StatusCode == http.StatusUnprocessableEntity {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.Errorf("Unable to read the response with status code: %d", resp.StatusCode)
			}
			return nil, errors.Errorf("Error: %s", string(bodyBytes))
		}
		return nil, errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read the body of the response")
	}
	return respBytes, nil
}
