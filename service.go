package imgurl

type Request struct {
	Url string
	Maxwidth,Maxheight int
	Payload interface{}
}

type Response struct {
	Image string
	Payload interface{}
}


type TranscodeService struct {
	In chan* Request
	Out chan* Response
}

func NewTranscodeService(worker int) (ret *TranscodeService) {
	in := make(chan *Request,10)
	out := make(chan *Response,10)

	for i:=0;i< worker;i++ {
		go func (in chan *Request) {
			for ;; {
				req := <-in
				img,_ :=Urlify(req.Url,req.Maxwidth,req.Maxheight);
				resp := &Response{Image: img, Payload: req.Payload}
				out <- resp
			}
		}(in)
	}

	return &TranscodeService{In: in, Out: out}
}
