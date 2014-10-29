package game

type cardrank uint8

func (cr cardrank) Next() cardrank {
	if uint8(cr) < 9 {
		return cr + 1
	}
	return cardrank(0)
}

func (cr cardrank) Prev() cardrank {
	if uint8(cr) > 0 {
		return cr - 1
	}
	return cardrank(9)
}

type card struct {
	rank cardrank
}

func NewCard(rank uint8) *card {
	c := card{}
	c.rank = cardrank(rank)
	return &c
}

func (c *card) NextTo(targetCard *card) bool {
	if c.rank == targetCard.rank.Next() || c.rank == targetCard.rank.Prev() {
		return true
	}
	return false
}
