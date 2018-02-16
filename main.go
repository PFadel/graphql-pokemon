package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/levigross/grequests"
)

type pokemonResponse struct {
	Forms []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"forms"`
	Abilities []struct {
		Slot     int  `json:"slot"`
		IsHidden bool `json:"is_hidden"`
		Ability  struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		Stat struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"stat"`
		Effort   int `json:"effort"`
		BaseStat int `json:"base_stat"`
	} `json:"stats"`
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Moves  []struct {
		VersionGroupDetails []struct {
			MoveLearnMethod struct {
				URL  string `json:"url"`
				Name string `json:"name"`
			} `json:"move_learn_method"`
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				URL  string `json:"url"`
				Name string `json:"name"`
			} `json:"version_group"`
		} `json:"version_group_details"`
		Move struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"move"`
	} `json:"moves"`
	Sprites struct {
		BackFemale       interface{} `json:"back_female"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		BackDefault      string      `json:"back_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
		BackShiny        string      `json:"back_shiny"`
		FrontDefault     string      `json:"front_default"`
		FrontShiny       string      `json:"front_shiny"`
	} `json:"sprites"`
	HeldItems              []interface{} `json:"held_items"`
	LocationAreaEncounters string        `json:"location_area_encounters"`
	Height                 int           `json:"height"`
	IsDefault              bool          `json:"is_default"`
	Species                struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"species"`
	ID          int `json:"id"`
	Order       int `json:"order"`
	GameIndices []struct {
		Version struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"version"`
		GameIndex int `json:"game_index"`
	} `json:"game_indices"`
	BaseExperience int `json:"base_experience"`
	Types          []struct {
		Slot int `json:"slot"`
		Type struct {
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}


var requestOptions grequests.RequestOptions

var pokemonType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Pokemon",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"weight": &graphql.Field{
				Type: graphql.Int,
			},
			"forms": &graphql.Field{
				Type: &graphql.List{
					OfType: formsType,
				},
			},
			"abilities": &graphql.Field{
				Type: &graphql.List{
					OfType: abilitiesType,
				},
			},
		},
	},
)

var formsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Forms",
		Fields: graphql.Fields{
			"url": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var abilitiesType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Abilities",
		Fields: graphql.Fields{
			"slot": &graphql.Field{
				Type: graphql.Int,
			},
			"is_hidden": &graphql.Field{
				Type: graphql.Boolean,
			},
			"ability": &graphql.Field{
				Type: formsType,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"pokemon": &graphql.Field{
				Type: pokemonType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					requestOptions.Headers = map[string]string{"Content-Type": "application/json"}
					resp, err := grequests.Get(fmt.Sprintf("http://pokeapi.co/api/v2/pokemon/%d", p.Args["id"].(int)), &requestOptions)
					var data pokemonResponse
					resp.JSON(&data)
					return data, err
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={pokemon(id:1){name}}'")
	http.ListenAndServe(":8080", nil)
}
