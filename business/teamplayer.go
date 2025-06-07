package business

type TeamPlayer struct {
	Team
	players []Player
}

type TeamPlayerStorer interface {
	TeamStorer
	PlayerStorer
}

type TeamPlayerBusiness struct {
	store TeamPlayerStorer
}

// it also might be postgres view
// but
// func (tp *TeamPlayerBusiness) GetTeamPlayers(ctx context.Context, teamId int) (TeamPlayer, error) {
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
