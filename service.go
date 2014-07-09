package imgurl

type Request struct {
	Url string
	Maxwidth,Maxheight int
	Payload interface{}
	Filters []Filter
}

type Response struct {
	Image string
	Payload interface{}
	Tags []interface{}
}


type TranscodeService struct {
	in chan* Request
	out chan* Response
	workerCount,bufferSize int
}

// NewTranscodeService creates worker routines that transcode images in parallel.
// The amount of worker routines is taken from parameters including the
// size of the input and output channel buffer.
func NewTranscodeService(worker,buffersize int) (ret *TranscodeService) {
	in := make(chan *Request,buffersize)
	out := make(chan *Response,buffersize)

	for i:=0;i< worker;i++ {
		go func (in chan *Request) {
			for ;; {
				req := <-in
				img,tags,_ :=Urlify(req.Url,req.Maxwidth,req.Maxheight,req.Filters...);
				resp := &Response{Image: img, Payload: req.Payload, Tags: tags}
				out <- resp
			}
		}(in)
	}

	return &TranscodeService{in: in, out: out,workerCount: worker,bufferSize: buffersize}
}

func (t *TranscodeService) Get() (*Response) {
	return <- t.out
}

func (t *TranscodeService) Push(r *Request) {
	t.in <- r
}

func (t *TranscodeService) Full() bool {
	return len(t.in) >= t.bufferSize
}

func (t *TranscodeService) Ready() bool {
	return len(t.out) > 0
}
