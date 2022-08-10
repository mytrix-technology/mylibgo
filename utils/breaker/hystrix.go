package breaker

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

// Client to do http requests with
var Client *http.Client

// RETRIES is the number of retries to do in the retrier.
var retries = 3

func CallUsingCircuitBreaker(breakername string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1) // declare the channel where the hystrix goroutine will put success responses.

	errors := hystrix.Go(breakername, // pass the name of the circuit breaker as first parameter.

		// 2nd parameter, the inlined func to run inside the breaker.
		func() error {
			// create the request. omitted err handling for brevity
			req, _ := http.NewRequest(method, url, nil)

			// for hystrix, forward the err from the retrier. it's nil if successful.
			return CallWithRetries(req, output)
		},

		// 3rd parameter, the fallback func. in this case, we just do a bit of logging and return the error.
		func(err error) error {
			logrus.Errorf("in fallback function for breaker %v, error: %v", breakername, err.Error())
			circuit, _, _ := hystrix.GetCircuit(breakername)
			logrus.Errorf("circuit state is: %v", circuit.IsOpen())
			return err
		})

	// response and error handling. if the call was successful, the output channel gets the response. otherwise,
	// the errors channel gives us the error.
	select {
	case out := <-output:
		logrus.Debugf("call in breaker %v successful", breakername)
		return out, nil

	case err := <-errors:
		return nil, err
	}
}

func CallWithRetries(req *http.Request, output chan []byte) error {

	// create a retrier with constant backoff, retries number of attempts (3) with a 100ms sleep between retries.
	r := retrier.New(retrier.ConstantBackoff(retries, 100*time.Millisecond), nil)

	// this counter is just for getting some logging for showcasing, remove in production code.
	attempt := 0

	// retrier works similar to hystrix, we pass the actual work (doing the http request) in a func.
	err := r.Run(func() error {
		attempt++

		// do http request and handle response. if successful, pass resp.body over output channel,
		// otherwise, do a bit of error logging and return the err.
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			responsebody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				output <- responsebody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("status was %v", resp.StatusCode)
		}

		logrus.Errorf("retrier failed, attempt %v", attempt)
		return err
	})
	return err
}

// PerformHTTPRequestCircuitBreaker performs the supplied http.Request within a circuit breaker.
func PerformHTTPRequestCircuitBreaker(ctx context.Context, breakerName string, req *http.Request) ([]byte, error) {
	//func PerformHTTPRequestCircuitBreaker(ctx context.Context, breakerName string, req *http.Request) ([]byte, error) {
	output := make(chan []byte, 1)
	errors := hystrix.Go(breakerName, func() error {
		err := CallWithRetries(req, output)
		return err // For hystrix, forward the err from the retrier. It's nil if OK.
	}, func(err error) error {
		logrus.Errorf("In fallback function for breaker %v, error: %v", breakerName, err.Error())
		return err
	})

	select {
	case out := <-output:
		logrus.Debugf("Call in breaker %v successful", breakerName)
		return out, nil

	case err := <-errors:
		logrus.Errorf("Got error on channel in breaker %v. Msg: %v", breakerName, err.Error())
		return nil, err
	}
}

// ConfigureHystrix sets up hystrix circuit breakers.
//func ConfigureHystrix(commands []string, amqpClient messaging.IMessagingClient) {
//
//	for _, command := range commands {
//		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
//			Timeout:                resolveProperty(command, "Timeout"),
//			MaxConcurrentRequests:  resolveProperty(command, "MaxConcurrentRequests"),
//			ErrorPercentThreshold:  resolveProperty(command, "ErrorPercentThreshold"),
//			RequestVolumeThreshold: resolveProperty(command, "RequestVolumeThreshold"),
//			SleepWindow:            resolveProperty(command, "SleepWindow"),
//		})
//		logrus.Printf("Circuit %v settings: %v", command, hystrix.GetCircuitSettings()[command])
//	}
//
//	hystrixStreamHandler := hystrix.NewStreamHandler()
//	hystrixStreamHandler.Start()
//	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
//	logrus.Infoln("Launched hystrixStreamHandler at 8181")
//
//	// Publish presence on RabbitMQ
//	publishDiscoveryToken(amqpClient)
//}

// Deregister publishes a Deregister token to Hystrix/Turbine
//func Deregister(amqpClient messaging.IMessagingClient) {
//	ip, err := util.ResolveIPFromHostsFile()
//	if err != nil {
//		ip = util.GetIPWithPrefix("10.0.")
//	}
//	token := DiscoveryToken{
//		State:   "DOWN",
//		Address: ip,
//	}
//	bytes, _ := json.Marshal(token)
//	amqpClient.PublishOnQueue(bytes, "discovery")
//	logrus.Infoln("Sent deregistration token over SpringCloudBus")
//}

//func publishDiscoveryToken(amqpClient messaging.IMessagingClient) {
//	ip, err := util.ResolveIPFromHostsFile()
//	if err != nil {
//		ip = util.GetIPWithPrefix("10.0.")
//	}
//	token := DiscoveryToken{
//		State:   "UP",
//		Address: ip,
//	}
//	bytes, _ := json.Marshal(token)
//	go func() {
//		for {
//			amqpClient.PublishOnQueue(bytes, "discovery")
//			amqpClient.PublishOnQueue(bytes, "discovery")
//			time.Sleep(time.Second * 30)
//		}
//	}()
//}

func resolveProperty(command string, prop string) int {
	if viper.IsSet("hystrix.command." + command + "." + prop) {
		return viper.GetInt("hystrix.command." + command + "." + prop)
	}
	return getDefaultHystrixConfigPropertyValue(prop)

}
func getDefaultHystrixConfigPropertyValue(prop string) int {
	switch prop {
	case "Timeout":
		return 1000 //hystrix.DefaultTimeout
	case "MaxConcurrentRequests":
		return 200 //hystrix.DefaultMaxConcurrent
	case "RequestVolumeThreshold":
		return hystrix.DefaultVolumeThreshold
	case "SleepWindow":
		return hystrix.DefaultSleepWindow
	case "ErrorPercentThreshold":
		return hystrix.DefaultErrorPercentThreshold
	}
	panic("Got unknown hystrix property: " + prop + ". Panicing!")
}

// DiscoveryToken defines a struct for transmitting the state of a hystrix stream producer.
type DiscoveryToken struct {
	State   string `json:"state"` // UP, RUNNING, DOWN ??
	Address string `json:"address"`
}
