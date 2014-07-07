package imgurl

type Request struct {
	url string
	maxwidth,maxheight int
	payload interface{}
}

type Response struct {
	image string
	payload interface{}
}


type TranscodeService struct {
	in chan* Request
	out chan* Response
}

func NewTranscodeService(worker int) (ret *TranscodeService) {
	in := make(chan *Request,10)
	out := make(chan *Response,10)

	for i:=0;i< worker;i++ {
		go func (in chan *Request) {
			for ;; {
				req := <-in
				img,_ :=Urlify(req.url,req.maxwidth,req.maxheight);
				resp := &Response{image: img, payload: req.payload}
				out <- resp
			}
		}(in)
	}

	return &TranscodeService{in: in, out: out}
}
