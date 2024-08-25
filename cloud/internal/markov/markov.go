package markov

const (
	beforeOffice = "https://raw.githubusercontent.com/MarkHershey/CompleteTrumpTweetsArchive/master/data/realDonaldTrump_bf_office.csv"
	inOffice     = "https://raw.githubusercontent.com/MarkHershey/CompleteTrumpTweetsArchive/master/data/realDonaldTrump_in_office.csv"
	order        = 2
)

type MarkovChain struct {
	chain map[string][]string
}

func NewMarkovChain() *MarkovChain {
	return &MarkovChain{
		chain: make(map[string][]string),
	}
}

func (mc *MarkovChain) Train(words []string, order int) {
	for i := 0; i < len(words)-order; i++ {
		key := words[i]
		for j := 1; j < order; j++ {
			key += " " + words[i+j]
		}
		if _, ok := mc.chain[key]; !ok {
			mc.chain[key] = make([]string, 0)
		}
		mc.chain[key] = append(mc.chain[key], words[i+order])
	}
}

func (mc *MarkovChain) Generate(seed string, order int, length int) string {
	words := make([]string, 0)
	words = append(words, seed)
	for i := 0; i < length; i++ {
		key := words[len(words)-order]
		for j := 1; j < order; j++ {
			key += " " + words[len(words)-order+j]
		}
		if _, ok := mc.chain[key]; !ok {
			break
		}
		words = append(words, mc.chain[key][0])
	}
	return join(words)
}

func join(words []string) string {
	str := ""
	for i, word := range words {
		if i == 0 {
			str += word
		} else {
			str += " " + word
		}
	}
	return str
}
