package main

import (
	"fmt"

	"github.com/blockloop/scan"
)

// below functions will be enough to facilitate a graph traversal
func getActorDetails(nconst string) (Actor, error) {
	var actor Actor

	rows, err := db.Query("SELECT * FROM actors WHERE nconst = $1", nconst)
	if err != nil {
		return actor, err
	}

	err = scan.Row(&actor, rows)
	if err != nil {
		return actor, err
	}

	return actor, nil
}

func actorsInTitle(tconst string) ([]Actor, error) {
	var actors []Actor

	rows, err := db.Query("select actors.* from actors join title_actors on actors.nconst = title_actors.nconst where title_actors.tconst = $1", tconst)
	if err != nil {
		return nil, err
	}

	err = scan.Rows(&actors, rows)
	if err != nil {
		return nil, err
	}

	return actors, nil
}

func titlesWithActor(nconst string) ([]Title, error) {
	// get all titles that the actor has been in
	var titles []Title

	rows, err := db.Query("select titles.* from titles join title_actors on titles.tconst = title_actors.tconst where title_actors.nconst = $1", nconst)
	if err != nil {
		return titles, err
	}

	err = scan.Rows(&titles, rows)
	if err != nil {
		return titles, err
	}

	return titles, nil
}

func getTitleDetails(tconst string) (Title, error) {
	var title Title

	rows, err := db.Query("SELECT * FROM titles WHERE tconst = $1", tconst)
	if err != nil {
		return title, err
	}

	err = scan.Row(&title, rows)
	if err != nil {
		return title, err
	}

	return title, nil
}

type edge struct {
	Nconst string `db:"nconst"`
	Tconst string `db:"tconst"`
}

func actorsWithActor(nconst string) []edge {

	// get all actors that the actor has been in a title with
	rows, err := db.Query("select actors.nconst, title_actors.tconst from actors join title_actors on actors.nconst = title_actors.nconst where title_actors.tconst in (select tconst from title_actors where nconst = $1)", nconst)
	handleError(err)

	var edges []edge

	err = scan.Rows(&edges, rows)
	handleError(err)

	return edges
}

func searchGraph(nconstStart string, nconstEnd string) map[string]string {
	// BFS
	// start with nconstStart
	// get all actors that the actor has been in a title actorsWithActor
	// if nconstEnd is in the list return the path
	// else add all actors to the queue
	// repeat until queue is empty
	q := []string{nconstStart}
	parent := map[string]string{}
	visited := map[string]bool{nconstStart: true}
	for len(q) != 0 {
		v := q[0] // dequeue
		q = q[1:]
		//fmt.Println("Visiting", v)
		if v == nconstEnd {
			// return path
			fmt.Println("found")
			break
		}
		for _, edge := range actorsWithActor(v) {
			//fmt.Printf("checking %s from %s \n", edge.Nconst, v)
			if !visited[edge.Nconst] {
				visited[edge.Nconst] = true
				q = append(q, edge.Nconst)
				parent[edge.Nconst] = v
			}
		}
	}
	//fmt.Println(parent)
	return parent
}

func findPathRecurse(nconstStart string, nconstEnd string, parents map[string]string) {
	if (nconstStart == nconstEnd) || (nconstEnd == "") {
		fmt.Println(nconstStart)
	} else {
		findPathRecurse(nconstStart, parents[nconstEnd], parents)
		fmt.Println("adding", nconstEnd)
	}
}

func findPath(nconstStart string, nconstEnd string) []string {
	path := []string{nconstStart}
	parents := searchGraph(nconstStart, nconstEnd)
	findPathRecurse(nconstStart, nconstEnd, parents)
	return path
}
