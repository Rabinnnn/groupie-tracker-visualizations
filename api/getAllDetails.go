package api

import(
	"net/http"
	"encoding/json"
)

func GetLocation(id string)(Location, error){
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/locations/" + id)
	if err != nil{
		return Location{}, err
	}
	defer results.Body.Close()

	var data Location
	if err := json.NewDecoder(results.Body).Decode(&data); err != nil{
		return Location{}, err
	}
	return data, nil
}

func GetDates(id string)(Date, error){
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/dates/" + id)
	if err != nil{
		return Date{}, err
	}
	defer results.Body.Close()

	var data Date
	if err := json.NewDecoder(results.Body).Decode(&data); err != nil {
		return Date{}, err
	}
	return data, nil
}

func GetRelations(id string)(Relations, error){
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + id)
	if err != nil {
		return Relations{}, err
	}
	defer results.Body.Close()

	var data Relations
	if err := json.NewDecoder(results.Body).Decode(&data); err != nil{
		return Relations{}, err
	}
	return data, nil
}

func GetDetails(id string)(Details, error){
	results, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil{
		return Details{}, err
	}
	defer results.Body.Close()

	var data Details
	if err := json.NewDecoder(results.Body).Decode(&data); err != nil{
		return Details{}, err
	}
	return data, nil
}

func GetAllDetails(id string)(AllDetails, error){
	details, err := GetDetails(id)
	if err != nil{
		return AllDetails{}, err
	}

	dates, err := GetDates(id)
	if err != nil{
		return AllDetails{}, err
	}

	relations, err := GetRelations(id)
	if err != nil{
		return AllDetails{}, err
	}

	location, err := GetLocation(id)
	if err != nil{
		return AllDetails{}, err
	}

	data := AllDetails{
		Details: details,
		Dates: dates,
		Location: location,
		Relations: relations,
		
	}
	return data, nil
}