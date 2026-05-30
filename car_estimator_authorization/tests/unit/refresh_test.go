package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nikita-itmo-gh-acc/car_estimator_authorization/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_authorization/services"
	"github.com/nikita-itmo-gh-acc/car_estimator_authorization/tests/unit/mocks"
	mock "github.com/stretchr/testify/mock"
)

func TestRefreshRejectsChangedSource(t *testing.T) {
	refreshToken := "refresh-token"
	userID := uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479")
	storedSource := domain.Source{
		IpAddress: "127.0.0.1:8000",
		UserAgent: "Chrome/137.0.0.0",
	}
	changedSource := domain.Source{
		IpAddress: storedSource.IpAddress,
		UserAgent: "Firefox/126.0",
	}
	session := &domain.Session{
		UserId: userID,
		Email: "test@test.ru",
		Source: storedSource,
	}

	userProvider := mocks.NewIUserProvider(t)
	sessionProvider := mocks.NewISessionProvider(t)
	sessionSaver := mocks.NewISessionSaver(t)
	sessionRemover := mocks.NewISessionRemover(t)

	sessionProvider.
		On("Get", mock.Anything, refreshToken).
		Return(session, nil)
	sessionRemover.
		On("Delete", mock.Anything, refreshToken).
		Return(nil)

	authService := services.NewAuthService(
		userProvider,
		sessionProvider,
		sessionSaver,
		sessionRemover,
		NullLogger(),
	)

	tokenPair, err := authService.Refresh(context.Background(), refreshToken, changedSource)

	if tokenPair != nil {
		t.Fatalf("expected no token pair, got %+v", tokenPair)
	}
	if !errors.Is(err, services.ErrSourceChanged) {
		t.Fatalf("expected ErrSourceChanged, got %v", err)
	}
}
