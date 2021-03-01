package main

import (
	"time"

	"design-pattern/observer/goodexample/example2/player"
)

func main() {
	rankProcessor := player.LevelUpRankProcessor{
		Type: player.TypeLevelUp,
	}
	rewardProcessor := player.LevelUpRewardProcessor{
		Type: player.TypeLevelUp,
	}
	announceProcessor := player.LevelUpAnnounceProcessor{
		Type: player.TypeLevelUp,
	}
	obs := player.GetObs()
	obs.AddProcessor(rankProcessor)
	obs.AddProcessor(rewardProcessor)
	obs.AddProcessor(announceProcessor)

	p := player.NewPlayer("xie4ever")
	_ = p.Attack()
	time.Sleep(time.Second)
	_ = p.Attack()
	time.Sleep(time.Second)
	_ = p.Attack()
}
