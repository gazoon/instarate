package competition

import "math"

const (
	significanceCoeff = 16
)

func recalculateEloRating(winnerRating, loserRating int) (int, int) {
	newWinnerRating := calculateNewRating(winnerRating, loserRating, true)
	newLoserRating := calculateNewRating(loserRating, winnerRating, false)
	return newWinnerRating, newLoserRating
}

func calculateNewRating(rating1, rating2 int, isFirstWinner bool) int {
	ea := expectedProbability(rating1, rating2)
	var actualResult float64
	if isFirstWinner {
		actualResult = 1
	} else {
		actualResult = 0
	}
	return rating1 + int(math.Round(significanceCoeff*(actualResult-ea)))

}

func expectedProbability(rating1, rating2 int) float64 {
	return 1 / (1 + math.Pow(10, float64(rating1-rating2)/400))
}
