package strategy

type Ocr struct {
	context  *OcrContext
	strategy OcrStrategy
}

type OcrContext struct {
	ImgUrl string
}

func NewOcr(imgUrl string, strategy OcrStrategy) *Ocr {
	return &Ocr{
		context: &OcrContext{
			ImgUrl: imgUrl,
		},
		strategy: strategy,
	}
}

func (p *Ocr) Ocr() (interface{}, error) {
	return p.strategy.Ocr(p.context)
}

type OcrStrategy interface {
	Ocr(*OcrContext) (interface{}, error)
}
