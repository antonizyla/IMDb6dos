package main

import "fmt"

func getActor(nconst string) Actor {
	var actor Actor
	res := db.First(&actor, "nconst = ?", nconst)
	if res.Error != nil {
		fmt.Println(res.Error)
	}
	return actor
}

func getActors(tconst string) []Actor {
    var actors []Actor
    res := db.Joins("JOIN movie_actors ON movie_actors.nconst = actors.nconst").Where("tconst = ?", tconst).Find(&actors)
    if res.Error != nil {
        fmt.Println(res.Error)
    }
    return actors
}
