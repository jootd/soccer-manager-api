package userbus

// type TeamPlayer struct {
// 	Team
// 	players []Player
// 	Value   float64
// }

// type QueryTeamPlayer struct {
// }

// type TeamPlayerStorer interface {
// 	Query(ctx context.Context, query QueryPlayer) ([]TeamPlayer, error)
// }

// type TeamPlayeruserbus struct {
// 	store TeamPlayerStorer
// }

// func (tb *TeamPlayeruserbus) Value(ctx, query QueryTeamPlayer) (float64, error) {
// 	teamplayers, err := tb.store.Query(ctx, query)
// 	if err != nil {
// 		return
// 	}
// 	for _, v := range teamplayers {

// 	}

// 	return

// }

// it also might be postgres view
// but
// func (tp *TeamPlayeruserbus) GetTeamPlayers(ctx context.Context, teamId int) (TeamPlayer, error) {
// teams, ok := tp.store.GetTeamsBy(ctx, QueryTeam{ID: &teamId})
// if !ok {
// 	return TeamPlayer{}, errors.New("")
// }
// team := teams[0]

// players, ok := tp.store.GetPlayersBy(ctx, QueryPlayer{
// 	teamId: &teamId,
// })

// if !ok {
// 	return TeamPlayer{}, errors.New("")
// }

// return TeamPlayer{
// 	players: players,
// 	Team: Team{
// 		ID:      team.ID,
// 		Name:    team.Name,
// 		Country: team.Country,
// 	},
// }, nil
// }
