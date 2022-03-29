package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/model"
	"github.com/authelia/authelia/v4/internal/oidc"
	"github.com/authelia/authelia/v4/internal/session"
	"github.com/authelia/authelia/v4/internal/utils"
)

func oidcConsent(ctx *middlewares.AutheliaCtx) {
	userSession, consent, client, ret := oidcConsentGetSessionsAndClient(ctx)
	if ret {
		return
	}

	if !client.IsAuthenticationLevelSufficient(userSession.AuthenticationLevel) {
		ctx.Logger.Errorf("Unable to perform consent without sufficient authentication for user '%s' and client id '%s'", userSession.Username, consent.ClientID)
		ctx.Logger.Debugf("Insufficient permissions to give consent %d -> %d", userSession.AuthenticationLevel, client.Policy)
		ctx.ReplyForbidden()

		return
	}

	if err := ctx.SetJSONBody(client.GetConsentResponseBody(consent)); err != nil {
		ctx.Error(fmt.Errorf("unable to set JSON body: %v", err), "Operation failed")
	}
}

func oidcConsentPOST(ctx *middlewares.AutheliaCtx) {
	var (
		body ConsentPostRequestBody
		err  error
	)

	if err = json.Unmarshal(ctx.Request.Body(), &body); err != nil {
		ctx.Logger.Errorf("Failed to parse JSON body in consent POST: %+v", err)
		ctx.SetJSONError(messageOperationFailed)

		return
	}

	userSession, consent, client, ret := oidcConsentGetSessionsAndClient(ctx)
	if ret {
		return
	}

	if !client.IsAuthenticationLevelSufficient(userSession.AuthenticationLevel) {
		ctx.Logger.Debugf("Insufficient permissions to give consent v1 %d -> %d", userSession.AuthenticationLevel, userSession.OIDCWorkflowSession.RequiredAuthorizationLevel)
		ctx.ReplyForbidden()

		return
	}

	if consent.ClientID != body.ClientID {
		ctx.Logger.Errorf("User '%s' consented to scopes of another client (%s) than expected (%s). Beware this can be a sign of attack",
			userSession.Username, body.ClientID, consent.ClientID)
		ctx.SetJSONError(messageOperationFailed)

		return
	}

	form, err := consent.GetForm()
	if err != nil {
		ctx.Logger.Errorf("Could not parse the form stored in the consent session with challenge id '%s' for user '%s': %v", consent.ChallengeID.String(), userSession.Username, err)
		ctx.SetJSONError(messageOperationFailed)

		return
	}

	var (
		externalRootURL string
		authorized      = true
	)

	switch body.AcceptOrReject {
	case accept:
		if externalRootURL, err = ctx.ExternalRootURL(); err != nil {
			ctx.Logger.Errorf("Could not determine the external URL during consent session processing with challenge id '%s' for user '%s': %v", consent.ChallengeID.String(), userSession.Username, err)
			ctx.SetJSONError(messageOperationFailed)

			return
		}

		// redirectURI = fmt.Sprintf("%s%s?%s", externalRootURL, oidc.AuthorizationPath, form.Encode()).

		consent.GrantedScopes = consent.RequestedScopes
		consent.GrantedAudience = consent.RequestedAudience

		if !utils.IsStringInSlice(consent.ClientID, consent.GrantedAudience) {
			consent.GrantedAudience = append(consent.GrantedAudience, consent.ClientID)
		}

		/*
			if err = ctx.Providers.StorageProvider.SaveOAuth2ConsentSessionResponse(ctx, consent, false); err != nil {
				ctx.Logger.Errorf("Failed to save the consent session to the database: %+v", err)
				ctx.SetJSONError(messageOperationFailed)

				return
			}

		*/
	case reject:
		authorized = false
		/*
			redirectURIForm := url.Values{
				"error":             []string{"access_denied"},
				"error_description": []string{"User rejected the consent request"},
			}

			if state := form.Get("state"); state != "" {
				redirectURIForm.Set("state", state)
			}

			redirectURI = fmt.Sprintf("%s?%s", form.Get("redirect_uri"), redirectURIForm.Encode())

			if err = ctx.Providers.StorageProvider.SaveOAuth2ConsentSessionResponse(ctx, consent, true); err != nil {
				ctx.Logger.Errorf("Failed to save the consent session to the database: %+v", err)
				ctx.SetJSONError(messageOperationFailed)

				return
			}

			userSession.ConsentChallengeID = nil

			if err = ctx.SaveSession(userSession); err != nil {
				ctx.Logger.Errorf("Failed to save the user session: %+v", err)
				ctx.SetJSONError(messageOperationFailed)

				return
			}

		*/
	default:
		ctx.Logger.Warnf("User '%s' tried to reply to consent with an unexpected verb", userSession.Username)
		ctx.ReplyBadRequest()

		return
	}

	if err = ctx.Providers.StorageProvider.SaveOAuth2ConsentSessionResponse(ctx, consent, authorized); err != nil {
		ctx.Logger.Errorf("Failed to save the consent session response to the database: %+v", err)
		ctx.SetJSONError(messageOperationFailed)

		return
	}

	response := ConsentPostResponseBody{RedirectURI: fmt.Sprintf("%s%s?%s", externalRootURL, oidc.AuthorizationPath, form.Encode())}

	if err = ctx.SetJSONBody(response); err != nil {
		ctx.Error(fmt.Errorf("unable to set JSON body in response"), "Operation failed")
	}
}

func oidcConsentGetSessionsAndClient(ctx *middlewares.AutheliaCtx) (userSession session.UserSession, consent *model.OAuth2ConsentSession, client *oidc.Client, ret bool) {
	var (
		err error
	)

	userSession = ctx.GetSession()

	if userSession.ConsentChallengeID == nil {
		ctx.Logger.Errorf("Cannot consent for user '%s' when OIDC consent session has not been initiated", userSession.Username)
		ctx.ReplyForbidden()

		return userSession, nil, nil, true
	}

	if consent, err = ctx.Providers.StorageProvider.LoadOAuth2ConsentSessionByChallengeID(ctx, *userSession.ConsentChallengeID); err != nil {
		ctx.Logger.Errorf("Unable to load consent session with challenge id '%s': %v", userSession.ConsentChallengeID.String(), err)
		ctx.ReplyForbidden()

		return userSession, nil, nil, true
	}

	if client, err = ctx.Providers.OpenIDConnect.Store.GetFullClient(consent.ClientID); err != nil {
		ctx.Logger.Errorf("Unable to find related client configuration with name '%s': %v", consent.ClientID, err)
		ctx.ReplyForbidden()

		return userSession, nil, nil, true
	}

	return userSession, consent, client, false
}
