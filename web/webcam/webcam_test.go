package webcam_test

// All stolen at https://github.com/robdimsdale/garagepi

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/joek/beerbot/web/webcam"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var (
	fakeResponseWriter *httptest.ResponseRecorder

	dummyRequest *http.Request
	w            webcam.Handler
)

var _ = Describe("Webcam", func() {
	var server *ghttp.Server

	BeforeEach(func() {
		server := ghttp.NewServer()
		webcamURL := server.URL()
		parsedURL, err := url.Parse(webcamURL)

		Ω(err).NotTo(HaveOccurred())

		fakeResponseWriter = httptest.NewRecorder()

		w = webcam.NewHandler(parsedURL.Host)

		dummyRequest = new(http.Request)
		dummyRequest.URL = &url.URL{}
		dummyRequest.Header = http.Header{}
	})

	AfterEach(func() {
		server.Close()
	})

	It("should make a request to fetch the image", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/", "action=stream"),
			),
		)

		w.Handle(fakeResponseWriter, dummyRequest)
		Ω(server.ReceivedRequests()).Should(HaveLen(1))
	})

})
