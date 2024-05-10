package data

import "github.com/kurakura967/testcontainers-for-elasticsearch/model"

var DummyDocuments = model.Documents{
	{
		ID:     "1",
		Title:  "title1",
		Artist: "artist1",
	},
	{
		ID:     "2",
		Title:  "title2",
		Artist: "artist2",
	},
}
