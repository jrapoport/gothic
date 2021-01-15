package api

import (
	"context"
	"net/http"

	"github.com/jrapoport/gothic/api/provider"
)

func (a *API) loadSAMLState(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	state := r.FormValue("RelayState")
	if state == "" {
		return nil, badRequestError("saml relay state is missing")
	}

	ctx := r.Context()

	return a.loadExternalState(ctx, state)
}

func (a *API) samlCallback(r *http.Request, ctx context.Context) (*provider.UserProvidedData, error) {
	config := a.config

	samlProvider, err := provider.NewSamlProvider(config.External.Saml, a.db)
	if err != nil {
		return nil, badRequestError("could not initialize saml provider: %+v", err).WithInternalError(err)
	}

	samlResponse := r.FormValue("SAMLResponse")
	if samlResponse == "" {
		return nil, badRequestError("saml Response is missing")
	}

	assertionInfo, err := samlProvider.ServiceProvider.RetrieveAssertionInfo(samlResponse)
	if err != nil {
		return nil, internalServerError("parsing saml assertion failed: %+v", err).WithInternalError(err)
	}

	if assertionInfo.WarningInfo.InvalidTime {
		return nil, forbiddenError("saml response has invalid time")
	}

	if assertionInfo.WarningInfo.NotInAudience {
		return nil, forbiddenError("saml response is not in audience")
	}

	if assertionInfo == nil {
		return nil, internalServerError("saml Assertion is missing")
	}
	userData := &provider.UserProvidedData{
		Emails: []provider.Email{{
			Email:    assertionInfo.NameID,
			Verified: true,
		}},
	}
	return userData, nil
}

func (a *API) SAMLMetadata(w http.ResponseWriter, r *http.Request) error {
	external := a.config.External
	samlProvider, err := provider.NewSamlProvider(external.Saml, a.db)
	if err != nil {
		return internalServerError("Could not create SAML Provider: %+v", err).WithInternalError(err)
	}

	metadata, err := samlProvider.SPMetadata()
	w.Header().Set("Content-Type", "application/xml")
	_, err = w.Write(metadata)
	return err
}
