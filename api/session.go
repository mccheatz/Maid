package api

import (
	"errors"
	"maid/api/rest"
	"net/http"
	"strconv"
	"strings"
)

type X19Session struct {
	PatchList   []rest.NeteasePatch
	Release     rest.NeteaseReleaseInfo
	AuthServers []rest.NeteaseAuthServer
	LatestPatch string
	UserAgent   string
}

func (s *X19Session) CheckSessionAbility() error {
	if s.Release.ServerStop == "1" || s.Release.TempServerStop == 1 {
		return errors.New("game server under maintenance")
	}

	if len(s.AuthServers) == 0 {
		return errors.New("auth server offline")
	}

	return nil
}

func (s *X19Session) UpdateLatestPatch() {
	ver := -1
	latest := s.PatchList[0]

	for _, patch := range s.PatchList {
		info := strings.Split(patch.Name, ".")
		versionSeq := 0
		i, _ := strconv.Atoi(info[0])
		versionSeq += i << 24
		i, _ = strconv.Atoi(info[1])
		versionSeq += i << 20
		i, _ = strconv.Atoi(info[2])
		versionSeq += i << 16
		i, _ = strconv.Atoi(info[3])
		versionSeq += i

		if versionSeq > ver {
			latest = patch
			ver = versionSeq
		}
	}

	s.LatestPatch = latest.Name
	s.UserAgent = "WPFLauncher/" + s.LatestPatch
}

func EstablishSession(client *http.Client) (X19Session, error) {
	session := X19Session{}

	err := rest.PatchList(client, &session.PatchList)
	if err != nil {
		return session, err
	}

	err = rest.GameReleaseInfo(client, &session.Release)
	if err != nil {
		return session, err
	}

	err = rest.AuthServerList(client, session.Release, &session.AuthServers)
	if err != nil {
		return session, err
	}

	session.UpdateLatestPatch()

	return session, nil
}
