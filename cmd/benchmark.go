package cmd

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	adapter_sdk "github.com/BrobridgeOrg/gravity-sdk/v2/adapter"
	product_sdk "github.com/BrobridgeOrg/gravity-sdk/v2/product"
	subscriber_sdk "github.com/BrobridgeOrg/gravity-sdk/v2/subscriber"
	gravity_sdk_types_product_event "github.com/BrobridgeOrg/gravity-sdk/v2/types/product_event"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

const (
	benchmarkProductStream = "GVT_%s_DP_%s"
)

func init() {

	RootCmd.AddCommand(domainBenchmarkCmd)
}

func assertDomainStream(js nats.JetStreamContext, streamName string, subject string) error {

	fmt.Printf("Check domain stream: %s\n", streamName)

	// Check if the stream already exists
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		return err
	}

	if stream == nil {

		_, err = js.AddStream(&nats.StreamConfig{
			Name:        streamName,
			Description: "Gravity domain event store",
			Subjects: []string{
				subject,
			},
		})

		fmt.Printf("Created domain stream: %s\n", streamName)

		if err != nil {
			return err
		}
	}

	return nil
}

var domainBenchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "measuring performance",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := runProductCmd(runDomainBenchmarkCmd, cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func runDomainBenchmarkCmd(cctx *ProductCommandContext) error {

	js, err := cctx.Connector.GetClient().GetJetStream()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	// Assert domain stream
	eventStream := fmt.Sprintf(domainEventStream, cctx.Connector.GetDomain())
	eventSubject := fmt.Sprintf(domainEventSubject, cctx.Connector.GetDomain(), "*")

	err = assertDomainStream(js, eventStream, eventSubject)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	// Preparing product for benchmarking
	productStreamName := fmt.Sprintf(benchmarkProductStream, cctx.Connector.GetDomain(), "gvt_benchmark")
	setting := product_sdk.ProductSetting{}
	setting.Name = "gvt_benchmark"
	setting.Description = "Gravity benchmarking"
	setting.Enabled = true
	setting.Stream = productStreamName
	setting.Schema = map[string]interface{}{
		"id": map[string]interface{}{
			"type": "uint",
		},
		"ts": map[string]interface{}{
			"type": "uint",
		},
	}

	// Preparing rule
	rule := product_sdk.NewRule()
	id, _ := uuid.NewUUID()
	rule.ID = id.String()

	rule.Name = "benchmark"
	rule.Product = setting.Name
	rule.UpdatedAt = time.Now()
	rule.CreatedAt = time.Now()
	rule.Event = "gvt_benchmark"
	rule.Method = "create"
	rule.PrimaryKey = []string{"id"}
	rule.Description = "benchmark"
	rule.Enabled = true
	rule.SchemaConfig = map[string]interface{}{
		"id": map[string]interface{}{
			"type": "uint",
		},
		"ts": map[string]interface{}{
			"type": "uint",
		},
	}
	rule.HandlerConfig = &product_sdk.HandlerConfig{
		Type: "script",
		Script: `
return {
	id: source.id,
	ts: source.ts
}
		`,
	}

	setting.Rules = map[string]*product_sdk.Rule{
		rule.Name: rule,
	}

	// Check whether product exists or not
	_, err = cctx.Product.GetClient().GetProduct(setting.Name)
	if err != nil {
		if err == product_sdk.ErrProductNotFound {

			// Create temporary product
			fmt.Printf("Creating product: %s\n", setting.Name)
			_, err = cctx.Product.GetClient().CreateProduct(&setting)
			if err != nil {
				cctx.Cmd.SilenceUsage = true
				return err
			}
		} else {
			return err
		}
	}

	// TODO: Better way to wait for product
	time.Sleep(time.Second * 3)

	doBenchmark(cctx)

	// Delete temporary product
	err = cctx.Product.GetClient().DeleteProduct(setting.Name)
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		fmt.Printf("Faield to deleting product: %s\n", setting.Name)
		return err
	}

	return nil
}

func doBenchmark(cctx *ProductCommandContext) error {

	var wg sync.WaitGroup
	var counter uint64
	var total uint64 = 10000
	var minLatency int64 = -1
	var maxLatency int64 = 0
	var sentTime time.Time

	productName = "gvt_benchmark"

	fmt.Printf("Subscribing to product: %s\n", productName)

	// Initializing gravity subscriber
	opts := subscriber_sdk.NewOptions()
	opts.Domain = cctx.Connector.GetDomain()
	s := subscriber_sdk.NewSubscriberWithClient("", cctx.Connector.GetClient(), opts)
	sub, err := s.Subscribe(productName, func(msg *nats.Msg) {

		var pe gravity_sdk_types_product_event.ProductEvent

		err := proto.Unmarshal(msg.Data, &pe)
		if err != nil {
			fmt.Printf("Failed to parsing product event: %v", err)
			msg.Ack()
			return
		}

		r, err := pe.GetContent()
		if err != nil {
			fmt.Printf("Failed to parsing content: %v", err)
			msg.Ack()
			return
		}

		var ts uint64
		for _, field := range r.Payload.Map.Fields {
			if field.Name == "ts" {
				ts = field.Value.GetData().(uint64)
				break
			}
		}

		//		md, _ := msg.Metadata()
		//		ts = uint64(md.Timestamp.UnixNano())

		end := time.Now().UnixNano()
		d := end - int64(ts)
		if minLatency == -1 {
			minLatency = d
		} else if minLatency > d {
			minLatency = d
		}

		if maxLatency < d {
			maxLatency = d
		}

		counter = atomic.AddUint64((*uint64)(&counter), 1)
		if counter == total {
			wg.Done()
		}
	}, subscriber_sdk.Partition(-1), subscriber_sdk.DeliverNew())
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		fmt.Println(err)
		return err
	}

	defer sub.Close()
	wg.Add(1)

	// Connect to gravity network
	publisher, err := cctx.Connector.CreateClient()
	if err != nil {
		cctx.Cmd.SilenceUsage = true
		return err
	}

	// Initializing adapter connector
	aopts := adapter_sdk.NewOptions()
	aopts.Domain = cctx.Connector.GetDomain()
	ac := adapter_sdk.NewAdapterConnectorWithClient(publisher, aopts)

	fmt.Println("benchmarking...")

	var startTime time.Time
	go func() {

		wg.Add(1)
		defer wg.Done()

		startTime = time.Now()
		for i := uint64(1); i <= total; i++ {
			payload := []byte(fmt.Sprintf(`{"id":%d,"ts":%d}`, i, time.Now().UnixNano()))
			_, err = ac.PublishAsync(productName, payload, nil)
			if err != nil {
				cctx.Cmd.SilenceUsage = true
				fmt.Println(err)
				os.Exit(1)
				return
			}
		}

		<-ac.PublishAsyncComplete()

		sentTime = time.Now()
	}()

	wg.Wait()

	endTime := time.Now()
	elapsed := endTime.UnixNano() - startTime.UnixNano()
	elapsedOfSent := sentTime.UnixNano() - startTime.UnixNano()

	fmt.Printf("Total number of messages: %d\n", total)
	fmt.Printf("Total execution time: %s\n", endTime.Sub(startTime))
	fmt.Printf("Duration time of publishing: %s\n", endTime.Sub(sentTime))
	fmt.Printf("Thoughput: %f msg/s\n", (float64(total) / float64(elapsed/1000000000)))

	if maxLatency > elapsed-elapsedOfSent {
		maxLatency -= elapsed - elapsedOfSent
	}

	fmt.Printf("Latency (Min): %dms\n", minLatency/1000000)
	fmt.Printf("Latency (Max): %dms\n", maxLatency/1000000)

	return nil
}
