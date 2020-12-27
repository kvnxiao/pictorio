package users

import (
	"errors"
	"net/http"

	"github.com/kvnxiao/pictorio/cookies"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/random"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
)

func Read(w http.ResponseWriter, req *http.Request) (ksuid.KSUID, string, error) {
	// Read user's unique ID or generate one if not exist

	userID, err := cookies.GetUserID(req)
	if err != nil || userID == "" {
		randomID, err := ksuid.NewRandom()
		if err != nil {
			return ksuid.KSUID{}, "", errors.New("could not generate a unique ID for a new user")
		}
		userID = randomID.String()
		cookies.SetUserID(w, userID)
	}

	// Ensure user ID parsed from cookies is of expected type
	userKSUID, err := ksuid.Parse(userID)
	if err != nil {
		return ksuid.KSUID{}, "", errors.New("could not parse user ID")
	}

	// Read user name
	userName, err := cookies.GetUserName(req)
	if err != nil || userName == "" {
		log.Debug().Msg("Generating random name for new user")
		userName = random.GenerateName()
		cookies.SetUserName(w, userName)
	}

	return userKSUID, userName, nil
}

func ReadName(w http.ResponseWriter, req *http.Request) (model.NameResponse, error) {
	_, userName, err := Read(w, req)
	if err != nil {
		return model.NameResponse{}, errors.New("unable to read name from cookies")
	}

	return model.NameResponse{
		Name:        userName,
		IsGenerated: random.IsGeneratedName(userName),
	}, nil
}

func ChangeName(newName string, w http.ResponseWriter) model.NameResponse {
	cookies.SetUserName(w, newName)
	return model.NameResponse{
		Name:        newName,
		IsGenerated: false,
	}
}
