package repository

import (
	"math/rand/v2"
)

type InMemoryRepository struct {
	data []string
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: []string{
			"We can know only that we know nothing. And that is the highest degree of human wisdom.",
			"The strongest of all warriors are these two — Time and Patience.",
			"All, everything that I understand, I understand only because I love. Everything is, in a sense, meaningless without love.",
			"If everyone fought for their own convictions, there would be no war.",
			"There is no greatness where there is no simplicity, goodness, and truth.",
			"In the case of a conflict between two countries, the truth is always on the side of the country which is the most powerful.",
			"The two most powerful warriors are patience and time.",
			"A man is never as happy as when he is in the process of achieving something, of making progress.",
			"True love is not about being happy together but about helping each other become better.",
			"Everything depends on the question: Can we forgive?",
			"It’s not the life that matters, but the courage you bring to it.",
			"Human beings are made for love, and love is the only thing that keeps them alive.",
			"People who fight for a good cause have always been misunderstood.",
			"Happiness does not depend on outward things. It depends on the way we look at things.",
			"Man’s real nature, as I see it, is to suffer, but to suffer for the sake of some greater good.",
			"You can be sure that the truth of a matter is always somewhere in between the two extremes.",
			"War is the continuation of politics by other means.",
			"The greatest happiness in life is to have love and to be loved in return.",
			"To be good and to be loved are two different things, but love is always what makes us better.",
			"Only those who can see the invisible can do the impossible.",
		},
	}
}

func (r *InMemoryRepository) GetWoWQuote() string {
	//nolint:gosec
	return r.data[rand.IntN(len(r.data))]
}
