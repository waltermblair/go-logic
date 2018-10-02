package logic_test

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/logic/logic"
	"io/ioutil"
	"os"
)

type MockRabbitClientImpl struct {
	URL			string
	thisQueue	string
}

func NewMockRabbitClient(url string, thisQueue string) RabbitClient {
	r := MockRabbitClientImpl{
		URL: 		url,
		thisQueue:	thisQueue,
	}
	return &r
}

func (r *MockRabbitClientImpl) RunConsumer(p Processor) {}
func (r *MockRabbitClientImpl) Publish (m MessageBody, s string) error {
	return errors.New("published")
}
func (r *MockRabbitClientImpl) InitRabbit() {}

var _ = Describe("Core", func() {

	var p 			Processor
	var cfg 		Config
	var msgConfig 	Message
	var msg			Message

	BeforeSuite(func() {
		fileConfig, _ := os.Open("../resources/json/msgConfig.json")
		fileMsg, _ := os.Open("../resources/json/msgTrue.json")
		bytesConfig, _ := ioutil.ReadAll(fileConfig)
		bytesMsg, _ := ioutil.ReadAll(fileMsg)
		json.Unmarshal(bytesConfig, &msgConfig)
		json.Unmarshal(bytesMsg, &msg)
	})

	Describe("Logic component is in up state", func() {
		BeforeEach(func() {
			cfg = Config{
				123,
				"up",
				"buffer",
				[]int{1, 2, 3},
			}
			p = NewProcessor()
			p.ApplyConfig(cfg)
		})

		Describe("Apply Function", func() {
			It("should apply buffer function", func() {
				result := p.ApplyFunction(msg.Body)
				Ω(result).Should(BeTrue())
			})
			It("should apply not function", func() {
				cfg.Function = "not"
				p.ApplyConfig(cfg)
				result := p.ApplyFunction(msg.Body)
				Ω(result).Should(BeFalse())
			})
		})

		Describe("Build Message", func() {
			It("should build message with output", func() {
				result := p.BuildMessage(msg.Body)
				Ω(result.Input[0]).Should(BeTrue())
			})
		})

		// TODO - test I'm sending right message to right nextKey
		Describe("Process Message", func() {

			var mockRabbit RabbitClient

			BeforeEach(func() {
				mockRabbit = NewMockRabbitClient("mock-endpoint", "1")
			})

			It("should mock publish message", func() {
				result := p.Process(msg.Body, mockRabbit)
				Ω(result.Error()).Should(Equal("published"))
			})
		})

	})

})
