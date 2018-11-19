package logic_test

import (
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/waltermblair/logic/logic"
	"io/ioutil"
	"os"
	"strconv"
)

type MockRabbitClientImpl struct {
	config 		Config
}

func NewMockRabbitClient(cfg Config) RabbitClient {
	r := MockRabbitClientImpl{
		cfg,
	}
	return &r
}

func (r *MockRabbitClientImpl) RunConsumer(p Processor) {}
func (r *MockRabbitClientImpl) Publish (m MessageBody, s string) error {
	output := strconv.FormatBool(m.Input[0])
	return errors.New("next-key: " + s + " output: " + output)
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

	Describe("Logic component is in up state and expects one input", func() {
		BeforeEach(func() {
			cfg = Config{
				123,
				"up",
				"buffer",
				1,
				[]int{1, 2, 3},
			}
			p = NewProcessor()
			p.ApplyConfig(cfg)
		})

		Describe("Apply Config", func() {
			It("should apply config", func() {
				Ω(len(p.GetConfig().NextKeys)).Should(Equal(3))
				p.ApplyConfig(Config{456, "", "", 1, []int{1,2}})
				Ω(len(p.GetConfig().NextKeys)).Should(Equal(2))
			})
		})

		Describe("Apply Function", func() {
			It("should apply buffer function", func() {
				p.ApplyFunction(msg.Body)
				Ω(p.GetOutput()).Should(BeTrue())
			})
			It("should apply not function", func() {
				cfg.Function = "not"
				p.ApplyConfig(cfg)
				p.ApplyFunction(msg.Body)
				Ω(p.GetOutput()).Should(BeFalse())
			})
		})

		Describe("Build Message", func() {
			It("should build message with output true", func() {
				p.ApplyFunction(msg.Body)
				result := p.BuildMessage()
				Ω(result.Input[0]).Should(BeTrue())
			})
		})

		Describe("Process Message", func() {

			var mockRabbit RabbitClient
			var lastKey    string
			var output	   string

			BeforeEach(func() {
				mockRabbit = NewMockRabbitClient(cfg)
				lastKey = strconv.Itoa(cfg.NextKeys[len(cfg.NextKeys)-1])
			})

			It("should apply config and mock publish output messages", func() {
				p.ApplyConfig(msgConfig.Body.Configs[0])
				p.ApplyFunction(msgConfig.Body)
				output = strconv.FormatBool(p.GetOutput())
				result := p.Process(msgConfig.Body, mockRabbit)
				Ω(result.Error()).Should(Equal("next-key: " + lastKey + " output: " + output))
			})
			It("should mock publish output messages", func() {
				p.ApplyFunction(msg.Body)
				output = strconv.FormatBool(p.GetOutput())
				result := p.Process(msg.Body, mockRabbit)
				Ω(result.Error()).Should(Equal("next-key: " + lastKey + " output: " + output))
			})
		})
	})

	Describe("Logic component is in up state and expects multiple inputs", func() {
		BeforeEach(func() {
			cfg = Config{
				123,
				"up",
				"and",
				2,
				[]int{1, 2, 3},
			}
			p = NewProcessor()
			p.ApplyConfig(cfg)
		})

		Describe("Process Message", func() {

			var mockRabbit RabbitClient
			var lastKey    string
			var output	   string

			BeforeEach(func() {
				mockRabbit = NewMockRabbitClient(cfg)
				lastKey = strconv.Itoa(cfg.NextKeys[len(cfg.NextKeys)-1])
			})

			It("should mock publish output messages", func() {
				result := p.Process(msg.Body, mockRabbit)
				result = p.Process(msg.Body, mockRabbit)
				output = strconv.FormatBool(p.GetOutput())
				Ω(result.Error()).Should(Equal("next-key: " + lastKey + " output: " + output))
			})
		})
	})
})
