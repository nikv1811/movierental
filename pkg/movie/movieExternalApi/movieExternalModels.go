package movieExternalApi

type ListMoviesResponse struct {
	Status        string                `json:"status"`
	StatusMessage string                `json:"status_message"`
	Data          MovieListResponseData `json:"data"`
}

type MovieListResponseData struct {
	MovieCount int     `json:"movie_count"`
	Limit      int     `json:"limit"`
	PageNumber int     `json:"page_number"`
	Movies     []Movie `json:"movies"`
}

type MovieResponse struct {
	Status        string            `json:"status"`
	StatusMessage string            `json:"status_message"`
	Data          MovieResponseData `json:"data"`
}

type MovieResponseData struct {
	Movie Movie `json:"movie"`
}

type Movie struct {
	ID       int      `json:"id"`
	IMDBCode string   `json:"imdb_code"`
	Year     int      `json:"year"`
	Genres   []string `json:"genres"`
	Title    string   `json:"title"`
}
