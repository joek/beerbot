package ws_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gorilla/websocket"
	. "github.com/joek/beerbot/web/ws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {

	var u url.URL
	var h *HubImpl
	var s *httptest.Server
	var command chan *BotCommand

	BeforeEach(func() {
		command = make(chan *BotCommand)
		h = NewHub(command)
		go h.Run()

		mux := http.NewServeMux()
		mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { h.ServeWs(w, r) })
		s = httptest.NewServer(mux)
		u = url.URL{Scheme: "ws", Host: s.Listener.Addr().String(), Path: "/ws"}
	})

	AfterEach(func() {
		h.Stop()
		s.Close()
	})

	It("Excepts websocket connections", func() {
		_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		Ω(err).ShouldNot(HaveOccurred())

	})

	It("Sends data", func() {
		c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
		type TestMessage struct {
			MessageId   int
			Description string
		}
		h.Broadcast(TestMessage{
			MessageId:   1002,
			Description: "Test Message"})
		_, p, _ := c.ReadMessage()

		Ω(string(p)).Should(MatchJSON("{\"MessageId\":1002,\"Description\":\"Test Message\"}"))
	})

	Describe("Dailer connected", func() {
		var c *websocket.Conn
		BeforeEach(func() {
			c, _, _ = websocket.DefaultDialer.Dial(u.String(), nil)
		})

		It("Receives data", func() {
			c.WriteMessage(websocket.TextMessage, []byte("{\"motor\": {\"left\": 0.5, \"right\": -0.5}}"))
			Eventually(command).Should(Receive())
		})

		It("Receives motor command", func() {
			c.WriteMessage(websocket.TextMessage, []byte("{\"motor\": {\"left\": 0.5, \"right\": -0.5}}"))
			m := <-command
			Ω(m.Motor.Left).Should(Equal(float32(0.5)))
			Ω(m.Motor.Right).Should(Equal(float32(-0.5)))
		})

		It("Send disconnect event", func() {
			c.Close()
			m := <-command
			Ω(m.Event).Should(Equal("Disconnect"))
		})
	})
})
