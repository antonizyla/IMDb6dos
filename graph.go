package main

// below functions will be enough to facilitate a graph traversal

func getActor(nconst string) Actor {
	var actor Actor
	//res := db.First(&actor, "nconst = ?", nconst)
	//if res.Error != nil {
//		fmt.Println(res.Error)
//	}
	return actor
}

func getConnectedActors(tconst string) []Actor {
	var actors []Actor
//	res := db.Joins("JOIN movie_actors ON movie_actors.nconst = actors.nconst").Where("tconst = ?", tconst).Find(&actors)
//	if res.Error != nil {
//		fmt.Println(res.Error)
//	}
	return actors
}

func getConnectedTitles(nconst string) []Title {
	// get all titles that the actor has been in
	return []Title{}
}

func getTitle(tconst string) Title {
	// get title by tconst
	return Title{}
}
